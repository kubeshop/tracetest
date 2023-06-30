package rules

import (
	"context"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/metadata"
	"github.com/kubeshop/tracetest/server/linter/results"
	"github.com/kubeshop/tracetest/server/model"
)

type requiredAttributesRule struct {
	BaseRule
}

func NewRequiredAttributesRule(metadata metadata.RuleMetadata) Rule {
	return requiredAttributesRule{
		BaseRule: NewRule(metadata),
	}
}

func (r requiredAttributesRule) Evaluate(ctx context.Context, trace model.Trace, config analyzer.LinterRule) (results.RuleResult, error) {
	res := make([]results.Result, 0)
	for _, span := range trace.Flat {
		res = append(res, r.validateSpan(span))
	}

	var allPassed bool = true
	for _, result := range res {
		if !result.Passed {
			allPassed = false
		}
	}

	return results.NewRuleResult(r.metadata, config, results.EvalRuleResult{Passed: allPassed, Results: res}), nil
}
