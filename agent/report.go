package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/utils"
	"github.com/dbrainhub/dbrainhub/utils/logger"
)

const StartupReportPath = "/report"

// report after startup
type StartupReporter interface {
	Report(ctx context.Context) error
}

func NewStartupReporter(agentConf *configs.AgentConfig) (StartupReporter, error) {
	dbtype, err := agentConf.ConvertDBType()
	if err != nil {
		logger.Errorf("get dbtype from agentconf error, err: %v", err)
		return nil, errors.AgentConfigError("db_type error")
	}

	httpClient := utils.NewHttpClient(time.Millisecond*time.Duration(agentConf.Server.Timeout),
		agentConf.Server.Retry,
		time.Duration(agentConf.Server.RetryInterval)*time.Millisecond)

	return &startupReportImpl{
		hostType:   agentConf.ConvertHostType(),
		dbType:     dbtype,
		port:       agentConf.DB.Port,
		serverAddr: agentConf.Server.Addr,
		client:     httpClient,
	}, nil
}

type startupReportImpl struct {
	dbType   api.StartupReportRequest_DBType
	hostType api.StartupReportRequest_HostType
	port     int

	serverAddr string
	client     utils.HttpClient
}

func (s *startupReportImpl) Report(ctx context.Context) error {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Errorf("startup report get hostname error, err:%v", err)
	}
	localip, err := utils.GetLocalIP()
	if err != nil {
		logger.Errorf("getlocalIP in startupReport error, err: %v", err)
		return err
	}
	req := &api.StartupReportRequest{
		DbType:    api.StartupReportRequest_DBType(s.dbType),
		HostType:  api.StartupReportRequest_HostType(s.hostType),
		Hostname:  hostname,
		IpAddr:    localip,
		Port:      int32(s.port),
		Os:        runtime.GOOS,
		OsVersion: "", // TODO: 不同系统的获取方式不同。
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("json marshal startup report request error, err: %v, req: %#v", err, req)
		return err
	}

	_, err = s.client.Send(ctx, s.reportUrl(), "POST", string(reqBytes))
	if err != nil {
		logger.Errorf("report to server error, err: %v", err)
	}
	return err
}

func (s *startupReportImpl) reportUrl() string {
	return fmt.Sprintf("http://%s%s", s.serverAddr, StartupReportPath)
}
