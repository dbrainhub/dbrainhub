package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
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

	client, err := model.NewEsClient(configs.GetGlobalServerConfig().OutputServer.EsAddresses)
	if err != nil {
		logger.Errorf("GetInstanceIndices NewEsClient error, err: %v", err)
		return nil, errors.ConnectToESFailed("NewEsClient error")
	}

	res := api.GetInstanceIndicesResponse{
		CpuRatios:  []*api.GetInstanceIndicesResponse_IndexValue{},
		DiskRatios: []*api.GetInstanceIndicesResponse_IndexValue{},
		MemRatios:  []*api.GetInstanceIndicesResponse_IndexValue{},
		Qps:        []*api.GetInstanceIndicesResponse_IndexValue{},
		Tps:        []*api.GetInstanceIndicesResponse_IndexValue{},
	}
	str, _ := json.Marshal(res.CpuRatios)
	fmt.Println(string(str))
	var wg sync.WaitGroup

	queryAvg := func(ctx context.Context, field string, res *[]*api.GetInstanceIndicesResponse_IndexValue) {
		defer wg.Done()
		avgs, err := client.AggregateAvg(ctx, field, begin, end, map[string]interface{}{"ip": req.Host, "port": req.Port})
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

func esAvgResultsToApiIndexValues(avgResults []*model.AvgAggsResult) []*api.GetInstanceIndicesResponse_IndexValue {
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
