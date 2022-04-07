package model

type (
	FilebeatConf struct {
		FilebeatModule FilebeatConfigModule `config:"filebeat.config.modules"`
		HttpInfo       HttpInfo             `config:"http"`
		RainhubOutput  DBRainhubOutput      `config:"output.dbrainhub"`
	}

	ModuleConf struct {
		Module  string `config:"module"`
		SlowLog struct {
			Enabled bool     `config:"enabled"`
			Paths   []string `config:"var.paths"`
		} `config:"slowlog"`
	}

	HttpInfo struct {
		Enabled bool   `config:"enabled"`
		Host    string `config:"host"`
		Port    int    `config:"port"`
	}

	// the module config in filebeat.yml
	// ReloadEnabled is the switch for filebeat reload.
	FilebeatConfigModule struct {
		Enabled       bool   `config:"enabled"`
		Path          string `config:"path"`
		ReloadEnabled bool   `config:"reload.enabled"`
		ReloadPeriod  string `config:"reload.period"`
	}

	DBRainhubOutput struct {
		Hosts      []string `config:"hosts"`
		BatchSize  int      `config:"batch_size"`
		RetryLimit int      `config:"retry_limit"`
		Timeout    int      `config:"timeout"`
	}
)

const (
	InputModuleType = "mysql"
)
