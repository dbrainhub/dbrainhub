package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/dbrainhub/dbrainhub/router"
	"github.com/dbrainhub/dbrainhub/utils/logger"
)

var configFilePath = flag.String("config", "", "config file path, refer to example_config.toml")

func main() {
	flag.Parse()
	if *configFilePath == "" {
		fmt.Fprintf(os.Stderr, "ERROR: config file required\n")
		os.Exit(1)
	}

	configs.InitConfigOrPanic(*configFilePath, configs.GetGlobalServerConfig())

	config := configs.GetGlobalServerConfig()
	logger.InitLog(config.LogInfo.LogDir, config.LogInfo.Name, config.LogInfo.Level)

	logger.Infof("Start server at: %s", config.Address)
	if err := http.ListenAndServe(config.Address, router.NewDefaultHandler()); err != nil {
		logger.Errorf("http ListenAndServe err: %v", err)
	}
}
