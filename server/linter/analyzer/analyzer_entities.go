package analyzer

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/linter/metadata"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"golang.org/x/exp/slices"
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

func (l Linter) Validate() error {
	if l.Name == "" {
		return fmt.Errorf("linter name cannot be empty")
	}

	for _, p := range l.Plugins {
		plugin, ok := findPlugin(p.Slug, metadata.Plugins)
		if !ok {
			return fmt.Errorf("plugin %s not supported, supported plugins are %s", p.Slug, metadata.AvailablePlugins)
		}

		for _, r := range p.Rules {
			index := slices.IndexFunc(plugin.Rules, func(rule metadata.RuleMetadata) bool { return rule.Slug == r.Slug })

			if index == -1 {
				availableRules := strings.Join(getAvailableRules(plugin), "|")
				return fmt.Errorf("rule %s not found for plugin %s, supported rules for plugin are %s", r.Slug, p.Slug, availableRules)
			}
		}
	}

	return nil
}

func findPlugin(slug string, plugins []metadata.PluginMetadata) (metadata.PluginMetadata, bool) {
	for _, p := range plugins {
		if p.Slug == slug {
			return p, true
		}
	}

	return metadata.PluginMetadata{}, false
}

func getAvailableRules(plugin metadata.PluginMetadata) []string {
	rules := make([]string, 0)

	for _, r := range plugin.Rules {
		rules = append(rules, r.Slug)
	}

	return rules
}

func getDefaultLinter() Linter {
	plugins := []LinterPlugin{}

	for _, plugin := range metadata.Plugins {
		rules := []LinterRule{}

		for _, rule := range plugin.Rules {
			rules = append(rules, LinterRule{
				Slug:       rule.Slug,
				Weight:     rule.DefaultWeight,
				ErrorLevel: rule.DefaultErrorLevel,
			})
		}

		plugins = append(plugins, LinterPlugin{
			Slug:    plugin.Slug,
			Name:    plugin.Name,
			Enabled: plugin.DefaultEnabled,
			Rules:   rules,
		})
	}

	return Linter{
		ID:           id.ID("current"),
		Name:         "analyzer",
		Enabled:      true,
		Plugins:      plugins,
		MinimumScore: 0,
	}
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
