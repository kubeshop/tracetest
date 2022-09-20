package installer

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/analytics"
)

var (
	Force = false
)

func Start() {
	analytics.Track("Start", "installer", map[string]string{})
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
		{"Using Kubernetes", kubernetes.Install},
	}, 0)

	option.fn(ui)
}

type installer struct {
	name      string
	preChecks []preChecker
	configs   []configurator
	installFn func(config configuration, ui UI)
}

func (i installer) PreCheck(ui UI) {
	ui.Title("Let's check if your system has everything we need")
	for _, pc := range i.preChecks {
		pc(ui)
	}

	ui.Title("Your system is ready! Now, let's configure TraceTest")
}

func (i installer) Configure(ui UI) configuration {
	config := newConfiguration(ui)
	config.set("installer", i.name)
	for _, confFn := range i.configs {
		config = confFn(config, ui)
	}

	return config
}

func (i installer) Install(ui UI) {
	analytics.Track("PreCheck", "installer", map[string]string{})
	i.PreCheck(ui)

	analytics.Track("Configure", "installer", map[string]string{})
	conf := i.Configure(ui)

	ui.Title("Thanks! We are ready to install TraceTest now")

	i.installFn(conf, ui)
}

type preChecker func(ui UI)

type configuration struct {
	db map[string]interface{}
	ui UI
}

func newConfiguration(ui UI) configuration {
	return configuration{
		db: map[string]interface{}{},
		ui: ui,
	}
}

func (c configuration) set(key string, value interface{}) {
	if _, exists := c.db[key]; exists {
		c.ui.Panic(fmt.Errorf("config key %s already exists", key))
	}

	c.db[key] = value
}

func (c configuration) get(key string) interface{} {
	v, exists := c.db[key]
	if !exists {
		c.ui.Panic(fmt.Errorf("config key %s not exists", key))
	}

	return v
}

func (c configuration) Bool(key string) bool {
	b, ok := c.get(key).(bool)
	if !ok {
		c.ui.Panic(fmt.Errorf("config key %s is not a bool", key))
	}

	return b
}

func (c configuration) String(key string) string {
	s, ok := c.get(key).(string)
	if !ok {
		c.ui.Panic(fmt.Errorf("config key %s is not a string", key))
	}

	return s
}

type configurator func(config configuration, ui UI) configuration

func trackInstall(name string, config configuration, extra map[string]string) {
	props := map[string]string{
		"type":                    name,
		"install_backend":         fmt.Sprintf("%t", config.Bool("tracetest.backend.install")),
		"install_collector":       fmt.Sprintf("%t", config.Bool("tracetest.collector.install")),
		"install_demo_pokeshop":   fmt.Sprintf("%t", config.Bool("demo.enable.pokeshop")),
		"install_demo_otel":       fmt.Sprintf("%t", config.Bool("demo.enable.otel")),
		"enable_server_analytics": fmt.Sprintf("%t", config.Bool("tracetest.analytics")),
		"backend_type":            config.String("tracetest.backend.type"),
	}

	for k, v := range extra {
		props[k] = v
	}

	analytics.Track("Apply", "installer", props)
}
