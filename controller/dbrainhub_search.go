package controller

import (
	"fmt"
	"net"
	"time"

	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
	"github.com/dbrainhub/dbrainhub/model/es"
	"github.com/dbrainhub/dbrainhub/utils/search_time"
	"github.com/gin-gonic/gin"
)

const (
	DefaultSize    int64 = 20
	DefaultBuckets int64 = 50
)

func DbRainhubSearchMemberLogsWithCount(c *gin.Context, req *api.SearchMemberLogCountRequest) (*api.SearchMemberLogCountResponse, error) {
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
	indexName := fmt.Sprintf("%s-%s-*", dbType, req.Type)

	size := req.Size
	if size == 0 {
		size = DefaultSize
	}

	// interval
	buckets := DefaultBuckets
	if req.Buckets > 0 {
		buckets = req.Buckets
	}
	param := &es.SlowlogQuerierParam{
		Index:   indexName,
		Begin:   req.StartTime,
		End:     req.EndTime,
		Buckets: buckets,
		Size:    size,
		From:    req.From,
		Cond: map[string]interface{}{
			"service.type": dbType,
			"instance":     fmt.Sprintf("%s:%d", dbIp, dbPort),
		},
	}
	searchRes, err := es.NewSlowlogQuerier(es.GetESClient()).Query(c, param)
	if err != nil {
		return nil, err
	}
	searchRes.From = req.From
	searchRes.Size = req.Size
	searchRes.Aggregations.Interval = search_time.GetInterval(req.StartTime, req.EndTime, buckets)
	return searchRes, nil
}

func validateReq(req *api.SearchMemberLogCountRequest) error {
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
