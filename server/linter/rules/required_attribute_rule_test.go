package rules_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/rules"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/traces"
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
		trace := traceWithSpans(
			spanWithAttributes("http", map[string]string{"http.method": "POST", "http.url": "http://localhost:11633"}),
			spanWithAttributes("http", map[string]string{"http.method": "GET", "http.url": "http://localhost:11633"}),
			spanWithAttributes("messaging", map[string]string{"messaging.operation": "user.sync", "messaging.system": "kafka"}),
			spanWithAttributes("database", map[string]string{"db.system": "postgres"}),
		)

		rule := rules.NewRequiredAttributesRule()
		result, _ := rule.Evaluate(context.Background(), trace, analyzer.LinterRule{})

		for _, result := range result.Results {
			assert.True(t, result.Passed)
		}

		assert.True(t, result.Passed)
	})

	t.Run("When some attribute is missing", func(t *testing.T) {
		rule := rules.NewRequiredAttributesRule()
		result, _ := rule.Evaluate(context.Background(), trace, analyzer.LinterRule{})

		assert.False(t, result.Passed)
	})

	t.Run("When all attributes are missing", func(t *testing.T) {
		rule := rules.NewRequiredAttributesRule()
		result, _ := rule.Evaluate(context.Background(), trace, analyzer.LinterRule{})

		assert.False(t, result.Passed)
	})
}

func traceWithSpans(spans ...traces.Span) traces.Trace {
	trace := traces.Trace{
		Flat: make(map[trace.SpanID]*traces.Span, 0),
	}

	for _, span := range spans {
		realSpan := span
		span.ID = id.NewRandGenerator().SpanID()
		trace.Flat[span.ID] = &realSpan
	}

	return trace
}

func spanWithAttributes(spanType string, attributes map[string]string) traces.Span {
	span := traces.Span{
		Attributes: make(traces.Attributes, 0),
	}

	for name, value := range attributes {
		span.Attributes[name] = value
	}

	span.Attributes["tracetest.span.type"] = spanType

	return span
}
