package linter

import (
	"context"
	"fmt"

	analyzer "github.com/kubeshop/tracetest/server/linter/analyzer"
	linter_plugin_common "github.com/kubeshop/tracetest/server/linter/plugins/common"
	linter_plugin_security "github.com/kubeshop/tracetest/server/linter/plugins/security"
	linter_plugin_standards "github.com/kubeshop/tracetest/server/linter/plugins/standards"
	"github.com/kubeshop/tracetest/server/model"
)

var (
	AvailablePlugins = []model.Plugin{
		linter_plugin_standards.NewStandardsPlugin(),
		linter_plugin_security.NewSecurityPlugin(),
		linter_plugin_common.NewCommonPlugin(),
	}
)

type Linter interface {
	Run(context.Context, model.Trace) (model.LinterResult, error)
	ShouldSkip() (bool, string)
	IsValid() error
}

type linter struct {
	plugins        []model.Plugin
	linterResource analyzer.Linter
}

func NewLinter(linterResource analyzer.Linter, plugins ...model.Plugin) Linter {
	return linter{plugins, linterResource}
}

var _ Linter = &linter{}

func (l linter) Run(ctx context.Context, trace model.Trace) (model.LinterResult, error) {
	pluginResults := make([]model.PluginResult, len(l.plugins))

	totalScore := 0
	passed := true
	for i, plugin := range l.plugins {
		pluginResult, err := plugin.Execute(ctx, trace)
		if err != nil {
			return model.LinterResult{}, err
		}

		passed = passed && pluginResult.Passed
		totalScore += pluginResult.Score
		pluginResults[i] = pluginResult
	}

	return model.LinterResult{
		Plugins: pluginResults,
		Score:   totalScore / len(l.plugins),
		Passed:  passed,
	}, nil
}

func (l linter) ShouldSkip() (bool, string) {
	if !l.linterResource.Enabled {
		return true, "linter is disabled"
	}

	return false, ""
}

func (l linter) IsValid() error {
	plugins, err := l.getPlugins()
	if err != nil {
		return err
	}

	if len(plugins) == 0 {
		return fmt.Errorf("No plugins found")
	}

	return nil
}

func (l linter) getPlugins() ([]model.Plugin, error) {
	plugins := []model.Plugin{}

	for _, cfgPlugin := range l.linterResource.Plugins {
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

func (l linter) findPlugin(pluginName string) (*model.Plugin, error) {
	for _, plugin := range AvailablePlugins {
		if plugin.Name() == pluginName {
			return &plugin, nil
		}
	}

	return nil, fmt.Errorf("plugin %s is not configured", pluginName)
}
