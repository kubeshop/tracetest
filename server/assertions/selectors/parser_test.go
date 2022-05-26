package selectors_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseExpressionOrFail(t *testing.T, expression string) *selectors.ParserSelector {
	parser, err := selectors.CreateParser()
	require.NoError(t, err, "parser should be created successfully")

	selector := &selectors.ParserSelector{}
	err = parser.ParseString("", expression, selector)
	require.NoError(t, err, "parser should be able to parse expression")

	return selector
}

func TestParserSingleSpanStringProperty(t *testing.T) {
	expression := "span[service.name=\"Pokeshop\"]"

	selector := parseExpressionOrFail(t, expression)

	assert.Equal(t, "service.name", selector.SpanSelectors[0].Filters[0].Property)
	assert.Equal(t, "=", selector.SpanSelectors[0].Filters[0].Operator)
	assert.Equal(t, "\"Pokeshop\"", *selector.SpanSelectors[0].Filters[0].Value.String)
}

func TestParserSingleSpanIntProperty(t *testing.T) {
	expression := "span[http.status_code=200]"

	selector := parseExpressionOrFail(t, expression)

	assert.Equal(t, "http.status_code", selector.SpanSelectors[0].Filters[0].Property)
	assert.Equal(t, "=", selector.SpanSelectors[0].Filters[0].Operator)
	assert.Equal(t, int64(200), *selector.SpanSelectors[0].Filters[0].Value.Int)
}

func TestParserSingleSpanMultipleAttributes(t *testing.T) {
	expression := "span[service.name=\"Pokeshop\" tracetest.span.type=\"http\"]"

	selector := parseExpressionOrFail(t, expression)

	assert.Equal(t, "service.name", selector.SpanSelectors[0].Filters[0].Property)
	assert.Equal(t, "\"Pokeshop\"", *selector.SpanSelectors[0].Filters[0].Value.String)
	assert.Equal(t, "tracetest.span.type", selector.SpanSelectors[0].Filters[1].Property)
	assert.Equal(t, "\"http\"", *selector.SpanSelectors[0].Filters[1].Value.String)
}

func TestParserSingleSpanUsingContainsComparator(t *testing.T) {
	expression := "span[service.name contains \"api-\"]"

	selector := parseExpressionOrFail(t, expression)

	assert.Equal(t, "service.name", selector.SpanSelectors[0].Filters[0].Property)
	assert.Equal(t, "contains", selector.SpanSelectors[0].Filters[0].Operator)
	assert.Equal(t, "\"api-\"", *selector.SpanSelectors[0].Filters[0].Value.String)
}

func TestParserSingleSpanWithNthChild(t *testing.T) {
	expression := "span[tracetest.span.type=\"http\"]:nth_child(2)"

	selector := parseExpressionOrFail(t, expression)

	assert.Equal(t, "tracetest.span.type", selector.SpanSelectors[0].Filters[0].Property)
	assert.Equal(t, "=", selector.SpanSelectors[0].Filters[0].Operator)
	assert.Equal(t, "\"http\"", *selector.SpanSelectors[0].Filters[0].Value.String)
	assert.Equal(t, "nth_child", selector.SpanSelectors[0].PseudoClass.Type)
	assert.Equal(t, int64(2), *selector.SpanSelectors[0].PseudoClass.Value.Int)
}

func TestParserWithSpanHierarchy(t *testing.T) {
	expression := "span[service.name=\"Pokeshop\"] span[tracetest.span.type=\"http\"]"

	selector := parseExpressionOrFail(t, expression)

	assert.Equal(t, "service.name", selector.SpanSelectors[0].Filters[0].Property)
	assert.Equal(t, "tracetest.span.type", selector.SpanSelectors[0].ChildSelector.Filters[0].Property)
	assert.Equal(t, "\"http\"", *selector.SpanSelectors[0].ChildSelector.Filters[0].Value.String)
}

func TestParserWithMultipleSpans(t *testing.T) {
	expression := "span[tracetest.span.type=\"http\"], span[tracetest.span.type=\"grpc\"]"

	selector := parseExpressionOrFail(t, expression)

	assert.Equal(t, "tracetest.span.type", selector.SpanSelectors[0].Filters[0].Property)
	assert.Equal(t, "\"http\"", *selector.SpanSelectors[0].Filters[0].Value.String)
	assert.Nil(t, selector.SpanSelectors[0].ChildSelector)

	assert.Equal(t, "tracetest.span.type", selector.SpanSelectors[1].Filters[0].Property)
	assert.Equal(t, "\"grpc\"", *selector.SpanSelectors[1].Filters[0].Value.String)
	assert.Nil(t, selector.SpanSelectors[0].ChildSelector)
}

func TestParserWithMultipleSpansAndHierarchy(t *testing.T) {
	expression := "span[tracetest.span.type=\"http\"], span[service.name=\"Pokeshop\"] span[tracetest.span.type=\"grpc\"]"

	selector := parseExpressionOrFail(t, expression)

	assert.Equal(t, "tracetest.span.type", selector.SpanSelectors[0].Filters[0].Property)
	assert.Equal(t, "\"http\"", *selector.SpanSelectors[0].Filters[0].Value.String)
	assert.Nil(t, selector.SpanSelectors[0].ChildSelector)

	assert.Equal(t, "service.name", selector.SpanSelectors[1].Filters[0].Property)
	assert.Equal(t, "\"Pokeshop\"", *selector.SpanSelectors[1].Filters[0].Value.String)
	assert.NotNil(t, selector.SpanSelectors[1].ChildSelector)
	assert.Equal(t, "tracetest.span.type", selector.SpanSelectors[1].ChildSelector.Filters[0].Property)
	assert.Equal(t, "\"grpc\"", *selector.SpanSelectors[1].ChildSelector.Filters[0].Value.String)
}
