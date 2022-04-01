package model

type (
	FilebeatConf struct {
		// 考虑配置本身复杂性，直接使用 filebeat 的配置结构，这个不是全量配置，但已包含大多数配置。
		Inputs        []*FilebeatInput `config:"filebeat.inputs"`
		HttpInfo      HttpInfo         `config:"http"`
		RainhubOutput DBRainhubOutput  `config:"output.dbrainhub"`
	}

	SlowLogModuleConf struct {
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

	FilebeatInput struct {
		Type    string   `config:"type"`
		Enabled bool     `config:"enabled"`
		Paths   []string `config:"paths"`
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
