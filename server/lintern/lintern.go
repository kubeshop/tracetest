package lintern

import (
	"context"

	lintern_plugin_standards "github.com/kubeshop/tracetest/server/lintern/plugins/standards"
	"github.com/kubeshop/tracetest/server/model"
)

var (
	DefaultPlugins = []model.Plugin{
		lintern_plugin_standards.NewStandardsPlugin(lintern_plugin_standards.DefaultRules...),
	}
)

type Lintern interface {
	Run(context.Context, model.Trace) (model.LinternResult, error)
}

type lintern struct {
	plugins []model.Plugin
}

func NewLintern(plugins ...model.Plugin) Lintern {
	return lintern{
		plugins: plugins,
	}
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
