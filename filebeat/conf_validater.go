package filebeat

import (
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
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
	conf, err := model.NewFileBeatConfFactory().NewFilebeatConf(template)
	if err != nil {
		return err
	}

	// can configure output?
	if err := NewFilebeatConfGenerator("", 0, nil).CanGenerate(conf); err != nil {
		return err
	}

	if !conf.HttpInfo.Enabled {
		return errors.FileBeatConfError("http conf is disabled")
	}

	if !conf.FilebeatModule.Enabled || !conf.FilebeatModule.ReloadEnabled {
		return errors.FileBeatConfError("filebeat module reload.enabled is disabled")
	}
	return nil
}

func (c *confValidator) ValidateModuleConf(template string) error {
	// yaml parser
	conf, err := model.NewFileBeatConfFactory().NewModuleConf(template, model.InputModuleType)
	if err != nil {
		return err
	}

	// can configure input?
	if err := NewModuleConfGenerator("").CanGenerate(conf); err != nil {
		return err
	}
	return nil
}
