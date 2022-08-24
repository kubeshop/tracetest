package installer

import (
	"errors"
)

func Start() {
	DockerCompose.Install(DefaultUI)
}

type Installer struct {
	preChecks []preChecker
	configs   []configurator
	installer installer
}

func (i Installer) PreCheck(ui UI) {
	for _, pc := range i.preChecks {
		pc(ui)
	}
}

func (i Installer) Configure(ui UI) configuration {
	config := configuration{}
	for _, confFn := range i.configs {
		config = confFn(config, ui)
	}

	return config
}

func (i Installer) Install(ui UI) {
	i.PreCheck(ui)
	conf := i.Configure(ui)
	i.installer(conf, ui)
}

type preChecker func(ui UI)

type configuration map[string]interface{}

var (
	errConfigNotExists = errors.New("configuration key not set")
	errConfigWrongType = errors.New("configuration value is not the correct type")
)

func (c configuration) Bool(key string) (bool, error) {
	v, exists := c[key]
	if !exists {
		return false, errConfigNotExists
	}
	b, ok := v.(bool)
	if !ok {
		return false, errConfigWrongType
	}

	return b, nil
}

func (c configuration) String(key string) (string, error) {
	v, exists := c[key]
	if !exists {
		return "", errConfigNotExists
	}
	b, ok := v.(string)
	if !ok {
		return "", errConfigWrongType
	}

	return b, nil
}

type configurator func(config configuration, ui UI) configuration
type installer func(config configuration, ui UI)
