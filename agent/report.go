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
	"github.com/dbrainhub/dbrainhub/dbs"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/utils"
	"github.com/dbrainhub/dbrainhub/utils/logger"
	osutils "github.com/dbrainhub/dbrainhub/utils/os"
)

const StartupReportPath = "/agent/report"

// report after startup
type StartupReporter interface {
	Report(ctx context.Context) error
}

func NewStartupReporter(agentConf *configs.AgentConfig, df dbs.DBOperationFactory) (StartupReporter, error) {
	dbtype, err := agentConf.ConvertDBType()
	if err != nil {
		logger.Errorf("get dbtype from agentconf error, err: %v", err)
		return nil, errors.AgentConfigError("db_type error")
	}

	dbVersionQuerier, err := df.CreateVersionQuerier()
	if err != nil {
		logger.Errorf("createVersionQuerier error, err: %v", err)
		return nil, err
	}

	httpClient := utils.NewHttpClient(time.Millisecond*time.Duration(agentConf.Server.TimeoutMs),
		agentConf.Server.Retry,
		time.Duration(agentConf.Server.RetryIntervalMs)*time.Millisecond)

	return &startupReportImpl{
		hostType:         agentConf.ConvertHostType(),
		dbType:           dbtype,
		env:              agentConf.DB.Env,
		port:             agentConf.DB.Port,
		serverAddr:       agentConf.Server.Addr,
		client:           httpClient,
		dbVersionQuerier: dbVersionQuerier,
	}, nil
}

type startupReportImpl struct {
	dbType   api.StartupReportRequest_DBType
	hostType api.StartupReportRequest_HostType
	port     int
	env      string

	serverAddr       string
	dbVersionQuerier dbs.DBVersionQuerier
	client           utils.HttpClient
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
	dbVersion, err := s.dbVersionQuerier.Query(ctx)
	if err != nil { // don't return
		logger.Errorf("queryDBVersion error, err: %v", err)
	}

	versionQuerier := osutils.NewVersionQuerier(runtime.GOOS)
	osVersion, err := versionQuerier.GetOsVersion()
	if err != nil { // don't return
		logger.Errorf("get os_version error, err: %v", err)
	}

	kernelVersion, err := versionQuerier.GetKernelVersion()
	if err != nil { // don't return
		logger.Errorf("get kernel_version error, err: %v", err)
	}

	req := &api.StartupReportRequest{
		DbType:        api.StartupReportRequest_DBType(s.dbType),
		HostType:      api.StartupReportRequest_HostType(s.hostType),
		Hostname:      hostname,
		IpAddr:        localip,
		Port:          int32(s.port),
		Os:            runtime.GOOS,
		OsVersion:     osVersion,
		DbVersion:     dbVersion.Version,
		KernelVersion: kernelVersion,
		Env:           s.env,
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
