package installer

import (
	"errors"
)

func Start() {
	ui := DefaultUI

	ui.Banner()

	ui.Println(`
Hi! Welcome to the TraceTest server installer. I'll help you set up your TraceTest server by asking you a few questions
and configuring your system with all the requirements, so you can start TraceTesting right away!

To get more info about TraceTest, you can check our docs at https://kubeshop.github.io/tracetest/

If you have any issues, please let us know by creating an issue (https://github.com/kubeshop/tracetest/issues/new/choose)
or reach us on Discord https://discord.gg/6zupCZFQbe

`)

	option := ui.Select("How do you want to run TraceTest?", []option{
		{"Using Docker Compose", dockerCompose.Install},
	})

	option.fn(ui)
}

type installer struct {
	preChecks []preChecker
	configs   []configurator
	installFn func(config configuration, ui UI)
}

func (i installer) PreCheck(ui UI) {
	for _, pc := range i.preChecks {
		pc(ui)
	}
}

func (i installer) Configure(ui UI) configuration {
	config := configuration{}
	for _, confFn := range i.configs {
		config = confFn(config, ui)
	}

	return config
}

func (i installer) Install(ui UI) {
	i.PreCheck(ui)
	conf := i.Configure(ui)
	i.installFn(conf, ui)
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
