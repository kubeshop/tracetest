package rules

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/metadata"
	"github.com/kubeshop/tracetest/server/linter/results"
	"github.com/kubeshop/tracetest/server/model"
)

type notEmptyRuleAttributesRule struct {
	BaseRule
}

func NewNotEmptyAttributesRule(metadata metadata.RuleMetadata) Rule {
	return &notEmptyRuleAttributesRule{
		BaseRule: NewRule(metadata),
	}
}

func (r notEmptyRuleAttributesRule) Evaluate(ctx context.Context, trace model.Trace, config analyzer.LinterRule) (results.RuleResult, error) {
	res := make([]results.Result, 0, len(trace.Flat))
	passed := true

	if config.ErrorLevel != metadata.ErrorLevelDisabled {
		for _, span := range trace.Flat {
			emptyAttributes := make([]string, 0)
			for name, value := range span.Attributes {
				if value == "" {
					emptyAttributes = append(emptyAttributes, name)
				}
			}

			errors := make([]results.Error, 0, len(emptyAttributes))
			for _, emptyAttribute := range emptyAttributes {
				errors = append(errors, results.Error{
					Value:       emptyAttribute,
					Description: fmt.Sprintf(`Attribute "%s" is empty`, emptyAttribute),
				},
				)
			}

			if len(errors) > 0 {
				passed = false
			}

			res = append(res, results.Result{
				SpanID: span.ID.String(),
				Passed: len(emptyAttributes) == 0,
				Errors: errors,
			})
		}
	}

	return results.NewRuleResult(r.metadata, config, results.EvalRuleResult{Passed: passed, Results: res}), nil
}
