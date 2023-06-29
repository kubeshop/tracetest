package analyzer

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
)

const (
	ResourceName       = "Analyzer"
	ResourceNamePlural = "Analyzers"
)

var Operations = []resourcemanager.Operation{
	resourcemanager.OperationGet,
	resourcemanager.OperationList,
	resourcemanager.OperationUpdate,
}

type (
	Linter struct {
		ID           id.ID          `json:"id"`
		Name         string         `json:"name"`
		Enabled      bool           `json:"enabled"`
		MinimumScore int            `json:"minimumScore"`
		Plugins      []LinterPlugin `json:"plugins"`
	}

	LinterPlugin struct {
		Slug    string       `json:"slug"`
		Name    string       `json:"name"`
		Enabled bool         `json:"enabled"`
		Rules   []LinterRule `json:"rules"`
	}

	LinterRule struct {
		Slug       string `json:"slug"`
		Weight     int    `json:"weight"`
		ErrorLevel string `json:"errorLevel"`
	}
)

// var (
// 	pluginMap        = plugins_registry.DefaultPluginRegistry.PluginMap()
// 	availablePlugins = strings.Join(plugins_registry.DefaultPluginRegistry.List(), "|")
// )

func (l Linter) Validate() error {
	if l.Name == "" {
		return fmt.Errorf("linter name cannot be empty")
	}

	// for _, p := range l.Plugins {
	// 	rules, ok := pluginMap[p.Slug]
	// 	if !ok {
	// 		return fmt.Errorf("plugin %s not supported, supported plugins are %s", p.Slug, availablePlugins, "|")
	// 	}

	// 	for _, r := range p.Rules {
	// 		index := slices.IndexFunc(rules, func(rule string) bool { return rule == r.Slug })

	// 		if index == -1 {
	// 			availableRules := strings.Join(rules, "|")
	// 			return fmt.Errorf("rule %s not found for plugin %s, supported rules for plugin are %s", r.Slug, p.Slug, availableRules)
	// 		}
	// 	}
	// }

	return nil
}

func (l Linter) HasID() bool {
	return l.ID != ""
}

func (l Linter) GetID() id.ID {
	return l.ID
}

func (l Linter) EnabledPlugins() []LinterPlugin {
	plugins := make([]LinterPlugin, 0)
	for _, p := range l.Plugins {
		if p.Enabled {
			plugins = append(plugins, p)
		}
	}

	return plugins
}

func (l Linter) ShouldSkip() bool {
	return !l.Enabled
}
