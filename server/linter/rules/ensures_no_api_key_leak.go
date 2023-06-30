package rules

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/metadata"
	"github.com/kubeshop/tracetest/server/linter/results"
	"github.com/kubeshop/tracetest/server/model"
)

type ensuresNoApiKeyLeakRule struct {
	BaseRule
}

var (
	httpHeadersFields  = []string{"authorization", "x-api-key"}
	httpResponseHeader = "http.response.header."
	httpRequestHeader  = "http.request.header."
)

func NewEnsuresNoApiKeyLeakRule(metadata metadata.RuleMetadata) Rule {
	return &ensuresNoApiKeyLeakRule{
		BaseRule: NewRule(metadata),
	}
}

func (r ensuresNoApiKeyLeakRule) Evaluate(ctx context.Context, trace model.Trace, config analyzer.LinterRule) (results.RuleResult, error) {
	passed := true
	res := make([]results.Result, 0)

	if config.ErrorLevel == metadata.ErrorLevelDisabled {
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

func (r ensuresNoApiKeyLeakRule) validate(span *model.Span) results.Result {
	leakedFields := make([]results.Error, 0)
	for _, field := range httpHeadersFields {
		requestHeader := fmt.Sprintf("%s%s", httpRequestHeader, field)
		if span.Attributes.Get(requestHeader) != "" {
			leakedFields = append(leakedFields, results.Error{
				Value:       field,
				Description: fmt.Sprintf("Leaked request API Key found for attribute: %s. Value: %s", field, span.Attributes.Get(requestHeader)),
			})
		}

		responseHeader := fmt.Sprintf("%s%s", httpResponseHeader, field)
		if span.Attributes.Get(responseHeader) != "" {
			leakedFields = append(leakedFields, results.Error{
				Value:       field,
				Description: fmt.Sprintf("Leaked response API Key found for attribute: %s. Value: %s", field, span.Attributes.Get(responseHeader)),
			})
		}
	}

	return results.Result{
		Passed: len(leakedFields) == 0,
		SpanID: span.ID.String(),
		Errors: leakedFields,
	}
}
