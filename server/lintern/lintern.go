package lintern

import (
	"context"
	"fmt"

	lintern_plugin_common "github.com/kubeshop/tracetest/server/lintern/plugins/common"
	lintern_plugin_security "github.com/kubeshop/tracetest/server/lintern/plugins/security"
	lintern_plugin_standards "github.com/kubeshop/tracetest/server/lintern/plugins/standards"
	lintern_resource "github.com/kubeshop/tracetest/server/lintern/resource"
	"github.com/kubeshop/tracetest/server/model"
)

var (
	AvailablePlugins = []model.Plugin{
		lintern_plugin_standards.NewStandardsPlugin(),
		lintern_plugin_security.NewSecurityPlugin(),
		lintern_plugin_common.NewCommonPlugin(),
	}
)

type Lintern interface {
	Run(context.Context, model.Trace) (model.LinternResult, error)
	ShouldSkip() (bool, string)
	IsValid() error
}

type lintern struct {
	plugins         []model.Plugin
	linternResource lintern_resource.Lintern
}

func NewLintern(linternResource lintern_resource.Lintern, plugins ...model.Plugin) Lintern {
	return lintern{plugins, linternResource}
}

var _ Lintern = &lintern{}

func (l lintern) Run(ctx context.Context, trace model.Trace) (model.LinternResult, error) {
	pluginResults := make([]model.PluginResult, len(l.plugins))

	totalScore := 0
	passed := true
	for i, plugin := range l.plugins {
		pluginResult, err := plugin.Execute(ctx, trace)
		if err != nil {
			return model.LinternResult{}, err
		}

		passed = passed && pluginResult.Passed
		totalScore += pluginResult.Score
		pluginResults[i] = pluginResult
	}

	return model.LinternResult{
		Plugins: pluginResults,
		Score:   totalScore / len(l.plugins),
		Passed:  passed,
	}, nil
}

func (l lintern) ShouldSkip() (bool, string) {
	if !l.linternResource.Enabled {
		return true, "Lintern is disabled"
	}

	return false, ""
}

func (l lintern) IsValid() error {
	plugins, err := l.getPlugins()
	if err != nil {
		return err
	}

	if len(plugins) == 0 {
		return fmt.Errorf("No plugins found")
	}

	return nil
}

func (l lintern) getPlugins() ([]model.Plugin, error) {
	plugins := []model.Plugin{}

	for _, cfgPlugin := range l.linternResource.Plugins {
		if cfgPlugin.Enabled {
			plugin, err := l.findPlugin(cfgPlugin.Name)
			if err != nil {
				return nil, err
			}

			plugins = append(plugins, *plugin)
		}
	}

	return plugins, nil
}

func (l lintern) findPlugin(pluginName string) (*model.Plugin, error) {
	for _, plugin := range AvailablePlugins {
		if plugin.Name() == pluginName {
			return &plugin, nil
		}
	}

	return nil, fmt.Errorf("plugin %s is not configured", pluginName)
}
