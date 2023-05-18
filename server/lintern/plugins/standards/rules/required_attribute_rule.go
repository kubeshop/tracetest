package rules

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
)

type requiredAttributesRule struct {
	model.BaseRule
	attrMap RequiredAttributesMap
}

var DefaultAttrMap = NewRequiredAttributesMap(map[string][]string{
	"http": {"http.route"},
})

func NewRequiredAttributesRule() model.Rule {
	return requiredAttributesRule{
		BaseRule: model.BaseRule{
			Name:        "Required Attributes by Span Type",
			Description: "Ensure all spans have required attributes",
			Tips:        []string{"This rule checks if all required attributes are present in spans of given type"},
			Weight:      25,
		},
		attrMap: DefaultAttrMap,
	}
}

func (r requiredAttributesRule) Evaluate(ctx context.Context, trace model.Trace) (model.RuleResult, error) {
	results := make([]model.Result, 0)
	for _, span := range trace.Flat {
		spanType := span.Attributes["tracetest.span.type"]
		attributeList := r.attrMap.TypeAttributes(spanType)

		if len(attributeList) == 0 {
			results = append(results, model.Result{
				SpanID: span.ID.String(),
				Passed: true,
				Errors: []string{},
			})

			continue
		}

		missingAttributes := make([]string, 0)
		for _, requiredAttribute := range attributeList {
			if _, attributeExists := span.Attributes[requiredAttribute]; !attributeExists {
				missingAttributes = append(missingAttributes, requiredAttribute)
			}
		}

		numberOfMissingRequiredAttributes := len(missingAttributes)
		errors := make([]string, 0, int(numberOfMissingRequiredAttributes))
		for _, missingAttribute := range missingAttributes {
			errors = append(errors, fmt.Sprintf(`Attribute "%s" is missing from span of type "%s"`, missingAttribute, spanType))
		}

		results = append(results, model.Result{
			SpanID: span.ID.String(),
			Passed: numberOfMissingRequiredAttributes == 0,
			Errors: errors,
		})
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
