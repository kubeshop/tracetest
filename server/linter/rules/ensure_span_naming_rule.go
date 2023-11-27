package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/traces"
)

type ensureSpanNamingRule struct{}

func NewEnsureSpanNamingRule() Rule {
	return &ensureSpanNamingRule{}
}

func (r ensureSpanNamingRule) ID() string {
	return analyzer.EnsureSpanNamingRuleID
}

func (r ensureSpanNamingRule) Evaluate(ctx context.Context, trace traces.Trace, config analyzer.LinterRule) (analyzer.RuleResult, error) {
	res := make([]analyzer.Result, 0)
	hasErrors := false

	if config.ErrorLevel != analyzer.ErrorLevelDisabled {
		for _, span := range trace.Flat {
			result := r.validateSpanName(ctx, span)
			res = append(res, result)

			hasErrors = hasErrors || !result.Passed
		}
	}

	return analyzer.NewRuleResult(config, analyzer.EvalRuleResult{Passed: !hasErrors, Results: res}), nil
}

func (r ensureSpanNamingRule) validateSpanName(ctx context.Context, span *traces.Span) analyzer.Result {
	switch span.Attributes.Get("tracetest.span.type") {
	case "http":
		return r.validateHTTPSpanName(ctx, span)
	case "database":
		return r.validateDatabaseSpanName(ctx, span)
	case "rpc":
		return r.validateRPCSpanName(ctx, span)
	case "messaging":
		return r.validateMessagingSpanName(ctx, span)
	}

	return analyzer.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}

func (r ensureSpanNamingRule) validateHTTPSpanName(ctx context.Context, span *traces.Span) analyzer.Result {
	expectedName := ""
	if span.Kind == traces.SpanKindClient {
		expectedName = span.Attributes.Get("http.method")
	} else {
		expectedName = fmt.Sprintf("%s %s", span.Attributes.Get("http.method"), span.Attributes.Get("http.route"))
	}

	if span.Name != expectedName {
		return analyzer.Result{
			SpanID: span.ID.String(),
			Passed: false,
			Errors: []analyzer.Error{
				{
					Value:       "tracetest.span.name",
					Expected:    expectedName,
					Description: fmt.Sprintf(`The Span name %s is not matching the naming convention`, span.Name),
				},
			},
		}
	}

	return analyzer.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}

func (r ensureSpanNamingRule) validateDatabaseSpanName(ctx context.Context, span *traces.Span) analyzer.Result {
	dbOperation := span.Attributes.Get("db.operation")
	dbName := span.Attributes.Get("db.name")
	tableName := span.Attributes.Get("db.sql.table")
	expectedName := fmt.Sprintf("%s %s.%s", dbOperation, dbName, tableName)

	// TODO: fix this by adding proper validation for all other db systems
	dbSystem := span.Attributes.Get("db.system")
	if dbSystem == "redis" {
		dbStatement := strings.ToLower(span.Attributes.Get("db.statement"))
		return analyzer.Result{
			Passed: dbStatement != "" && strings.HasPrefix(dbStatement, strings.ToLower(span.Name)),
			SpanID: span.ID.String(),
		}
	}

	if strings.Trim(span.Name, "") != strings.Trim(expectedName, "") {
		return analyzer.Result{
			SpanID: span.ID.String(),
			Passed: false,
			Errors: []analyzer.Error{
				{
					Value:       "tracetest.span.name",
					Expected:    expectedName,
					Description: fmt.Sprintf(`The Span name %s is not matching the naming convention`, span.Name),
				},
			},
		}
	}

	return analyzer.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}

func (r ensureSpanNamingRule) validateRPCSpanName(ctx context.Context, span *traces.Span) analyzer.Result {
	rpcService := span.Attributes.Get("rpc.service")
	rpcMethod := span.Attributes.Get("rpc.method")

	expectedName := fmt.Sprintf("%s/%s", rpcService, rpcMethod)

	if span.Name != expectedName {
		return analyzer.Result{
			SpanID: span.ID.String(),
			Passed: false,
			Errors: []analyzer.Error{
				{
					Value:       "tracetest.span.name",
					Expected:    expectedName,
					Description: fmt.Sprintf(`The Span name %s is not matching the naming convention`, span.Name),
				},
			},
		}
	}

	return analyzer.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}

func (r ensureSpanNamingRule) validateMessagingSpanName(ctx context.Context, span *traces.Span) analyzer.Result {
	destination := span.Attributes.Get("messaging.destination")
	operation := span.Attributes.Get("messaging.operation")

	expectedName := fmt.Sprintf("%s %s", destination, operation)

	if span.Name != expectedName {
		return analyzer.Result{
			SpanID: span.ID.String(),
			Passed: false,
			Errors: []analyzer.Error{
				{
					Value:       "tracetest.span.name",
					Expected:    expectedName,
					Description: fmt.Sprintf(`The Span name %s is not matching the naming convention`, span.Name),
				},
			},
		}
	}

	return analyzer.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}
