package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
	"github.com/dbrainhub/dbrainhub/utils/search_time"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
)

const (
	DefaultSize = 20
)

func DbRainhubSearchMemberLogsWithCount(c *gin.Context, req api.SearchMemberLogCountRequest) (*api.SearchMemberLogCountResponse, error) {
	err := validateReq(req)
	if err != nil {
		return nil, err
	}

	dbIp := req.DbIp
	dbPort := req.DbPort
	db := model.GetDB(c)
	member, err := model.GetDbClusterMemberByIpAndPort(c, db, dbIp, int16(dbPort))
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.DbClusterMemberNotFoundByIpAndPort(dbIp, int16(dbPort))
	}

	dbType := member.DbType
	indexName := fmt.Sprintf("%s-%s", dbType, req.Type)

	size := req.Size
	if size == 0 {
		size = DefaultSize
	}

	// interval
	interval := search_time.GetInterval(req.StartTime, req.EndTime)

	var buf bytes.Buffer
	query := map[string]interface{}{
		"from":             req.From,
		"size":             size,
		"track_total_hits": true,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []interface{}{
					// bool-filter
					map[string]interface{}{"bool": map[string]interface{}{
						"filter": []interface{}{
							map[string]interface{}{"match": map[string]interface{}{
								"service.type": dbType,
							}},
							map[string]interface{}{"match": map[string]interface{}{
								"instance": fmt.Sprintf("%s:%d", dbIp, dbPort),
							}},
						},
					}},

					// range
					map[string]interface{}{"range": map[string]interface{}{
						"@timestamp": map[string]interface{}{
							"format": "strict_date_optional_time",
							"gte":    req.StartTime,
							"lte":    req.EndTime,
						},
					}},
				},
			},
		},
		"aggs": map[string]interface{}{
			"logs": map[string]interface{}{
				"date_histogram": map[string]interface{}{
					"field":          "@timestamp",
					"fixed_interval": interval,
					"time_zone":      "Asia/Shanghai",
					"min_doc_count":  1,
					"format":         "yyyy-MM-dd hh:mm:ss",
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	_, es := GetRateLimiterAndEsClient(c)
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(indexName),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			return nil, errors.FilebeatSearchError(e["error"].(map[string]interface{})["reason"].(string))
		}
	}

	searchRes := new(api.SearchMemberLogCountResponse)
	unm := jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
	if err := unm.Unmarshal(res.Body, searchRes); err != nil {
		return nil, err
	}
	searchRes.From = req.From
	searchRes.Size = req.Size
	return searchRes, nil
}

func validateReq(req api.SearchMemberLogCountRequest) error {
	// time validate
	s := req.StartTime
	st, err := time.ParseInLocation("2006-01-02T15:04:05Z", s, time.Local)
	if err != nil {
		return errors.FilebeatSearchParamsError("start_time format should be like '2006-01-02T15:04:05Z'.")
	}
	e := req.EndTime
	et, err := time.ParseInLocation("2006-01-02T15:04:05Z", e, time.Local)
	if err != nil {
		return errors.FilebeatSearchParamsError("end_time format should be like '2006-01-02T15:04:05Z'.")
	}
	if et.Sub(st) < 0 {
		return errors.FilebeatSearchParamsError("end_time should be larger than start_time.")
	}

	// ip validate
	ip := net.ParseIP(req.DbIp)
	if ip == nil {
		return errors.FilebeatSearchParamsError("db_ip param is wrong.")
	}

	// port validate
	if req.DbPort <= 0 || req.DbPort > 65535 {
		return errors.FilebeatSearchParamsError("db_port should be in [1, 65535].")
	}

	// size and from validate
	if req.Size < 0 || req.From < 0 {
		return errors.FilebeatSearchParamsError("size or from should be positive.")
	}
	return nil
}
