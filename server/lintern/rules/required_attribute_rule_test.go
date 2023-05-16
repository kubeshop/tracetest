package rules_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/lintern/rules"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestRequiredAttributesRule(t *testing.T) {
	trace := traceWithSpans(
		spanWithAttributes("http", map[string]string{"http.method": "POST", "http.url": "http://localhost:11633"}),
		spanWithAttributes("http", map[string]string{"http.method": "GET", "http.url": "http://localhost:11633"}),
		spanWithAttributes("messaging", map[string]string{"messaging.target": "user.sync", "messaging.system": "kafka"}),
		spanWithAttributes("database", map[string]string{"db.statement": "INSERT INTO users (name, email) VALUES ($1, $2)"}),
	)

	t.Run("When all required attributes are found", func(t *testing.T) {
		rule := rules.NewRequiredAttributesRule([]string{"http"}, []string{"http.method", "http.url"})
		result := rule.Run(context.Background(), trace)

		assert.True(t, result.Passed)
		for _, spanResult := range result.SpansResults {
			// expect score 50 because 2/2 attributes were found
			assert.Equal(t, uint(100), spanResult.Score)
		}
	})

	t.Run("When some attribute is missing", func(t *testing.T) {
		rule := rules.NewRequiredAttributesRule([]string{"database"}, []string{"database.kind", "db.statement"})
		result := rule.Run(context.Background(), trace)

		assert.False(t, result.Passed)

		for _, spanReresult := range result.SpansResults {
			// expect score 50 because 1/2 attributes were found
			assert.Equal(t, uint(50), spanReresult.Score)
		}
	})

	t.Run("When all attributes are missing", func(t *testing.T) {
		rule := rules.NewRequiredAttributesRule([]string{"http"}, []string{"http.protocol", "http.cors"})
		result := rule.Run(context.Background(), trace)

		assert.False(t, result.Passed)

		for _, spanReresult := range result.SpansResults {
			// expect score 50 because 1/2 attributes were found
			assert.Equal(t, uint(0), spanReresult.Score)
		}
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
