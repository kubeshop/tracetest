package rules

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/traces"
)

type ensureAttributeNamingRule struct{}

func NewEnsureAttributeNamingRule() Rule {
	return &ensureAttributeNamingRule{}
}

func (r ensureAttributeNamingRule) ID() string {
	return analyzer.EnsureAttributeNamingRuleID
}

func (r ensureAttributeNamingRule) Evaluate(ctx context.Context, trace traces.Trace, config analyzer.LinterRule) (analyzer.RuleResult, error) {
	regex := regexp.MustCompile(`^([a-z0-9_]+\.)+[a-z0-9_]+$`)
	res := make([]analyzer.Result, 0)
	passed := true

	if config.ErrorLevel != analyzer.ErrorLevelDisabled {

		for _, span := range trace.Flat {
			errors := make([]analyzer.Error, 0)
			namespaces := make([]string, 0)
			for name := range span.Attributes {
				if !regex.MatchString(name) {
					errors = append(errors, analyzer.Error{
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
						errors = append(errors, analyzer.Error{
							Value:       name,
							Description: fmt.Sprintf(`Attribute "%s" uses the same name as an existing namespace in the same span`, name),
						})
					}
				}
			}

			if len(errors) > 0 {
				passed = false
			}

			res = append(res, analyzer.Result{
				SpanID: span.ID.String(),
				Passed: len(errors) == 0,
				Errors: errors,
			})
		}
	}

	return analyzer.NewRuleResult(config, analyzer.EvalRuleResult{Passed: passed, Results: res}), nil
}
