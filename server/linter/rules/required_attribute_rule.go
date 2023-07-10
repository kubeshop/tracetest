package rules

import (
	"context"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
)

type requiredAttributesRule struct{}

func NewRequiredAttributesRule() Rule {
	return requiredAttributesRule{}
}

func (r requiredAttributesRule) ID() string {
	return analyzer.RequiredAttributesRuleID
}

func (r requiredAttributesRule) Evaluate(ctx context.Context, trace model.Trace, config analyzer.LinterRule) (analyzer.RuleResult, error) {
	res := make([]analyzer.Result, 0)
	var allPassed bool = true

	if config.ErrorLevel != analyzer.ErrorLevelDisabled {
		for _, span := range trace.Flat {
			analyzerResult := r.validateSpan(span)
			res = append(res, analyzerResult)

			if !analyzerResult.Passed {
				allPassed = false
			}
		}
	}

	return analyzer.NewRuleResult(config, analyzer.EvalRuleResult{Passed: allPassed, Results: res}), nil
}
