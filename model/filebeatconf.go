package model

type (
	FilebeatConf struct {
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
