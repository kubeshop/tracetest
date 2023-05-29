package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/model"
)

type enforceHttpsProtocolRule struct {
	model.BaseRule
}

var (
	httpFields = []string{"http.scheme", "http.url"}
)

func NewEnforceHttpsProtocolRule() model.Rule {
	return &enforceHttpsProtocolRule{
		BaseRule: model.BaseRule{
			Name:        "Enforce HTTPS Protocol",
			Description: "Ensure all request use https",
			Tips:        []string{},
			Weight:      40,
		},
	}
}

func (r enforceHttpsProtocolRule) Evaluate(ctx context.Context, trace model.Trace) (model.RuleResult, error) {
	passed := true
	results := make([]model.Result, 0)
	for _, span := range trace.Flat {
		if span.Attributes.Get("tracetest.span.type") == "http" {
			result := r.validate(span)
			if !result.Passed {
				passed = false
			}
			results = append(results, result)
		}
	}

	return model.RuleResult{
		BaseRule: r.BaseRule,
		Passed:   passed,
		Results:  results,
	}, nil
}

func (r enforceHttpsProtocolRule) validate(span *model.Span) model.Result {
	insecureFields := make([]string, 0)
	for _, field := range httpFields {
		if span.Attributes.Get(field) != "" && !strings.Contains(span.Attributes.Get(field), "https") {
			insecureFields = append(insecureFields, fmt.Sprintf("Insecure http schema found for attribute: %s. Value: %s", field, span.Attributes.Get(field)))
		}
	}

	return model.Result{
		Passed: len(insecureFields) == 0,
		SpanID: span.ID.String(),
		Errors: insecureFields,
	}
}
