package filebeat

import (
	"testing"

	"github.com/dbrainhub/dbrainhub/model"
	"github.com/stretchr/testify/assert"
)

func TestConfGenImpl_CanGenerateFilebeatConf_Error(t *testing.T) {
	var confContent = `
output.dbrainhub:
  hosts: ["127.0.0.1:10010",""]
  batch_size: 20480
  retry_limit: 5
  timeout: 2`
	conf, err := model.NewFileBeatConfFactory().NewFilebeatConf(confContent)
	assert.Nil(t, err)
	assert.Error(t, NewConfGenerator("", nil).CanGenerateFilebeatConf(conf))
}

func TestConfGenImpl_CanGenerateFilebeatConf_True(t *testing.T) {
	var confContent = `
output.dbrainhub:
  hosts: ["127.0.0.1:10010","$output_hosts"]
  batch_size: 20480
  retry_limit: 5
  timeout: 2`
	conf, err := model.NewFileBeatConfFactory().NewFilebeatConf(confContent)
	assert.Nil(t, err)
	assert.Nil(t, NewConfGenerator("", nil).CanGenerateFilebeatConf(conf))
}

func TestConfGenImpl_CanGenerateModuleConf_Error(t *testing.T) {
	var confContent = `
- module: mysql
  slowlog:
    enabled: true
    var.paths: ["/path/to/log/mysql/mysql-slow.log*"]`
	conf, err := model.NewFileBeatConfFactory().NewSlowLogModuleConf(confContent, model.InputModuleType)
	assert.Nil(t, err)
	assert.Error(t, NewConfGenerator("", nil).CanGenerateModuleConf(conf))
}

func TestConfGenImpl_CanGenerateModuleConf_True(t *testing.T) {
	var confContent = `
- module: mysql
  slowlog:
    enabled: true
    var.paths: ["$input_path"]`
	conf, err := model.NewFileBeatConfFactory().NewSlowLogModuleConf(confContent, model.InputModuleType)
	assert.Nil(t, err)
	assert.Nil(t, NewConfGenerator("", nil).CanGenerateModuleConf(conf))
}

func TestConfGenImpl_GenerateModuleConf(t *testing.T) {
	var confContent = `
- module: mysql
  slowlog:
    enabled: true
    var.paths: ["$input_path"]`

	gen := NewConfGenerator("123", nil)
	assert.Equal(t, gen.GenerateModuleConf(confContent), `
- module: mysql
  slowlog:
    enabled: true
    var.paths: ["123"]`)
}

func TestConfGenImpl_GenerateFilebeatConf(t *testing.T) {
	var confContent = `
output.dbrainhub:
  hosts: ["127.0.0.1:10010","$output_hosts"]
  batch_size: 20480
  retry_limit: 5
  timeout: 2`
	gen := NewConfGenerator("", []string{"123", "456"})
	assert.Equal(t, gen.GenerateFilebeatConf(confContent), `
output.dbrainhub:
  hosts: ["127.0.0.1:10010","123","456"]
  batch_size: 20480
  retry_limit: 5
  timeout: 2`)
}
