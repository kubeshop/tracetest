package linter_plugin_common

import (
	"context"

	"github.com/kubeshop/tracetest/server/linter/plugins/common/rules"
	"github.com/kubeshop/tracetest/server/model"
)

type CommonPlugin struct {
	model.BasePlugin
}

var (
	_ model.Plugin = &CommonPlugin{}
)

func NewCommonPlugin() model.Plugin {
	return CommonPlugin{
		BasePlugin: model.BasePlugin{
			Name:        "Common Issues",
			Description: "Helps you find common problems with your application",
			Rules: []model.Rule{
				rules.NewEnforceDnsUsageRule(),
			},
		},
	}
}

func (p CommonPlugin) Name() string {
	return "common"
}

func (p CommonPlugin) Execute(ctx context.Context, trace model.Trace) (model.PluginResult, error) {
	results := make([]model.RuleResult, len(p.Rules))
	for i, rule := range p.Rules {
		result, err := rule.Evaluate(ctx, trace)
		if err != nil {
			return model.PluginResult{}, err
		}

		results[i] = result
	}

	return model.PluginResult{
		BasePlugin: p.BasePlugin,
		Rules:      results,
	}.CalculateResults(), nil
}
