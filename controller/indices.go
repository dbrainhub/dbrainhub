package controller

import (
	"context"
	"sync"
	"time"

	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model/es"
	"github.com/dbrainhub/dbrainhub/server"
	"github.com/dbrainhub/dbrainhub/utils/logger"
	"github.com/gin-gonic/gin"
)

func GetInstanceIndices(c *gin.Context, req *api.GetInstanceIndicesRequest) (*api.GetInstanceIndicesResponse, error) {
	begin, err := time.ParseInLocation("2006-01-02 15:04:05", req.From, time.Local)
	if err != nil {
		return nil, errors.InvalidTimeDuration("invalid from: %s", req.From)
	}
	end, err := time.ParseInLocation("2006-01-02 15:04:05", req.To, time.Local)
	if err != nil {
		return nil, errors.InvalidTimeDuration("invalid to: %s", req.To)
	}

	var res api.GetInstanceIndicesResponse
	var wg sync.WaitGroup

	queryAvg := func(ctx context.Context, field string, res *[]*api.GetInstanceIndicesResponse_IndexValue) {
		defer wg.Done()

		client := es.NewIndexQuerier(es.GetESClient())
		avgs, err := client.Query(ctx, &es.IndexQuerierParam{
			Index:    server.GetIndicesIndexName(),
			AggField: field,
			Begin:    begin,
			End:      end,
			Buckets:  req.Buckets,
			Cond:     map[string]interface{}{"ip": req.Host, "port": req.Port},
		})
		if err != nil {
			logger.Errorf("GetInstanceIndices AggregateAvg %s error, err: %v, req: %v", field, err, req)
			return
		}
		*res = esAvgResultsToApiIndexValues(avgs)
	}

	wg.Add(5)
	go queryAvg(c, "cpu_ratio", &res.CpuRatios)
	go queryAvg(c, "mem_ratio", &res.MemRatios)
	go queryAvg(c, "disk_ratio", &res.DiskRatios)
	go queryAvg(c, "qps", &res.Qps)
	go queryAvg(c, "tps", &res.Tps)
	wg.Wait()

	return &res, err

}

func esAvgResultsToApiIndexValues(avgResults []*es.IndexQuerierResult) []*api.GetInstanceIndicesResponse_IndexValue {
	var res []*api.GetInstanceIndicesResponse_IndexValue
	for _, avgRes := range avgResults {
		res = append(res, &api.GetInstanceIndicesResponse_IndexValue{
			Count:     int64(avgRes.Count),
			Value:     avgRes.Value,
			StartTime: avgRes.StartTime,
		})
	}
	return res
}
