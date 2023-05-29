package linter_plugin_standards

import (
	"context"

	"github.com/kubeshop/tracetest/server/linter/plugins/standards/rules"
	"github.com/kubeshop/tracetest/server/model"
)

type StandardsPlugin struct {
	model.BasePlugin
}

var (
	_ model.Plugin = &StandardsPlugin{}
)

func NewStandardsPlugin() model.Plugin {
	return StandardsPlugin{
		BasePlugin: model.BasePlugin{
			Name:        "Global Standards",
			Description: "Enforce standards for spans and attributes",
			Rules: []model.Rule{
				rules.NewEnsureSpanNamingRule(),
				rules.NewRequiredAttributesRule(),
				rules.NewEnsureAttributeNamingRule(),
				rules.NewNotEmptyAttributesRule(),
			},
		},
	}
}

func (p StandardsPlugin) Name() string {
	return "standards"
}

func (p StandardsPlugin) Execute(ctx context.Context, trace model.Trace) (model.PluginResult, error) {
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
