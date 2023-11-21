package linter_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/linter"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/plugins"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestLinter(t *testing.T) {
	trace := traceWithSpans(
		spanWithAttributes("http", map[string]string{"http.method": "POST", "http.url": "http://localhost:11633"}),
		spanWithAttributes("http", map[string]string{"http.method": "GET", "http.url": "http://localhost:11633"}),
		spanWithAttributes("messaging", map[string]string{"messaging.target": "user.sync", "messaging.system1": "kafka"}),
		spanWithAttributes("database", map[string]string{"db.statement": "INSERT INTO users (name, email) VALUES ($1, $2)"}),
	)

	t.Run("RunDefaults", func(t *testing.T) {
		linter := linter.NewLinter(linter.DefaultPluginRegistry)

		result, err := linter.Run(context.TODO(), trace, analyzer.GetDefaultLinter())
		assert.Nil(t, err)
		assert.Equal(t, 73, result.Score)
		assert.Equal(t, false, result.Passed)
		assert.Equal(t, 3, len(result.Plugins))
	})

	t.Run("RunEmpty", func(t *testing.T) {
		registry := plugins.NewRegistry()
		linter := linter.NewLinter(registry)

		result, err := linter.Run(context.TODO(), trace, analyzer.Linter{})
		assert.Nil(t, err)
		assert.Equal(t, 100, result.Score)
		assert.Equal(t, true, result.Passed)
		assert.Equal(t, 0, len(result.Plugins))
	})

	t.Run("RunMissingPluginRegistration", func(t *testing.T) {
		registry := plugins.NewRegistry()
		linter := linter.NewLinter(registry)

		_, err := linter.Run(context.TODO(), trace, analyzer.GetDefaultLinter())
		assert.NotNil(t, err)
		assert.Contains(t, "plugin standards not found", err.Error())
	})

	t.Run("RunEmptyEnabledPlugins", func(t *testing.T) {
		config := analyzer.GetDefaultLinter()
		plugins := make([]analyzer.LinterPlugin, 0)
		for _, plugin := range config.Plugins {
			plugin.Enabled = false
		}
		config.Plugins = plugins
		linter := linter.NewLinter(linter.DefaultPluginRegistry)

		result, err := linter.Run(context.TODO(), trace, config)
		assert.Nil(t, err)

		// if no plugins are enabled, then the score should be 100
		assert.Equal(t, 100, result.Score)
		assert.Equal(t, true, result.Passed)
		assert.Equal(t, 0, len(result.Plugins))
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
