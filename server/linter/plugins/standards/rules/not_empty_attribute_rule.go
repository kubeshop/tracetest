package rules

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
)

type notEmptyRuleAttributesRule struct {
	model.BaseRule
}

func NewNotEmptyAttributesRule() model.Rule {
	return &notEmptyRuleAttributesRule{
		BaseRule: model.BaseRule{
			Name:        "Not Empty Attributes",
			Description: "Does not allow empty attribute values in any span",
			Tips:        []string{"Empty attributes don't provide any information about the operation and should be removed"},
			Weight:      25,
		},
	}
}

func (r notEmptyRuleAttributesRule) Evaluate(ctx context.Context, trace model.Trace) (model.RuleResult, error) {
	results := make([]model.Result, 0, len(trace.Flat))
	passed := true
	for _, span := range trace.Flat {
		emptyAttributes := make([]string, 0)
		for name, value := range span.Attributes {
			if value == "" {
				emptyAttributes = append(emptyAttributes, name)
			}
		}

		errors := make([]string, 0, len(emptyAttributes))
		for _, emptyAttribute := range emptyAttributes {
			errors = append(errors, fmt.Sprintf(`Attribute "%s" is empty`, emptyAttribute))
		}

		if len(errors) > 0 {
			passed = false
		}

		results = append(results, model.Result{
			SpanID: span.ID.String(),
			Passed: len(emptyAttributes) == 0,
			Errors: errors,
		})
	}

	return model.RuleResult{
		BaseRule: r.BaseRule,
		Passed:   passed,
		Results:  results,
	}, nil
}
