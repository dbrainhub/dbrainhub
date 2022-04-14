package configs

import (
	"fmt"

	"github.com/dbrainhub/dbrainhub/api"
)

type AgentConfig struct {
	DB struct {
		HostType string `toml:"host_type"`
		DBType   string `toml:"type"`
		Port     int    `toml:"port"`
		User     string `toml:"user"`
		Passwd   string `toml:"password"`
	} `toml:"db"`

	Server struct {
		Addr                string `toml:"addr"` // TODO: 扩展成多个
		TimeoutMs           int    `toml:"http_timeout_ms"`
		Retry               int    `toml:"http_retry"`
		RetryIntervalMs     int    `toml:"http_retry_interval_ms"`
		HeartbeatIntervalMs int    `toml:"heartbeat_interval_ms"`
	} `toml:"server"`

	LogInfo struct {
		Level  string `toml:"level"`
		LogDir string `toml:"log_dir"`
		Name   string `toml:"name"`
	} `toml:"log_info"`

	Filebeat struct {
		FilebeatConfTemplate         string `toml:"filebeat_conf_template"`
		ModuleConfTemplate           string `toml:"module_conf_template"`
		HomePath                     string `toml:"home_path"`
		AliveListenerIntervalMs      int    `toml:"alive_listener_interval_ms"`
		AliveListenerTimeoutMs       int    `toml:"alive_listener_http_timeout_ms"`
		AliveListenerHttpRetry       int    `toml:"alive_listener_http_retry"`
		AliveListenerRetryIntervalMs int    `toml:"alive_listener_http_retry_interval_ms"`
		SlowlogListenerIntervalMs    int    `toml:"slowlog_listener_interval_ms"`
		StartupTimeoutMs             int    `toml:"startup_timeout_ms"`
	}
}

var globalAgentConfig AgentConfig

func GetGlobalAgentConfig() *AgentConfig {
	return &globalAgentConfig
}

func (a *AgentConfig) ConvertHostType() api.StartupReportRequest_HostType {
	switch a.DB.HostType {
	case "self":
		return api.StartupReportRequest_SELF
	case "aliyun":
		return api.StartupReportRequest_ALIYUN
	case "tencentyun":
		return api.StartupReportRequest_TENCENTYUN
	default:
		return api.StartupReportRequest_UNKNOWN
	}
}

func (a *AgentConfig) ConvertDBType() (api.StartupReportRequest_DBType, error) {
	switch a.DB.DBType {
	case "mysql":
		return api.StartupReportRequest_MYSQL, nil
	case "tidb":
		return api.StartupReportRequest_TIDB, nil
	case "redis":
		return api.StartupReportRequest_REDIS, nil
	case "mongodb":
		return api.StartupReportRequest_MONGODB, nil
	}
	return 0, fmt.Errorf("unknown db_type in agent config")
}
