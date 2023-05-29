package linter_plugin_security

import (
	"context"

	"github.com/kubeshop/tracetest/server/linter/plugins/security/rules"
	"github.com/kubeshop/tracetest/server/model"
)

type SecurityPlugin struct {
	model.BasePlugin
}

var (
	_ model.Plugin = &SecurityPlugin{}
)

func NewSecurityPlugin() model.Plugin {
	return SecurityPlugin{
		BasePlugin: model.BasePlugin{
			Name:        "Security",
			Description: "Enforce security for spans and attributes",
			Rules: []model.Rule{
				rules.NewEnforceHttpsProtocolRule(),
				rules.NewEnsuresNoApiKeyLeakRule(),
			},
		},
	}
}

func (p SecurityPlugin) Name() string {
	return "security"
}

func (p SecurityPlugin) Execute(ctx context.Context, trace model.Trace) (model.PluginResult, error) {
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
