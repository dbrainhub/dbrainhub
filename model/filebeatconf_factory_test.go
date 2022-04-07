package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFilebeatConf = `
output.dbrainhub:
  hosts: ["127.0.0.1:10010","$output_hosts"]
  batch_size: 20480
  retry_limit: 5
  timeout: 2

http.enabled: true
http.host: localhost
http.port: 5066

filebeat.config.modules:
    enabled: true
    path: modules.d/*.yml
    reload.enabled: true
    reload.period: 10s`

func TestNewFilebeatConf(t *testing.T) {
	var factory filebeatConfFactory
	config, err := factory.NewFilebeatConf(testFilebeatConf)
	assert.Nil(t, err)

	assert.Equal(t, config.HttpInfo.Enabled, true)
	assert.Equal(t, config.HttpInfo.Port, 5066)
	assert.Equal(t, config.HttpInfo.Host, "localhost")

	assert.Equal(t, config.RainhubOutput.Hosts[0], "127.0.0.1:10010")
	assert.Equal(t, config.RainhubOutput.BatchSize, 20480)

	assert.Equal(t, config.FilebeatModule.Path, "modules.d/*.yml")
	assert.Equal(t, config.FilebeatModule.ReloadEnabled, true)
	assert.Equal(t, config.FilebeatModule.ReloadPeriod, "10s")

}

var testModuleFileConf = `
# Module: mysql
# Docs: https://www.elastic.co/guide/en/beats/filebeat/master/filebeat-module-mysql.html

- module: mysql
  # Slow logs
  slowlog:
    enabled: true

    # Set custom paths for the log files. If left empty,
    # Filebeat will choose the paths depending on your OS.
    var.paths: ["/path/to/log/mysql/mysql-slow.log*"]

  # Error logs
  error:
    enabled: false

    # Set custom paths for the log files. If left empty,
    # Filebeat will choose the paths depending on your OS.
    # var.paths:`

func TestGetModuleConf(t *testing.T) {
	var factory filebeatConfFactory
	conf, err := factory.NewModuleConf(testModuleFileConf, InputModuleType)
	assert.Nil(t, err)

	assert.Equal(t, conf.Module, "mysql")
	assert.Equal(t, conf.SlowLog.Enabled, true)
	assert.Equal(t, conf.SlowLog.Paths[0], "/path/to/log/mysql/mysql-slow.log*")

}
