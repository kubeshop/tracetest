package plugins

import (
	"context"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/rules"
	"github.com/kubeshop/tracetest/server/model"
)

type BasePlugin struct {
	id           string
	ruleRegistry *rules.RuleRegistry
}

func NewPlugin(ID string, ruleRegistry *rules.RuleRegistry) Plugin {
	return BasePlugin{id: ID, ruleRegistry: ruleRegistry}
}

func (p BasePlugin) ID() string {
	return p.id
}

func (p BasePlugin) RuleRegistry() *rules.RuleRegistry {
	return p.ruleRegistry
}

func (p BasePlugin) Execute(ctx context.Context, trace model.Trace, config analyzer.LinterPlugin) (analyzer.PluginResult, error) {
	res := make([]analyzer.RuleResult, 0, len(config.Rules))

	for _, cfgRule := range config.Rules {
		rule, err := p.ruleRegistry.Get(cfgRule.ID)
		if err != nil {
			return analyzer.PluginResult{}, err
		}

		result, err := rule.Evaluate(ctx, trace, cfgRule)
		if err != nil {
			return analyzer.PluginResult{}, err
		}

		res = append(res, result)
	}

	passed := true
	for _, result := range res {
		if !result.Passed {
			passed = false
		}
	}

	return analyzer.PluginResult{
		//config
		ID:          p.id,
		Name:        config.Name,
		Description: config.Description,

		// results
		Passed: passed,
		Rules:  res,
	}.CalculateResults(), nil
}
