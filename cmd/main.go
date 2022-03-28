package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/dbrainhub/dbrainhub/router"
	"github.com/dbrainhub/dbrainhub/utils/logger"
)

var configFilePath = flag.String("config", "", "config file path")

func main() {
	fmt.Printf("Hello World! \n")

	flag.Parse()
	configs.InitConfigOrPanic(*configFilePath)

	config := configs.GetGlobalConfig()
	logger.InitLog(config.LogInfo.LogDir, config.LogInfo.Name, config.LogInfo.Level)

	fmt.Printf("[INFO] Start server at: %s \n", config.Address)

	if err := http.ListenAndServe(config.Address, router.NewDefaultHandler()); err != nil {
		logger.Errorf("http ListenAndServe err: %v", err)
	}
}
