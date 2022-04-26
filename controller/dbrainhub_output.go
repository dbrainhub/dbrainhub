package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
	"github.com/dbrainhub/dbrainhub/utils/rate_limit"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/beat/events"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/publisher"
	esClient "github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

var limiter *rate_limit.SlidingWindowRateLimiter
var es *esClient.Client
var once = sync.Once{}

func GetRateLimiterAndEsClient(ctx context.Context) (*rate_limit.SlidingWindowRateLimiter, *esClient.Client) {
	var err error
	once.Do(func() {
		cfg := configs.GetGlobalServerConfig().OutputServer
		limiter = rate_limit.NewSlidingWindowRateLimiter(int64(cfg.QpsThreshold))

		esCfg := esClient.Config{
			Addresses: cfg.EsAddresses,
		}
		es, err = esClient.NewClient(esCfg)
		if err != nil {
			panic(err)
		}
	})
	return limiter, es
}

func DbRainhubOutput(c *gin.Context, req DbRainhubRequest) (*DbRainhubResponse, error) {
	// limiter check
	limiter, es := GetRateLimiterAndEsClient(c)
	err := limiter.Limit()
	if err == rate_limit.ErrRateLimited {
		return nil, errors.FilebeatRateLimited("dbrainhub output triggers the rate limit.")
	}

	events := req.Events
	dbIp := req.DbIp
	dbPort := req.DbPort
	// check db member instance belongs to a cluster
	db := model.GetDB(c)
	member, err := model.GetDbClusterMemberByIpAndPort(c, db, dbIp, int16(dbPort))
	if err != nil {
		return nil, err
	}
	if member == nil || member.ClusterId == 0 {
		return nil, errors.DbClusterMemberNotAssigned("db member should be assigned to a cluster first.")
	}

	cluster, err := model.GetDbClusterById(c, db, member.ClusterId)
	if err != nil {
		return nil, err
	}
	clusterName := cluster.Name
	clusterType := cluster.DbType

	var failedEvents []int32
	var buf bytes.Buffer

	for i, eve := range events {
		data := eve.Content.Fields
		data["cluster"] = clusterName
		data["instance"] = fmt.Sprintf("%s:%d", dbIp, dbPort)
		eb, err := json.Marshal(data)
		if err != nil {
			failedEvents = append(failedEvents, int32(i))
			continue
		}

		indexName, err := genIndex(clusterType, &eve.Content)
		if err != nil {
			failedEvents = append(failedEvents, int32(i))
			continue
		}
		pipeline, err := getPipeline(&eve.Content)
		if err != nil {
			failedEvents = append(failedEvents, int32(i))
			continue
		}
		meta := []byte(fmt.Sprintf(`{ "create" : { "_index":"%s", "pipeline":"%s" } }%s`, indexName, pipeline, "\n"))
		buf.Grow(len(meta) + len(eb) + 1)
		buf.Write(meta)
		buf.Write(eb)
		buf.WriteByte('\n')
	}
	_, err = es.Bulk(&buf)
	if err != nil {
		// bulk failed, all events failed to retry in filebeat.
		for i := 0; i < len(events); i++ {
			failedEvents = append(failedEvents, int32(i))
		}
		return nil, err
	}

	return &DbRainhubResponse{
		FailedEvents: failedEvents,
	}, nil
}

type DbRainhubRequest struct {
	Events []publisher.Event `json:"events"`
	DbIp   string            `json:"db_ip"`
	DbPort int               `json:"db_port"`
}

type DbRainhubResponse struct {
	FailedEvents []int32 `json:"failed_events"`
}

// Index name
func genIndex(cluster string, event *beat.Event) (string, error) {
	logType, err := event.Fields.GetValue("fileset.name")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s", cluster, logType), nil
}

func getPipeline(event *beat.Event) (string, error) {
	if event.Meta != nil {
		pipeline, err := events.GetMetaStringValue(*event, events.FieldMetaPipeline)
		if err == common.ErrKeyNotFound {
			return "", nil
		}
		if err != nil {
			return "", err
		}
		return strings.ToLower(pipeline), nil
	}
	return "", nil
}
