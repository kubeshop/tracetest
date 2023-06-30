package plugins

import (
	"context"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/metadata"
	"github.com/kubeshop/tracetest/server/linter/results"
	"github.com/kubeshop/tracetest/server/linter/rules"
	"github.com/kubeshop/tracetest/server/model"
)

type BasePlugin struct {
	metadata     metadata.PluginMetadata
	ruleRegistry rules.RuleRegistry
}

func NewPlugin(metadata metadata.PluginMetadata, ruleRegistry rules.RuleRegistry) Plugin {
	return BasePlugin{
		metadata:     metadata,
		ruleRegistry: ruleRegistry,
	}
}

func (p BasePlugin) Slug() string {
	return p.metadata.Slug
}

func (p BasePlugin) RuleRegistry() rules.RuleRegistry {
	return p.ruleRegistry
}

func (p BasePlugin) Execute(ctx context.Context, trace model.Trace, config analyzer.LinterPlugin) (results.PluginResult, error) {
	res := make([]results.RuleResult, len(config.Rules))

	for _, cfgRule := range config.Rules {
		rule, err := p.ruleRegistry.Get(cfgRule.Slug)
		if err != nil {
			return results.PluginResult{}, err
		}

		result, err := rule.Evaluate(ctx, trace, cfgRule)
		if err != nil {
			return results.PluginResult{}, err
		}

		res = append(res, result)
	}

	var allPassed bool = true
	for _, result := range res {
		if !result.Passed {
			allPassed = false
		}
	}

	return results.PluginResult{
		// metadata
		Slug:        p.metadata.Slug,
		Name:        p.metadata.Name,
		Description: p.metadata.Description,

		// results
		Passed: allPassed,
		Rules:  res,
	}.CalculateResults(), nil
}
