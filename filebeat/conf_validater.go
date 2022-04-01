package filebeat

import (
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
	"github.com/dbrainhub/dbrainhub/model/dao"
)

type (
	ConfValidator interface {
		ValidateFilebeatConf(template string) error
		ValidateModuleConf(template string) error
	}
)

func NewConfValidator() ConfValidator {
	return &confValidator{}
}

type confValidator struct {
}

func (c *confValidator) ValidateFilebeatConf(template string) error {
	// yaml parser
	conf, err := dao.NewFileBeatConfFactory().NewFilebeatConf(template)
	if err != nil {
		return err
	}

	// can configure output?
	if err := NewConfGenerator("", nil).CanGenerateFilebeatConf(conf); err != nil {
		return err
	}

	if !conf.HttpInfo.Enabled {
		return errors.FileBeatConfError("http conf is disabled")
	}
	return nil
}

func (c *confValidator) ValidateModuleConf(template string) error {
	// yaml parser
	conf, err := dao.NewFileBeatConfFactory().NewSlowLogModuleConf(template, model.InputModuleType)
	if err != nil {
		return err
	}

	// can configure input?
	if err := NewConfGenerator("", nil).CanGenerateModuleConf(conf); err != nil {
		return err
	}
	return nil
}
