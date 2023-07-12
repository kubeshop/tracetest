package analyzer

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"golang.org/x/exp/slices"
)

const (
	ResourceName       = "Analyzer"
	ResourceNamePlural = "Analyzers"
)

type (
	Linter struct {
		ID           id.ID          `json:"id"`
		Name         string         `json:"name"`
		Enabled      bool           `json:"enabled"`
		MinimumScore int            `json:"minimumScore"`
		Plugins      []LinterPlugin `json:"plugins"`
	}

	LinterPlugin struct {
		ID      string       `json:"id"`
		Enabled bool         `json:"enabled"`
		Rules   []LinterRule `json:"rules"`

		// internal fields
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	LinterRule struct {
		ID         string `json:"id"`
		Weight     int    `json:"weight"`
		ErrorLevel string `json:"errorLevel"`

		// internal fields
		Name             string   `json:"name"`
		ErrorDescription string   `json:"errorDescription"`
		Description      string   `json:"description"`
		Tips             []string `json:"tips"`
	}
)

func (l Linter) Validate() error {
	if l.Name == "" {
		return fmt.Errorf("linter name cannot be empty")
	}

	for _, p := range l.Plugins {
		plugin, ok := findPlugin(p.ID, DefaultPlugins)
		availableRules := strings.Join(getAvailableRules(plugin), " | ")

		if !ok {
			availablePlugins := strings.Join(AvailablePlugins, " | ")
			return fmt.Errorf("plugin %s not supported, supported plugins are %s", p.ID, availablePlugins)
		}

		if len(p.Rules) != len(plugin.Rules) {
			return fmt.Errorf("plugin %s requires %d rules, but %d provided, supported rules for plugin are %s", p.ID, len(plugin.Rules), len(p.Rules), availableRules)
		}

		for _, r := range p.Rules {
			index := slices.IndexFunc(plugin.Rules, func(rule LinterRule) bool { return rule.ID == r.ID })

			if index == -1 {
				return fmt.Errorf("rule %s not found for plugin %s, supported rules for plugin are %s", r.ID, p.ID, availableRules)
			}
		}
	}

	return nil
}

func findPlugin(ID string, plugins []LinterPlugin) (LinterPlugin, bool) {
	for _, p := range plugins {
		if p.ID == ID {
			return p, true
		}
	}

	return LinterPlugin{}, false
}

func findRule(ID string, rules []LinterRule) (LinterRule, bool) {
	for _, r := range rules {
		if r.ID == ID {
			return r, true
		}
	}

	return LinterRule{}, false
}

func getAvailableRules(plugin LinterPlugin) []string {
	rules := make([]string, 0)

	for _, r := range plugin.Rules {
		rules = append(rules, r.ID)
	}

	return rules
}

func GetDefaultLinter() Linter {
	return Linter{
		ID:           id.ID("current"),
		Name:         "analyzer",
		Enabled:      true,
		Plugins:      DefaultPlugins,
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

func (l Linter) WithMetadata() (Linter, error) {
	plugins := make([]LinterPlugin, 0)

	for _, p := range l.Plugins {
		metadataPlugin, ok := findPlugin(p.ID, DefaultPlugins)
		if !ok {
			return l, fmt.Errorf("plugin %s not supported, supported plugins are %s", p.ID, strings.Join(AvailablePlugins, " | "))
		}

		rules := make([]LinterRule, 0)
		for _, r := range p.Rules {
			metadataRule, ok := findRule(r.ID, metadataPlugin.Rules)
			if !ok {
				return l, fmt.Errorf("rule %s not found for plugin %s, supported rules for plugin are %s", r.ID, p.ID, strings.Join(getAvailableRules(metadataPlugin), " | "))
			}

			rules = append(rules, LinterRule{
				// config
				ID:         r.ID,
				ErrorLevel: r.ErrorLevel,
				Weight:     r.Weight,

				// metadata
				Description:      metadataRule.Description,
				Name:             metadataRule.Name,
				ErrorDescription: metadataRule.ErrorDescription,
				Tips:             metadataRule.Tips,
			})
		}

		plugins = append(plugins, LinterPlugin{
			Rules: rules,

			// config
			ID:      p.ID,
			Enabled: p.Enabled,

			// metadata
			Name:        metadataPlugin.Name,
			Description: metadataPlugin.Description,
		})
	}

	l.Plugins = plugins
	return l, nil
}
