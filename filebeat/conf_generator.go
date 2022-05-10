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

func NewFilebeatConfGenerator(localip string, port int, outputHosts []string) FilebeatConfGenerator {
	return &filebeatConfGenImpl{
		localip:     localip,
		port:        port,
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
	PlaceHolderForLocalIP     = "$localip"
	PlaceHolderForPort        = "$port"
)

type filebeatConfGenImpl struct {
	localip     string
	port        int
	outputHosts []string
}

// check output
func (c *filebeatConfGenImpl) CanGenerate(conf *model.FilebeatConf) error {
	if !c.hasHostPlaceholder(conf) {
		logger.Error("yml cannot configure output.dbrainhub.hosts, there is no '%s' placeholder", PlaceHolderForOutputHosts)
		return errors.FileBeatConfError("cannot configure output hosts")
	}

	if conf.RainhubOutput.DBIP != PlaceHolderForLocalIP {
		logger.Error("yml cannot configure output.dbrainhub.db_ip, there is no '%s' placeholder", PlaceHolderForLocalIP)
		return errors.FileBeatConfError("cannot configure output localip")
	}

	if conf.RainhubOutput.DBPort != PlaceHolderForPort {
		logger.Error("yml cannot configure output.dbrainhub.db_port, there is no '%s' placeholder", PlaceHolderForPort)
		return errors.FileBeatConfError("cannot configure output port")
	}
	return nil
}

func (c *filebeatConfGenImpl) hasHostPlaceholder(conf *model.FilebeatConf) bool {
	for _, host := range conf.RainhubOutput.Hosts {
		if strings.TrimSpace(host) == PlaceHolderForOutputHosts {
			return true
		}
	}
	return false
}

func (c *filebeatConfGenImpl) Generate(template string) string {
	// for hosts
	var hosts string
	for _, host := range c.outputHosts {
		hosts += fmt.Sprintf(`,"%s"`, host)
	}
	hosts = strings.TrimPrefix(hosts, ",")
	template = strings.ReplaceAll(template, fmt.Sprintf(`"%s"`, PlaceHolderForOutputHosts), hosts)

	// for localip
	template = strings.ReplaceAll(template, PlaceHolderForLocalIP, c.localip)

	// for port
	template = strings.ReplaceAll(template, PlaceHolderForPort, fmt.Sprintf("%d", c.port))
	return template
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
