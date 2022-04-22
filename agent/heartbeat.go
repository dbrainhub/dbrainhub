package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/dbrainhub/dbrainhub/dbs"
	"github.com/dbrainhub/dbrainhub/utils"
	"github.com/dbrainhub/dbrainhub/utils/logger"
	osutils "github.com/dbrainhub/dbrainhub/utils/os"
)

const HeartbeatPath = "/agent/heartbeat"

type HeartbeatService interface {
	Run(ctx context.Context)
}

func NewHeartbeatService(ctx context.Context, agentConf *configs.AgentConfig, dbfactory dbs.DBOperationFactory) (HeartbeatService, error) {
	httpClient := utils.NewHttpClient(time.Millisecond*time.Duration(agentConf.Server.TimeoutMs),
		agentConf.Server.Retry,
		time.Duration(agentConf.Server.RetryIntervalMs)*time.Millisecond)

	dbVariableCenter, err := dbfactory.CreateVariablesCenter()
	if err != nil {
		logger.Errorf("create DBVariablesCenter failed, err: %v", err)
		return nil, err
	}

	dbIndexManager, err := createDBIndexManager(ctx, dbfactory)
	if err != nil {
		logger.Errorf("create DBIndexManager failed, err: %v", err)
		return nil, err
	}

	return &heartbeatImpl{
		interval:         time.Duration(agentConf.Server.HeartbeatIntervalMs) * time.Millisecond,
		serverAddr:       agentConf.Server.Addr,
		dbType:           agentConf.DB.DBType,
		Port:             agentConf.DB.Port,
		client:           httpClient,
		cpuUtils:         osutils.NewCPUUitls(ctx, 0),
		diskUtils:        osutils.NewDiskUitls(),
		memUtils:         osutils.NewMemoryUitls(),
		dbVariableCenter: dbVariableCenter,
		dbIndexManager:   dbIndexManager,
	}, nil
}

type heartbeatImpl struct {
	interval   time.Duration
	serverAddr string

	dbVariableCenter dbs.VariablesCenter
	dbIndexManager   dbs.DBIndexManager

	client    utils.HttpClient
	cpuUtils  osutils.CPUUtils
	memUtils  osutils.MemoryUtils
	diskUtils osutils.DiskUtils

	dbType string
	Port   int
}

func (h *heartbeatImpl) Run(ctx context.Context) {
	tick := time.NewTicker(h.interval) // time.Ticker (not time.Sleep) avoids heartbeat timeout
	for {
		select {
		case <-tick.C:
			h.heartbeat(ctx)
		}
	}
}

func (h *heartbeatImpl) heartbeat(ctx context.Context) {
	// for os
	ip, err := utils.GetLocalIP()
	if err != nil {
		logger.Errorf("heartbeat get local ip error, err: %v", err)
		return
	}
	cpuUsage, err := h.cpuUtils.Usage(ctx)
	if err != nil {
		logger.Errorf("heartbeat get cpuUsage error, err: %v", err)
	}
	memUsage, err := h.memUtils.Usage(ctx)
	if err != nil {
		logger.Errorf("heartbeat get memUsage error, err: %v", err)
	}
	diskUsage, err := h.diskUsage(ctx)
	if err != nil {
		logger.Errorf("heartbeat get diskUsage error, err: %v", err)
	}

	// for db
	dbQps, err := h.dbIndexManager.GetQPS(ctx)
	if err != nil {
		logger.Errorf("heartbeat get dbQPS error, err: %v", err)
	}
	dbTps, err := h.dbIndexManager.GetTPS(ctx)
	if err != nil {
		logger.Errorf("heartbeat get dbTPS error, err: %v", err)
	}

	req := &api.HeartbeatRequest{
		AgentInfo: &api.HeartbeatRequest_Agent{
			Localip:  ip,
			CpuRate:  cpuUsage,
			MemRate:  memUsage,
			DiskRate: diskUsage,
		},
		DbInfo: &api.HeartbeatRequest_DB{
			Dbtype: h.dbType,
			Port:   int32(h.Port),
			Qps:    dbQps,
			Tps:    dbTps,
		},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("heartbeat marshal request error, err: %v, req: %#v", err, req)
		return
	}
	_, err = h.client.Send(ctx, h.heartbeatUrl(), "POST", string(reqBytes))
	if err != nil {
		logger.Errorf("heartbeat send error, err: %v, req: %#v", err, req)
		return
	}

}

func (h *heartbeatImpl) heartbeatUrl() string {
	return fmt.Sprintf("http://%s%s", h.serverAddr, HeartbeatPath)
}

// the disk of db data path
func (h *heartbeatImpl) diskUsage(ctx context.Context) (float64, error) {
	diskDir, err := h.dbVariableCenter.DataDir(ctx)
	if err != nil {
		return 0.0, err
	}
	return h.diskUtils.Usage(ctx, diskDir)
}

func createDBIndexManager(ctx context.Context, dbfactory dbs.DBOperationFactory) (dbs.DBIndexManager, error) {
	statusQuerier, err := dbfactory.CreateStatusCenter()
	if err != nil {
		return nil, err
	}
	return dbs.NewDBIndexManager(ctx, statusQuerier)
}
