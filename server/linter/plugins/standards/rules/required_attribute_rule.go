package rules

import (
	"context"

	"github.com/kubeshop/tracetest/server/model"
)

type requiredAttributesRule struct {
	model.BaseRule
}

func NewRequiredAttributesRule() model.Rule {
	return requiredAttributesRule{
		BaseRule: model.BaseRule{
			Name:             "Required Attributes by Span Type",
			Description:      "Ensure all spans have required attributes",
			ErrorDescription: "This span is missing the following required attributes:",
			Tips:             []string{"This rule checks if all required attributes are present in spans of given type"},
			Weight:           25,
		},
	}
}

func (r requiredAttributesRule) Evaluate(ctx context.Context, trace model.Trace) (model.RuleResult, error) {
	results := make([]model.Result, 0)
	for _, span := range trace.Flat {
		results = append(results, r.validateSpan(span))
	}

	var allPassed bool = true
	for _, result := range results {
		if !result.Passed {
			allPassed = false
		}
	}

	return model.RuleResult{
		BaseRule: r.BaseRule,
		Passed:   allPassed,
		Results:  results,
	}, nil
}
