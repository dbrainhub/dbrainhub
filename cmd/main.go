package main

import (
	"flag"
	"fmt"

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

	httpRouter := router.NewDefaultRouter(configs.GetGlobalConfig().Address)
	if err := httpRouter.Run(); err != nil {
		logger.Fatal("http server err: %v", err)
	}
}
