package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type GlobalConfig struct {
	Address string `json:"address"`
	LogInfo struct {
		Level  string `json:"level"`
		LogDir string `json:"log_dir"`
		Name   string `json:"name"`
	} `json:"log_info"`
}

var globalConfig *GlobalConfig

func GetGlobalConfig() *GlobalConfig {
	return globalConfig
}

func InitConfigOrPanic(path string) {
	configPath := getConfigPath(path)
	config, err := loadConfigFromFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("InitConfig error, err: %v", err))
	}
	globalConfig = config
}

func loadConfigFromFile(path string) (*GlobalConfig, error) {
	log.Printf("loadConfigFromFile: path=%s", path)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config GlobalConfig
	err = json.Unmarshal(bytes, &config)
	return &config, err
}

// 按 用户指定 -> 环境变量 -> 默认值 的方式获取配置文件路径。
func getConfigPath(path string) string {
	const ConfigPathEnv = "config_path"
	const ConfigPath = "config.json"

	if path != "" {
		return path
	}

	path = os.Getenv(ConfigPathEnv)
	if path != "" {
		return path
	}
	return ConfigPath
}
