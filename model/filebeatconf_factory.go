package model

import (
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/utils/logger"
	"github.com/elastic/beats/v7/libbeat/common"
)

type (
	FileBeatConfFactory interface {
		NewFilebeatConf(confContent string) (*FilebeatConf, error)
		NewModuleConf(confContent string, moduleName string) (*ModuleConf, error)
	}
)

func NewFileBeatConfFactory() FileBeatConfFactory {
	return &filebeatConfFactory{}
}

type filebeatConfFactory struct{}

func (f *filebeatConfFactory) NewFilebeatConf(confContent string) (*FilebeatConf, error) {
	var res FilebeatConf
	if err := f.parse(confContent, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (f *filebeatConfFactory) NewModuleConf(confContent string, moduleName string) (*ModuleConf, error) {
	var confs []*ModuleConf
	if err := f.parse(confContent, &confs); err != nil {
		return nil, err
	}

	for _, conf := range confs {
		if conf.Module == moduleName {
			return conf, nil
		}
	}
	return nil, errors.FileBeatConfError("there is no slowlog conf in config file")
}

func (f *filebeatConfFactory) parse(confContent string, conf interface{}) error {
	cfg, err := common.NewConfigFrom(confContent)
	if err != nil {
		logger.Infof("yaml parse err: %v", err)
		return errors.FileBeatConfError("yaml parse error")
	}
	err = cfg.Unpack(conf)
	if err != nil {
		logger.Infof("config unpack err: %v", err)
		return errors.FileBeatConfError("cfg unpack error:%v", err)
	}

	return nil
}
