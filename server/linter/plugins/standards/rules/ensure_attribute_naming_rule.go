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
		errors := make([]model.Error, 0)
		namespaces := make([]string, 0)
		for name := range span.Attributes {
			if !regex.MatchString(name) {
				errors = append(errors, model.Error{
					Error:       "attribute_naming_error",
					Value:       name,
					Description: fmt.Sprintf(`Attribute "%s" does not follow the naming convention`, name),
				})
				continue
			}

			namespaces = append(namespaces, name[0:strings.LastIndex(name, ".")])
		}

		// ensure no attribute is named after a namespace
		for name := range span.Attributes {
			for _, namespace := range namespaces {
				if name == namespace {
					errors = append(errors, model.Error{
						Error:       "attribute_naming_error",
						Value:       name,
						Description: fmt.Sprintf(`Attribute "%s" uses the same name as an existing namespace in the same span`, name),
					})
				}
			}
		}

		if len(errors) > 0 {
			passed = false
		}

		results = append(results, model.Result{
			SpanID: span.ID.String(),
			Passed: len(errors) == 0,
			Errors: errors,
		})
	}

	return model.RuleResult{
		BaseRule: r.BaseRule,
		Passed:   passed,
		Results:  results,
	}, nil
}
