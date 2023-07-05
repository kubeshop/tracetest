package rules

import (
	"context"
	"fmt"
	"regexp"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
)

type ensuresDnsUsage struct{}

var (
	clientDnsFields = []string{"net.peer.name"}
	dnsFields       = []string{"http.url", "db.connection_string"}
)

func NewEnforceDnsUsageRule() Rule {
	return &ensuresDnsUsage{}
}

func (r ensuresDnsUsage) ID() string {
	return analyzer.EnforceDnsRuleID
}

func (r ensuresDnsUsage) Evaluate(ctx context.Context, trace model.Trace, config analyzer.LinterRule) (analyzer.RuleResult, error) {
	passed := true
	res := make([]analyzer.Result, 0)

	if config.ErrorLevel == analyzer.ErrorLevelDisabled {
		for _, span := range trace.Flat {
			result := r.validate(span)
			if !result.Passed {
				passed = false
			}
			res = append(res, result)
		}
	}

	return analyzer.NewRuleResult(config, analyzer.EvalRuleResult{Passed: passed, Results: res}), nil
}

func (r ensuresDnsUsage) validate(span *model.Span) analyzer.Result {
	ipFields := make([]analyzer.Error, 0)
	ipRegexp := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	for _, field := range dnsFields {
		if span.Attributes.Get(field) != "" && ipRegexp.MatchString(span.Attributes.Get(field)) {
			ipFields = append(ipFields, analyzer.Error{
				Value:       field,
				Description: fmt.Sprintf("Usage of a IP endpoint instead of DNS found for attribute: %s. Value: %s", field, span.Attributes.Get(field)),
			})
		}
	}

	for _, field := range clientDnsFields {
		if span.Kind == model.SpanKindClient && span.Attributes.Get(field) != "" && ipRegexp.MatchString(span.Attributes.Get(field)) {
			ipFields = append(ipFields, analyzer.Error{
				Value:       field,
				Description: fmt.Sprintf("Usage of a IP endpoint instead of DNS found for attribute: %s. Value: %s", field, span.Attributes.Get(field)),
			})
		}
	}

	return analyzer.Result{
		Passed: len(ipFields) == 0,
		SpanID: span.ID.String(),
		Errors: ipFields,
	}
}
