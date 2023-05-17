package rules_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/lintern/plugins/standards/rules"
	"github.com/stretchr/testify/assert"
)

func TestEnsureNamingRule(t *testing.T) {
	t.Run("compliant trace", func(t *testing.T) {
		trace := traceWithSpans(
			spanWithAttributes("http", map[string]string{"http.method": "POST", "http.url": "http://localhost:11633"}),
			spanWithAttributes("http", map[string]string{"http.method": "GET", "http.url": "http://localhost:11633"}),
			spanWithAttributes("messaging", map[string]string{"messaging.target": "user.sync", "messaging.system1": "kafka"}),
			spanWithAttributes("database", map[string]string{"db.statement": "INSERT INTO users (name, email) VALUES ($1, $2)"}),
		)

		rule := rules.NewEnsureNamingRule()
		result, _ := rule.Evaluate(context.Background(), trace)

		assert.True(t, result.Passed)
	})

	t.Run("no namespace", func(t *testing.T) {
		trace := traceWithSpans(
			spanWithAttributes("http", map[string]string{"method": "POST"}),
		)

		rule := rules.NewEnsureNamingRule()
		result, _ := rule.Evaluate(context.Background(), trace)

		assert.False(t, result.Passed)
		assert.Len(t, result.Results, 1)
		assert.Equal(t, `Attribute "method" doesnt follow naming convention`, result.Results[0].Errors[0])
	})

	t.Run("span name with camel case", func(t *testing.T) {
		trace := traceWithSpans(
			spanWithAttributes("http", map[string]string{"http.statusCode": "200"}),
		)

		rule := rules.NewEnsureNamingRule()
		result, _ := rule.Evaluate(context.Background(), trace)

		assert.False(t, result.Passed)
		assert.Len(t, result.Results, 1)
		assert.Equal(t, `Attribute "http.statusCode" doesnt follow naming convention`, result.Results[0].Errors[0])
	})

	t.Run("attribute named after namespace", func(t *testing.T) {
		trace := traceWithSpans(
			spanWithAttributes("http", map[string]string{"tracetest.span": "POST"}),
		)

		rule := rules.NewEnsureNamingRule()
		result, _ := rule.Evaluate(context.Background(), trace)

		assert.False(t, result.Passed)
		assert.Len(t, result.Results, 1)
		assert.Equal(t, `Attribute "tracetest.span" uses the same name as an existing namespace in the same span`, result.Results[0].Errors[0])
	})
}
