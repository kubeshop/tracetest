package rules_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/lintern/plugins/standards/rules"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestRequiredAttributesRule(t *testing.T) {
	trace := traceWithSpans(
		spanWithAttributes("http", map[string]string{"http.method": "POST", "http.url": "http://localhost:11633"}),
		spanWithAttributes("http", map[string]string{"http.method": "GET", "http.url": "http://localhost:11633"}),
		spanWithAttributes("messaging", map[string]string{"messaging.target1": "user.sync", "messaging.system1": "kafka"}),
		spanWithAttributes("database", map[string]string{"db.statement": "INSERT INTO users (name, email) VALUES ($1, $2)"}),
	)

	t.Run("When all required attributes are found", func(t *testing.T) {
		attrMap := rules.NewRequiredAttributesMap(map[string][]string{
			"http": {"http.method", "http.url"},
		})

		rule := rules.NewRequiredAttributesRule(attrMap)
		result, _ := rule.Evaluate(context.Background(), trace)

		for _, result := range result.Results {
			assert.True(t, result.Passed)
		}

		assert.True(t, result.Passed)
	})

	t.Run("When some attribute is missing", func(t *testing.T) {
		attrMap := rules.NewRequiredAttributesMap(map[string][]string{
			"database": {"database.kind", "db.statement"},
		})

		rule := rules.NewRequiredAttributesRule(attrMap)
		result, _ := rule.Evaluate(context.Background(), trace)

		assert.False(t, result.Passed)
	})

	t.Run("When all attributes are missing", func(t *testing.T) {
		attrMap := rules.NewRequiredAttributesMap(map[string][]string{
			"messaging": {"messaging.system", "messaging.target"},
		})

		rule := rules.NewRequiredAttributesRule(attrMap)
		result, _ := rule.Evaluate(context.Background(), trace)

		assert.False(t, result.Passed)
	})
}

func traceWithSpans(spans ...model.Span) model.Trace {
	trace := model.Trace{
		Flat: make(map[trace.SpanID]*model.Span, 0),
	}

	for _, span := range spans {
		realSpan := span
		span.ID = id.NewRandGenerator().SpanID()
		trace.Flat[span.ID] = &realSpan
	}

	return trace
}

func spanWithAttributes(spanType string, attributes map[string]string) model.Span {
	span := model.Span{
		Attributes: make(model.Attributes, 0),
	}

	for name, value := range attributes {
		span.Attributes[name] = value
	}

	span.Attributes["tracetest.span.type"] = spanType

	return span
}
