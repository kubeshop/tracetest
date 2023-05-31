package rules

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
)

type ensureSpanNamingRule struct {
	model.BaseRule
}

func NewEnsureSpanNamingRule() model.Rule {
	return &ensureSpanNamingRule{
		BaseRule: model.BaseRule{
			Name:        "Span Name Convention",
			Description: "Ensure all span follow the name convention",
			Tips:        []string{},
			Weight:      25,
		},
	}
}

func (r ensureSpanNamingRule) Evaluate(ctx context.Context, trace model.Trace) (model.RuleResult, error) {
	results := make([]model.Result, 0)
	for _, span := range trace.Flat {
		results = append(results, r.validateSpanName(ctx, span))
	}

	return model.RuleResult{
		BaseRule: r.BaseRule,
		Results:  results,
	}, nil
}

func (r ensureSpanNamingRule) validateSpanName(ctx context.Context, span *model.Span) model.Result {
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

	return model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}

func (r ensureSpanNamingRule) validateHTTPSpanName(ctx context.Context, span *model.Span) model.Result {
	expectedName := ""
	if span.Kind == model.SpanKindServer {
		expectedName = fmt.Sprintf("%s %s", span.Attributes.Get("http.method"), span.Attributes.Get("http.route"))
	}

	if span.Kind == model.SpanKindClient {
		expectedName = span.Attributes.Get("http.method")
	}

	if span.Name != expectedName {
		return model.Result{
			SpanID: span.ID.String(),
			Passed: false,
			Errors: []string{fmt.Sprintf(`Span name is not matching the naming convention. Expected: %s`, expectedName)},
		}
	}

	return model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}

func (r ensureSpanNamingRule) validateDatabaseSpanName(ctx context.Context, span *model.Span) model.Result {
	var expectedName string
	dbOperation := span.Attributes.Get("db.operation")
	dbName := span.Attributes.Get("db.operation")
	tableName := span.Attributes.Get("db.sql.table")

	if tableName != "" {
		expectedName = fmt.Sprintf("%s %s %s", dbOperation, dbName, tableName)
	} else {
		expectedName = fmt.Sprintf("%s %s", dbOperation, dbName)
	}

	if span.Name != expectedName {
		return model.Result{
			SpanID: span.ID.String(),
			Passed: false,
			Errors: []string{fmt.Sprintf(`Span name is not matching the naming convention. Expected: %s`, expectedName)},
		}
	}

	return model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}

func (r ensureSpanNamingRule) validateRPCSpanName(ctx context.Context, span *model.Span) model.Result {
	rpcPackage := span.Attributes.Get("rpc.package")
	rpcService := span.Attributes.Get("rpc.service")
	rpcMethod := span.Attributes.Get("rpc.method")

	expectedName := fmt.Sprintf("%s.%s/%s", rpcPackage, rpcService, rpcMethod)

	if span.Name != expectedName {
		return model.Result{
			SpanID: span.ID.String(),
			Passed: false,
			Errors: []string{fmt.Sprintf(`Span name is not matching the naming convention. Expected: %s`, expectedName)},
		}
	}

	return model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}

func (r ensureSpanNamingRule) validateMessagingSpanName(ctx context.Context, span *model.Span) model.Result {
	destination := span.Attributes.Get("messaging.destination")
	operation := span.Attributes.Get("messaging.operation")

	expectedName := fmt.Sprintf("%s %s", destination, operation)

	if span.Name != expectedName {
		return model.Result{
			SpanID: span.ID.String(),
			Passed: false,
			Errors: []string{fmt.Sprintf(`Span name is not matching the naming convention. Expected: %s`, expectedName)},
		}
	}

	return model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}
