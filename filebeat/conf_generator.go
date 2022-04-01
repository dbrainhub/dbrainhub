package filebeat

import (
	"fmt"
	"github.com/dbrainhub/dbrainhub/model"
	"github.com/dbrainhub/dbrainhub/utils/logger"
	"strings"

	"github.com/dbrainhub/dbrainhub/errors"
)

type (
	ConfGenerator interface {
		// 用于辅助校验配置文件有效性
		CanGenerateFilebeatConf(conf *model.FilebeatConf) error
		CanGenerateModuleConf(conf *model.SlowLogModuleConf) error

		GenerateFilebeatConf(template string) string
		GenerateModuleConf(template string) string
	}
)

func NewConfGenerator(inputPath string, outputHosts []string) ConfGenerator {
	return &confGenImpl{
		inputPath:   inputPath,
		outputHosts: outputHosts,
	}
}

const (
	PlaceHolderForInputPath   = "$input_path"
	PlaceHolderForOutputHosts = "$output_hosts"
)

type confGenImpl struct {
	inputPath   string
	outputHosts []string
}

// check output
func (c *confGenImpl) CanGenerateFilebeatConf(conf *model.FilebeatConf) error {
	for _, host := range conf.RainhubOutput.Hosts {
		if strings.TrimSpace(host) == PlaceHolderForOutputHosts {
			return nil
		}
	}
	logger.Infof("yml cannot configure output hosts, there is no '%s' placeholder", PlaceHolderForOutputHosts)
	return errors.FileBeatConfError("cannot configure output hosts")
}

// check input
func (c *confGenImpl) CanGenerateModuleConf(conf *model.SlowLogModuleConf) error {
	if !conf.SlowLog.Enabled {
		return errors.FileBeatConfError("slowlog is disabled")
	}
	for _, path := range conf.SlowLog.Paths {
		if strings.TrimSpace(path) == PlaceHolderForInputPath {
			return nil
		}
	}
	return errors.FileBeatConfError("cannot configure input path")
}

func (c *confGenImpl) GenerateFilebeatConf(template string) string {
	var hosts string
	for _, host := range c.outputHosts {
		hosts += fmt.Sprintf(`,"%s"`, host)
	}
	hosts = strings.TrimPrefix(hosts, ",")

	return strings.ReplaceAll(template, fmt.Sprintf(`"%s"`, PlaceHolderForOutputHosts), hosts)
}

func (c *confGenImpl) GenerateModuleConf(template string) string {
	return strings.ReplaceAll(template, PlaceHolderForInputPath, c.inputPath)
}
