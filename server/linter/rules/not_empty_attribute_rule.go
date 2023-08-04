package rules

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/traces"
)

type notEmptyRuleAttributesRule struct{}

func NewNotEmptyAttributesRule() Rule {
	return &notEmptyRuleAttributesRule{}
}

func (r notEmptyRuleAttributesRule) ID() string {
	return analyzer.NotEmptyAttributesRuleID
}

func (r notEmptyRuleAttributesRule) Evaluate(ctx context.Context, trace traces.Trace, config analyzer.LinterRule) (analyzer.RuleResult, error) {
	res := make([]analyzer.Result, 0, len(trace.Flat))
	passed := true

	if config.ErrorLevel != analyzer.ErrorLevelDisabled {
		for _, span := range trace.Flat {
			emptyAttributes := make([]string, 0)
			for name, value := range span.Attributes {
				if value == "" {
					emptyAttributes = append(emptyAttributes, name)
				}
			}

			errors := make([]analyzer.Error, 0, len(emptyAttributes))
			for _, emptyAttribute := range emptyAttributes {
				errors = append(errors, analyzer.Error{
					Value:       emptyAttribute,
					Description: fmt.Sprintf(`Attribute "%s" is empty`, emptyAttribute),
				},
				)
			}

			if len(errors) > 0 {
				passed = false
			}

			res = append(res, analyzer.Result{
				SpanID: span.ID.String(),
				Passed: len(emptyAttributes) == 0,
				Errors: errors,
			})
		}
	}

	return analyzer.NewRuleResult(config, analyzer.EvalRuleResult{Passed: passed, Results: res}), nil
}
