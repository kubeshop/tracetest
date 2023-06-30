package rules_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/metadata"
	"github.com/kubeshop/tracetest/server/linter/rules"
	"github.com/stretchr/testify/assert"
)

func TestNotEmptyAttributeRule(t *testing.T) {
	t.Run("no empty attributes", func(t *testing.T) {
		trace := traceWithSpans(
			spanWithAttributes("http", map[string]string{"http.method": "POST", "http.url": "http://localhost:11633"}),
			spanWithAttributes("http", map[string]string{"http.method": "GET", "http.url": "http://localhost:11633"}),
			spanWithAttributes("messaging", map[string]string{"messaging.target1": "user.sync", "messaging.system1": "kafka"}),
			spanWithAttributes("database", map[string]string{"db.statement": "INSERT INTO users (name, email) VALUES ($1, $2)"}),
		)

		rule := rules.NewNotEmptyAttributesRule(metadata.NotEmptyAttributesRule)
		result, _ := rule.Evaluate(context.Background(), trace, analyzer.LinterRule{})

		assert.True(t, result.Passed)
	})

	t.Run("empty attribute", func(t *testing.T) {
		trace := traceWithSpans(
			spanWithAttributes("http", map[string]string{"http.method": ""}),
		)

		rule := rules.NewNotEmptyAttributesRule(metadata.NotEmptyAttributesRule)
		result, _ := rule.Evaluate(context.Background(), trace, analyzer.LinterRule{})

		assert.False(t, result.Passed)
		assert.Len(t, result.Results, 1)
		assert.Equal(t, `Attribute "http.method" is empty`, result.Results[0].Errors[0].Description)
	})
}
