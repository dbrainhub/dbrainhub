package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dbrainhub/dbrainhub/agent"
	"github.com/dbrainhub/dbrainhub/configs"
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
	reporter, err := agent.NewStartupReporter(config)
	if err != nil {
		logger.Errorf("new server reporter error, exit...")
		return
	}
	if err := reporter.Report(ctx); err != nil {
		logger.Errorf("connect to server error, exit...")
		return
	}

	filebeatService, err := filebeat.NewFilebeatService(config)
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