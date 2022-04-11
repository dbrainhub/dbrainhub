package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/elastic/beats/v7/libbeat/publisher"
	esClient "github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

func (d *DefaultController) DbRainhubOutput(c *gin.Context, req DbRainhubRequest) (*DbRainhubResponse, error) {
	// TODO: limiter

	events := req.Events
	//dbIp := req.DbIp
	//dbPort := req.DbPort
	//// check db member instance belongs to a cluster
	//db := model.GetDB(c)
	//member, err := model.GetDbClusterMemberByIpAndPort(c, db, dbIp, int16(dbPort))
	//if err != nil {
	//	return nil, err
	//}
	//if member == nil || member.ClusterId == 0 {
	//	return nil, errors.DbMemberNotClassified("db member should be assigned to a cluster first.")
	//}

	cfg := configs.GetGlobalConfig().Es
	esCfg := esClient.Config{
		Addresses: cfg.Addresses,
	}
	es, err := esClient.NewClient(esCfg)
	if err != nil {
		return nil, err
	}

	var failedEvents []int32
	var buf bytes.Buffer
	for i, eve := range events {
		eb, err := json.Marshal(eve)
		if err != nil {
			failedEvents = append(failedEvents, int32(i))
			continue
		}

		eb = append(eb, "\n"...)
		// FIXME: index name
		meta := []byte(fmt.Sprintf(`{ "index" : { "_index":"test-bulk123" } }%s`, "\n"))
		buf.Grow(len(meta) + len(eb))
		buf.Write(meta)
		buf.Write(eb)
	}
	_, err = es.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		// bulk failed, all events failed to retry in filebeat.
		for i:=0; i<len(events); i++ {
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
