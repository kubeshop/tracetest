package installer

import (
	"fmt"

	cliUI "github.com/kubeshop/tracetest/cli/ui"
)

type configuration struct {
	db map[string]interface{}
	ui cliUI.UI
}

func newConfiguration(ui cliUI.UI) configuration {
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
	return c.db[key]
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
		return ""
	}

	return s
}

type configurator func(config configuration, ui cliUI.UI) configuration
