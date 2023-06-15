package rules

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/kubeshop/tracetest/server/model"
)

type ensureAttributeNamingRule struct {
	model.BaseRule
}

func NewEnsureAttributeNamingRule() model.Rule {
	return &ensureAttributeNamingRule{
		BaseRule: model.BaseRule{
			Name:        "Attribute Naming",
			Description: "Ensure all attributes follow the naming convention",
			Tips: []string{
				"You should always add namespaces to your span names to ensure they will not be overwritten",
				"Use snake_case to separate multi-words. Ex: http.status_code instead of http.statusCode"},
			Weight: 25,
		},
	}
}

func (r ensureAttributeNamingRule) Evaluate(ctx context.Context, trace model.Trace) (model.RuleResult, error) {
	regex := regexp.MustCompile(`^([a-z0-9_]+\.)+[a-z0-9_]+$`)
	results := make([]model.Result, 0)
	passed := true

	for _, span := range trace.Flat {
		errors := make([]string, 0)
		namespaces := make([]string, 0)
		groupedErrors := make([]model.GroupedError, 0)
		values := make([]string, 0)
		for name := range span.Attributes {
			if !regex.MatchString(name) {
				errors = append(errors, fmt.Sprintf(`Attribute "%s" doesnt follow naming convention`, name))
				values = append(values, name)
				continue
			}

			namespaces = append(namespaces, name[0:strings.LastIndex(name, ".")])
		}

		if len(values) > 1 {
			groupedErrors = append(groupedErrors, model.GroupedError{
				Error:  "The following attributes do not adhere to the naming convention:",
				Values: values,
			})
		}

		values = make([]string, 0)
		// ensure no attribute is named after a namespace
		for name := range span.Attributes {
			for _, namespace := range namespaces {
				if name == namespace {
					errors = append(errors, fmt.Sprintf(`Attribute "%s" uses the same name as an existing namespace in the same span`, name))
					values = append(values, name)
				}
			}
		}

		if len(values) > 1 {
			groupedErrors = append(groupedErrors, model.GroupedError{
				Error:  "The following attributes use the same name as an existing namespace in the same span:",
				Values: values,
			})
		}

		if len(errors) > 0 {
			passed = false
		}

		results = append(results, model.Result{
			SpanID:        span.ID.String(),
			Passed:        len(errors) == 0,
			Errors:        errors,
			GroupedErrors: groupedErrors,
		})
	}

	return model.RuleResult{
		BaseRule: r.BaseRule,
		Passed:   passed,
		Results:  results,
	}, nil
}
