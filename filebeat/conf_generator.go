package filebeat

import (
	"fmt"
	"strings"

	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
	"github.com/dbrainhub/dbrainhub/utils/logger"
)

type (
	FilebeatConfGenerator interface {
		// 用于辅助校验配置文件有效性
		CanGenerate(conf *model.FilebeatConf) error
		Generate(template string) string
	}

	ModuleConfGenerator interface {
		// 用于辅助校验配置文件有效性
		CanGenerate(conf *model.ModuleConf) error
		Generate(template string) string
	}
)

func NewFilebeatConfGenerator(outputHosts []string) FilebeatConfGenerator {
	return &filebeatConfGenImpl{
		outputHosts: outputHosts,
	}
}

func NewModuleConfGenerator(inputPath string) ModuleConfGenerator {
	return &moduleConfGenImpl{
		inputPath: inputPath,
	}
}

const (
	PlaceHolderForInputPath   = "$input_path"
	PlaceHolderForOutputHosts = "$output_hosts"
)

type filebeatConfGenImpl struct {
	outputHosts []string
}

// check output
func (c *filebeatConfGenImpl) CanGenerate(conf *model.FilebeatConf) error {
	for _, host := range conf.RainhubOutput.Hosts {
		if strings.TrimSpace(host) == PlaceHolderForOutputHosts {
			return nil
		}
	}
	logger.Infof("yml cannot configure output hosts, there is no '%s' placeholder", PlaceHolderForOutputHosts)
	return errors.FileBeatConfError("cannot configure output hosts")
}

func (c *filebeatConfGenImpl) Generate(template string) string {
	var hosts string
	for _, host := range c.outputHosts {
		hosts += fmt.Sprintf(`,"%s"`, host)
	}
	hosts = strings.TrimPrefix(hosts, ",")

	return strings.ReplaceAll(template, fmt.Sprintf(`"%s"`, PlaceHolderForOutputHosts), hosts)
}

type moduleConfGenImpl struct {
	inputPath string
}

// check input
func (c *moduleConfGenImpl) CanGenerate(conf *model.ModuleConf) error {
	if !conf.SlowLog.Enabled {
		return errors.FileBeatConfError("slowlog is disabled")
	}
	for _, path := range conf.SlowLog.Paths {
		if strings.TrimSpace(path) == PlaceHolderForInputPath {
			return nil
		}
	}
	logger.Infof("yml cannot configure input path, there is no '%s' placeholder", PlaceHolderForInputPath)
	return errors.FileBeatConfError("cannot configure input path")
}

func (c *moduleConfGenImpl) Generate(template string) string {
	return strings.ReplaceAll(template, PlaceHolderForInputPath, c.inputPath)
}
