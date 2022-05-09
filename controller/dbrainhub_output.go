package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
	"github.com/dbrainhub/dbrainhub/model/es"
	"github.com/dbrainhub/dbrainhub/server"
	"github.com/dbrainhub/dbrainhub/utils/rate_limit"
	"github.com/dbrainhub/dbrainhub/utils/search_time"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/beat/events"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/publisher"
	"github.com/gin-gonic/gin"
)

func DbRainhubOutput(c *gin.Context, req DbRainhubRequest) (*DbRainhubResponse, error) {
	// limiter check
	client := server.GetDefaultESClientWithRateLimiter()
	err := client.Limit()
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
	var msgs []*es.ESMessage

	for i, eve := range events {
		data := eve.Content.Fields
		data["cluster"] = clusterName
		data["instance"] = fmt.Sprintf("%s:%d", dbIp, dbPort)
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
		msgs = append(msgs, &es.ESMessage{
			Meta: &es.ESMeta{
				Index:    indexName,
				Pipeline: pipeline,
			},
			Data: data,
		})
	}

	err = client.Send(msgs)
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
	// get week
	_, week, err := search_time.GetYearAndWeek(time.Now().Format("2006-01-02T15:04:05Z"))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s-%dw", cluster, logType, week), nil
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
