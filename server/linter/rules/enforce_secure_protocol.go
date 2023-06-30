package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/metadata"
	"github.com/kubeshop/tracetest/server/linter/results"
	"github.com/kubeshop/tracetest/server/model"
)

type enforceHttpsProtocolRule struct {
	BaseRule
}

var (
	httpFields = []string{"http.scheme", "http.url"}
)

func NewEnforceHttpsProtocolRule(metadata metadata.RuleMetadata) Rule {
	return &enforceHttpsProtocolRule{
		BaseRule: NewRule(metadata),
	}
}

func (r enforceHttpsProtocolRule) Evaluate(ctx context.Context, trace model.Trace, config analyzer.LinterRule) (results.RuleResult, error) {
	passed := true
	res := make([]results.Result, 0)

	if config.ErrorLevel != metadata.ErrorLevelDisabled {
		for _, span := range trace.Flat {
			if span.Attributes.Get("tracetest.span.type") == "http" {
				result := r.validate(span)
				if !result.Passed {
					passed = false
				}
				res = append(res, result)
			}
		}
	}

	return results.NewRuleResult(r.metadata, config, results.EvalRuleResult{Passed: passed, Results: res}), nil
}

func (r enforceHttpsProtocolRule) validate(span *model.Span) results.Result {
	insecureFields := make([]results.Error, 0)
	for _, field := range httpFields {
		if span.Attributes.Get(field) != "" && !strings.Contains(span.Attributes.Get(field), "https") {
			insecureFields = append(insecureFields, results.Error{
				Value:       field,
				Description: fmt.Sprintf("Insecure http schema found for attribute: %s. Value: %s", field, span.Attributes.Get(field)),
			})
		}
	}

	return results.Result{
		Passed: len(insecureFields) == 0,
		SpanID: span.ID.String(),
		Errors: insecureFields,
	}
}
