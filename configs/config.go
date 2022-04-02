package configs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type GlobalConfig struct {
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
	err = toml.Unmarshal(bytes, &config)
	return &config, err
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
