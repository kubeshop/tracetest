package lintern_plugin_standards

import (
	"context"

	lintern_plugins_standards_rules "github.com/kubeshop/tracetest/server/lintern/plugins/standards/rules"
	"github.com/kubeshop/tracetest/server/model"
)

type StandardsPlugin struct {
	model.BasePlugin
}

var (
	_ model.Plugin = &StandardsPlugin{}

	DefaultRules = []model.Rule{
		lintern_plugins_standards_rules.NewRequiredAttributesRule(lintern_plugins_standards_rules.DefaultAttrMap),
	}
)

func NewStandardsPlugin(rules ...model.Rule) model.Plugin {
	return StandardsPlugin{
		BasePlugin: model.BasePlugin{
			Name:        "Standards",
			Description: "Standards plugin",
			Rules:       rules,
		},
	}
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
