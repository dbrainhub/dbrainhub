package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/dbrainhub/dbrainhub/utils"
	"github.com/dbrainhub/dbrainhub/utils/logger"
)

const HeartbeatPath = "/heartbeat"

type HeartbeatService interface {
	Run(ctx context.Context)
}

func NewHeartbeatService(agentConf *configs.AgentConfig) HeartbeatService {
	httpClient := utils.NewHttpClient(time.Millisecond*time.Duration(agentConf.Heartbeat.Timeout),
		agentConf.Heartbeat.Retry,
		time.Duration(agentConf.Heartbeat.RetryInterval)*time.Millisecond)
	return &heartbeatImpl{
		interval:   time.Duration(agentConf.Heartbeat.Interval) * time.Millisecond,
		serverAddr: agentConf.Server.Addr,
		dbType:     agentConf.DB.DBType,
		Port:       agentConf.DB.Port,
		client:     httpClient,
	}
}

type heartbeatImpl struct {
	interval   time.Duration
	client     utils.HttpClient
	serverAddr string

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
	ip, err := utils.GetLocalIP()
	if err != nil {
		logger.Errorf("heartbeat get local ip error, err: %v", err)
		return
	}

	req := &api.HeartbeatRequest{
		AgentInfo: &api.HeartbeatRequest_AgentInfo{
			Localip: ip,
		},
		DbInfo: &api.HeartbeatRequest_DBInfo{
			Dbtype: h.dbType,
			Port:   int32(h.Port),
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
