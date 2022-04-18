package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dbrainhub/dbrainhub/agent"
	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/dbrainhub/dbrainhub/dbs"
	"github.com/dbrainhub/dbrainhub/dbs/mysql"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/filebeat"
	"github.com/dbrainhub/dbrainhub/utils/logger"
)

var configFilePath = flag.String("config", "", "config file path, refer to example_agent_config.toml")

func main() {
	flag.Parse()
	if *configFilePath == "" {
		fmt.Fprintf(os.Stderr, "ERROR: config file required\n")
		os.Exit(1)
	}

	configs.InitConfigOrPanic(*configFilePath, configs.GetGlobalAgentConfig())

	config := configs.GetGlobalAgentConfig()
	logger.InitLog(config.LogInfo.LogDir, config.LogInfo.Name, config.LogInfo.Level)

	ctx := context.Background()
	dbOperationFactory, err := newDBOperationFactory(config)
	if err != nil {
		logger.Errorf("new db handler error(%v), exit...", err)
		return
	}
	reporter, err := agent.NewStartupReporter(config, dbOperationFactory)
	if err != nil {
		logger.Errorf("new server reporter error, exit...")
		return
	}
	if err := reporter.Report(ctx); err != nil {
		logger.Errorf("connect to server error, exit...")
		return
	}

	filebeatService, err := filebeat.NewFilebeatService(config, dbOperationFactory)
	if err != nil {
		logger.Errorf("create FilebeatService failed , exit...")
		return
	}
	if err := filebeatService.StartGatherSlowlog(ctx); err != nil {
		logger.Errorf("filebeat service StartGatherSlowlog error, exit...")
		return
	}

	heartbeatService := agent.NewHeartbeatService(config)
	heartbeatService.Run(ctx)
}

func newDBOperationFactory(config *configs.AgentConfig) (dbs.DBOperationFactory, error) {
	switch config.DB.DBType {
	case mysql.MysqlType:
		return mysql.NewMysqlOperationFactory(&dbs.DBInfo{
			IP:     "127.0.0.1",
			Port:   config.DB.Port,
			User:   config.DB.User,
			Passwd: config.DB.Passwd,
		}), nil
	}
	return nil, errors.AgentConfigError("unsupported dbtype: %s", config.DB.DBType)
}
