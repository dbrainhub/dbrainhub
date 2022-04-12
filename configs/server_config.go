package configs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	Address string `toml:"address"`
	DB      struct {
		Dialect  string `toml:"dialect"`
		User     string `toml:"user"`
		Password string `toml:"password"`
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		Database string `toml:"database"`
	} `toml:"db"`
	LogInfo struct {
		Level  string `toml:"level"`
		LogDir string `toml:"log_dir"`
		Name   string `toml:"name"`
	} `toml:"log_info"`
	OutputServer struct {
		EsAddresses  []string `toml:"es_addresses"`
		QpsThreshold int      `toml:"qps_threshold"`
	} `toml:"output_server"`
}

var globalServerConfig ServerConfig

func GetGlobalServerConfig() *ServerConfig {
	return &globalServerConfig
}

func InitConfigOrPanic(path string, conf interface{}) {
	configPath := getConfigPath(path)
	err := loadConfigFromFile(configPath, conf)
	if err != nil {
		panic(fmt.Sprintf("InitConfig error, err: %v", err))
	}
}

func loadConfigFromFile(path string, conf interface{}) error {
	log.Printf("loadConfigFromFile: path=%s", path)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return toml.Unmarshal(bytes, conf)
}

// 按 用户指定 -> 环境变量 -> 默认值 的方式获取配置文件路径。
func getConfigPath(path string) string {
	const ConfigPathEnv = "config_path"
	const ConfigPath = "config.yaml"

	if path != "" {
		return path
	}

	path = os.Getenv(ConfigPathEnv)
	if path != "" {
		return path
	}
	return ConfigPath
}
