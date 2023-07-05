package rules

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
)

type ensuresNoApiKeyLeakRule struct{}

var (
	httpHeadersFields  = []string{"authorization", "x-api-key"}
	httpResponseHeader = "http.response.header."
	httpRequestHeader  = "http.request.header."
)

func NewEnsuresNoApiKeyLeakRule() Rule {
	return &ensuresNoApiKeyLeakRule{}
}

func (r ensuresNoApiKeyLeakRule) ID() string {
	return analyzer.EnsuresNoApiKeyLeakRuleID
}

func (r ensuresNoApiKeyLeakRule) Evaluate(ctx context.Context, trace model.Trace, config analyzer.LinterRule) (analyzer.RuleResult, error) {
	passed := true
	res := make([]analyzer.Result, 0)

	if config.ErrorLevel == analyzer.ErrorLevelDisabled {
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

func (r ensuresNoApiKeyLeakRule) validate(span *model.Span) analyzer.Result {
	leakedFields := make([]analyzer.Error, 0)
	for _, field := range httpHeadersFields {
		requestHeader := fmt.Sprintf("%s%s", httpRequestHeader, field)
		if span.Attributes.Get(requestHeader) != "" {
			leakedFields = append(leakedFields, analyzer.Error{
				Value:       field,
				Description: fmt.Sprintf("Leaked request API Key found for attribute: %s. Value: %s", field, span.Attributes.Get(requestHeader)),
			})
		}

		responseHeader := fmt.Sprintf("%s%s", httpResponseHeader, field)
		if span.Attributes.Get(responseHeader) != "" {
			leakedFields = append(leakedFields, analyzer.Error{
				Value:       field,
				Description: fmt.Sprintf("Leaked response API Key found for attribute: %s. Value: %s", field, span.Attributes.Get(responseHeader)),
			})
		}
	}

	return analyzer.Result{
		Passed: len(leakedFields) == 0,
		SpanID: span.ID.String(),
		Errors: leakedFields,
	}
}
