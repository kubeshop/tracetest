package plugins

import (
	"context"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/rules"
	"github.com/kubeshop/tracetest/server/model"
)

type BasePlugin struct {
	slug         string
	ruleRegistry rules.RuleRegistry
}

func NewPlugin(slug string, ruleRegistry rules.RuleRegistry) Plugin {
	return BasePlugin{

		ruleRegistry: ruleRegistry,
	}
}

func (p BasePlugin) Slug() string {
	return p.slug
}

func (p BasePlugin) RuleRegistry() rules.RuleRegistry {
	return p.ruleRegistry
}

func (p BasePlugin) Execute(ctx context.Context, trace model.Trace, config analyzer.LinterPlugin) (analyzer.PluginResult, error) {
	res := make([]analyzer.RuleResult, len(config.Rules))

	for _, cfgRule := range config.Rules {
		rule, err := p.ruleRegistry.Get(cfgRule.Slug)
		if err != nil {
			return analyzer.PluginResult{}, err
		}

		result, err := rule.Evaluate(ctx, trace, cfgRule)
		if err != nil {
			return analyzer.PluginResult{}, err
		}

		res = append(res, result)
	}

	var allPassed bool = true
	for _, result := range res {
		if !result.Passed {
			allPassed = false
		}
	}

	return analyzer.PluginResult{
		//config
		Slug:        config.Slug,
		Name:        config.Name,
		Description: config.Description,

		// results
		Passed: allPassed,
		Rules:  res,
	}.CalculateResults(), nil
}
