package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/dbrainhub/dbrainhub/router"
	"github.com/dbrainhub/dbrainhub/server"
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
	if err := printConfig(config); err != nil {
		logger.Errorf("print global config error, err: %v, exit...", err)
		os.Exit(1)
	}

	logger.InitLog(config.LogInfo.LogDir, config.LogInfo.Name, config.LogInfo.Level)

	if err := server.InitDefaultEsClientAsync(config); err != nil {
		logger.Errorf("InitDefaultEsClientSync err: %v, exit...", err)
		os.Exit(1)
	}

	logger.Infof("Start server at: %s", config.Server.Address)
	if err := http.ListenAndServe(config.Server.Address, router.NewDefaultHandler()); err != nil {
		logger.Errorf("http ListenAndServe err: %v", err)
	}
}

func printConfig(config *configs.ServerConfig) error {
	var out bytes.Buffer
	if err := toml.NewEncoder(&out).Encode(config); err != nil {
		return err
	}

	fmt.Printf("global server config: \n%v\n", out.String())
	return nil
}
