package rules

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/metadata"
	"github.com/kubeshop/tracetest/server/linter/results"
	"github.com/kubeshop/tracetest/server/model"
)

type ensureAttributeNamingRule struct {
	BaseRule
}

func NewEnsureAttributeNamingRule(metadata metadata.RuleMetadata) Rule {
	return &ensureAttributeNamingRule{
		BaseRule: NewRule(metadata),
	}
}

func (r ensureAttributeNamingRule) Evaluate(ctx context.Context, trace model.Trace, config analyzer.LinterRule) (results.RuleResult, error) {
	regex := regexp.MustCompile(`^([a-z0-9_]+\.)+[a-z0-9_]+$`)
	res := make([]results.Result, 0)
	passed := true

	if config.ErrorLevel != metadata.ErrorLevelDisabled {

		for _, span := range trace.Flat {
			errors := make([]results.Error, 0)
			namespaces := make([]string, 0)
			for name := range span.Attributes {
				if !regex.MatchString(name) {
					errors = append(errors, results.Error{
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
						errors = append(errors, results.Error{
							Value:       name,
							Description: fmt.Sprintf(`Attribute "%s" uses the same name as an existing namespace in the same span`, name),
						})
					}
				}
			}

			if len(errors) > 0 {
				passed = false
			}

			res = append(res, results.Result{
				SpanID: span.ID.String(),
				Passed: len(errors) == 0,
				Errors: errors,
			})
		}
	}

	return results.NewRuleResult(r.metadata, config, results.EvalRuleResult{Passed: passed, Results: res}), nil
}
