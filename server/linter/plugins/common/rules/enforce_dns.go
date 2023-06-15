package rules

import (
	"context"
	"fmt"
	"regexp"

	"github.com/kubeshop/tracetest/server/model"
)

type ensuresDnsUsage struct {
	model.BaseRule
}

var (
	clientDnsFields = []string{"net.peer.name"}
	dnsFields       = []string{"http.url", "db.connection_string"}
)

func NewEnforceDnsUsageRule() model.Rule {
	return &ensuresDnsUsage{
		BaseRule: model.BaseRule{
			Name:        "DNS Over IP",
			Description: "Enforce DNS usage over IP addresses",
			Tips:        []string{},
			Weight:      100,
		},
	}
}

func (r ensuresDnsUsage) Evaluate(ctx context.Context, trace model.Trace) (model.RuleResult, error) {
	passed := true
	results := make([]model.Result, 0)
	for _, span := range trace.Flat {
		result := r.validate(span)
		if !result.Passed {
			passed = false
		}
		results = append(results, result)
	}

	return model.RuleResult{
		BaseRule: r.BaseRule,
		Passed:   passed,
		Results:  results,
	}, nil
}

func (r ensuresDnsUsage) validate(span *model.Span) model.Result {
	ipFields := make([]string, 0)
	ipRegexp := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	for _, field := range dnsFields {
		if span.Attributes.Get(field) != "" && ipRegexp.MatchString(span.Attributes.Get(field)) {
			ipFields = append(ipFields, fmt.Sprintf("Usage of a IP endpoint instead of DNS found for attribute: %s. Value: %s", field, span.Attributes.Get(field)))
		}
	}

	for _, field := range clientDnsFields {
		if span.Kind == model.SpanKindClient && span.Attributes.Get(field) != "" && ipRegexp.MatchString(span.Attributes.Get(field)) {
			ipFields = append(ipFields, fmt.Sprintf("Usage of a IP endpoint instead of DNS found for attribute: %s. Value: %s", field, span.Attributes.Get(field)))
		}
	}

	return model.Result{
		Passed: len(ipFields) == 0,
		SpanID: span.ID.String(),
		Errors: ipFields,
	}
}
