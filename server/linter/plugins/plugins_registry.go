package plugins

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/rules"
	"github.com/kubeshop/tracetest/server/model"
)

type Plugin interface {
	Execute(context.Context, model.Trace, analyzer.LinterPlugin) (analyzer.PluginResult, error)
	Slug() string
	RuleRegistry() rules.RuleRegistry
}

type Registry struct {
	plugins map[string]Plugin
}

func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]Plugin),
	}
}

func (r *Registry) Register(p Plugin) *Registry {
	r.plugins[p.Slug()] = p
	return r
}

func (r *Registry) Get(resourceName string) (Plugin, error) {
	plugin, ok := r.plugins[resourceName]
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", resourceName)
	}

	return plugin, nil
}

func (r *Registry) PluginMap() map[string][]string {
	plugins := make(map[string][]string)
	for k, plugin := range r.plugins {
		plugins[k] = plugin.RuleRegistry().List()
	}

	return plugins
}

func (r *Registry) List() []string {
	var plugins []string
	for k := range r.plugins {
		plugins = append(plugins, k)
	}

	return plugins
}
