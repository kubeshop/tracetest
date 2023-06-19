package rules

import (
	"context"
	"fmt"
	"strings"

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
	hasErrors := false
	for _, span := range trace.Flat {
		result := r.validateSpanName(ctx, span)
		results = append(results, result)

		hasErrors = hasErrors || !result.Passed
	}

	return model.RuleResult{
		BaseRule: r.BaseRule,
		Results:  results,
		Passed:   !hasErrors,
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
			Errors: []model.Error{
				{
					Error:       "span_naming_error",
					Value:       span.Name,
					Expected:    expectedName,
					Description: fmt.Sprintf(`The Span name %s is not matching the naming convention`, span.Name),
				},
			},
		}
	}

	return model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}

func (r ensureSpanNamingRule) validateDatabaseSpanName(ctx context.Context, span *model.Span) model.Result {
	dbOperation := span.Attributes.Get("db.operation")
	dbName := span.Attributes.Get("db.name")
	tableName := span.Attributes.Get("db.sql.table")
	expectedName := fmt.Sprintf("%s %s.%s", dbOperation, dbName, tableName)

	// TODO: fix this by adding proper validation for all other db systems
	dbSystem := span.Attributes.Get("db.system")
	if dbSystem == "redis" {
		dbStatement := strings.ToLower(span.Attributes.Get("db.statement"))
		return model.Result{
			Passed: dbStatement != "" && strings.HasPrefix(dbStatement, strings.ToLower(span.Name)),
			SpanID: span.ID.String(),
		}
	}

	if strings.Trim(span.Name, "") != strings.Trim(expectedName, "") {
		return model.Result{
			SpanID: span.ID.String(),
			Passed: false,
			Errors: []model.Error{
				{
					Error:       "span_naming_error",
					Value:       span.Name,
					Expected:    expectedName,
					Description: fmt.Sprintf(`The Span name %s is not matching the naming convention`, span.Name),
				},
			},
		}
	}

	return model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}

func (r ensureSpanNamingRule) validateRPCSpanName(ctx context.Context, span *model.Span) model.Result {
	rpcService := span.Attributes.Get("rpc.service")
	rpcMethod := span.Attributes.Get("rpc.method")

	expectedName := fmt.Sprintf("%s/%s", rpcService, rpcMethod)

	if span.Name != expectedName {
		return model.Result{
			SpanID: span.ID.String(),
			Passed: false,
			Errors: []model.Error{
				{
					Error:       "span_naming_error",
					Value:       span.Name,
					Expected:    expectedName,
					Description: fmt.Sprintf(`The Span name %s is not matching the naming convention`, span.Name),
				},
			},
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
			Errors: []model.Error{
				{
					Error:       "span_naming_error",
					Value:       span.Name,
					Expected:    expectedName,
					Description: fmt.Sprintf(`The Span name %s is not matching the naming convention`, span.Name),
				},
			},
		}
	}

	return model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}
