package rules

import (
	"context"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
)

type requiredAttributesRule struct {
	BaseRule
}

func NewRequiredAttributesRule() Rule {
	return requiredAttributesRule{
		BaseRule: NewRule(analyzer.RequiredAttributesRuleSlug),
	}
}

func (r requiredAttributesRule) Evaluate(ctx context.Context, trace model.Trace, config analyzer.LinterRule) (analyzer.RuleResult, error) {
	res := make([]analyzer.Result, 0)
	var allPassed bool = true

	if config.ErrorLevel != analyzer.ErrorLevelDisabled {
		for _, span := range trace.Flat {
			res = append(res, r.validateSpan(span))
		}

		for _, result := range res {
			if !result.Passed {
				allPassed = false
			}
		}
	}

	return analyzer.NewRuleResult(config, analyzer.EvalRuleResult{Passed: allPassed, Results: res}), nil
}
