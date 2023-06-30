package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
)

type enforceHttpsProtocolRule struct {
	BaseRule
}

var (
	httpFields = []string{"http.scheme", "http.url"}
)

func NewEnforceHttpsProtocolRule() Rule {
	return &enforceHttpsProtocolRule{
		BaseRule: NewRule(analyzer.EnforceHttpsProtocolRuleSlug),
	}
}

func (r enforceHttpsProtocolRule) Evaluate(ctx context.Context, trace model.Trace, config analyzer.LinterRule) (analyzer.RuleResult, error) {
	passed := true
	res := make([]analyzer.Result, 0)

	if config.ErrorLevel != analyzer.ErrorLevelDisabled {
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

	return analyzer.NewRuleResult(config, analyzer.EvalRuleResult{Passed: passed, Results: res}), nil
}

func (r enforceHttpsProtocolRule) validate(span *model.Span) analyzer.Result {
	insecureFields := make([]analyzer.Error, 0)
	for _, field := range httpFields {
		if span.Attributes.Get(field) != "" && !strings.Contains(span.Attributes.Get(field), "https") {
			insecureFields = append(insecureFields, analyzer.Error{
				Value:       field,
				Description: fmt.Sprintf("Insecure http schema found for attribute: %s. Value: %s", field, span.Attributes.Get(field)),
			})
		}
	}

	return analyzer.Result{
		Passed: len(insecureFields) == 0,
		SpanID: span.ID.String(),
		Errors: insecureFields,
	}
}
