package plugins_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/plugins"
	"github.com/kubeshop/tracetest/server/linter/rules"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestAnalyzerEntities(t *testing.T) {
	plugin := plugins.NewPlugin(
		analyzer.StandardsID,
		rules.NewRegistry().
			Register(rules.NewEnsureSpanNamingRule()).
			Register(rules.NewRequiredAttributesRule()).
			Register(rules.NewEnsureAttributeNamingRule()).
			Register(rules.NewNotEmptyAttributesRule()),
	)

	trace := traceWithSpans(
		spanWithAttributes("http", map[string]string{"http.method": "POST", "http.url": "http://localhost:11633"}),
		spanWithAttributes("http", map[string]string{"http.method": "GET", "http.url": "http://localhost:11633"}),
		spanWithAttributes("messaging", map[string]string{"messaging.target": "user.sync", "messaging.system1": "kafka"}),
		spanWithAttributes("database", map[string]string{"db.statement": "INSERT INTO users (name, email) VALUES ($1, $2)"}),
	)

	t.Run("PluginExecute", func(t *testing.T) {
		result, err := plugin.Execute(context.TODO(), trace, analyzer.StandardsPlugin)

		assert.Nil(t, err)
		assert.Equal(t, analyzer.StandardsID, result.ID)
		assert.Equal(t, 50, result.Score)
		assert.Equal(t, false, result.Passed)
	})

	t.Run("EmptyPluginExecute", func(t *testing.T) {
		emptyPlugin := plugins.NewPlugin(
			analyzer.StandardsID,
			rules.NewRegistry(),
		)
		result, err := emptyPlugin.Execute(context.TODO(), trace, analyzer.LinterPlugin{})

		assert.Nil(t, err)
		assert.Equal(t, analyzer.StandardsID, result.ID)
		assert.Equal(t, 100, result.Score)
		assert.Equal(t, true, result.Passed)
	})

	t.Run("PluginConfigMismatchExecute", func(t *testing.T) {
		emptyPlugin := plugins.NewPlugin(
			analyzer.StandardsID,
			rules.NewRegistry(),
		)
		_, err := emptyPlugin.Execute(context.TODO(), trace, analyzer.StandardsPlugin)

		assert.NotNil(t, err)
		assert.Contains(t, "rule span-naming not found", err.Error())
	})
}

func spanWithAttributes(spanType string, attributes map[string]string) traces.Span {
	span := traces.Span{
		Attributes: traces.NewAttributes(),
	}

	for name, value := range attributes {
		span.Attributes.Set(name, value)
	}

	span.Attributes.Set("tracetest.span.type", spanType)

	return span
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
