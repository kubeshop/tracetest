package parser_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/conversion/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAssertion(t *testing.T) {
	testCases := []struct {
		Name           string
		Query          string
		ExpectedOutput parser.Assertion
	}{
		{
			Name:  "should_parse_equals_operator",
			Query: "http.status_code = 200",
			ExpectedOutput: parser.Assertion{
				Attribute: "http.status_code",
				Operator:  "=",
				Value:     "200",
			},
		},
		{
			Name:  "should_parse_less_than_operator",
			Query: "tracetest.span.duration < 100",
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.span.duration",
				Operator:  "<",
				Value:     "100",
			},
		},
		{
			Name:  "should_parse_less_or_equal_to_operator",
			Query: "tracetest.span.duration <= 100",
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.span.duration",
				Operator:  "<=",
				Value:     "100",
			},
		},
		{
			Name:  "should_parse_not_equal_operator",
			Query: `tracetest.span.type != "http"`,
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.span.type",
				Operator:  "!=",
				Value:     "http",
			},
		},
		{
			Name:  "should_parse_greater_than_operator",
			Query: `tracetest.span.duration > 0`,
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.span.duration",
				Operator:  ">",
				Value:     "0",
			},
		},
		{
			Name:  "should_parse_greater_or_equal_to_operator",
			Query: `tracetest.span.duration >= 0`,
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.span.duration",
				Operator:  ">=",
				Value:     "0",
			},
		},
		{
			Name:  "should_parse_contains_operator",
			Query: `db.statement contains "INSERT INTO"`,
			ExpectedOutput: parser.Assertion{
				Attribute: "db.statement",
				Operator:  "contains",
				Value:     "INSERT INTO",
			},
		},
		{
			Name:  "should_parse_quoted_string_wrapped_on_single_quotes",
			Query: `tracetest.response.body contains '"id":"${TEST_ID}"'`,
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.response.body",
				Operator:  "contains",
				Value:     `"id":"${TEST_ID}"`,
			},
		},
		{
			Name:  "should_parse_escaped_quoted_string",
			Query: `tracetest.response.body contains "\"id\":\"${TEST_ID}\""`,
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.response.body",
				Operator:  "contains",
				Value:     `"id":"${TEST_ID}"`,
			},
		},
		{
			Name:  "should_parse_single_quoted_string_wrapped_on_double_quotes",
			Query: `tracetest.response.body contains "'single quoted value'"`,
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.response.body",
				Operator:  "contains",
				Value:     `'single quoted value'`,
			},
		},
		{
			Name:  "should_parse_escaped_single_quoted_string",
			Query: `tracetest.response.body contains '\'id\':\'${TEST_ID}\''`,
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.response.body",
				Operator:  "contains",
				Value:     `'id':'${TEST_ID}'`,
			},
		},
		{
			Name:  "should_parse_string_values",
			Query: `db.statement = "create"`,
			ExpectedOutput: parser.Assertion{
				Attribute: "db.statement",
				Operator:  "=",
				Value:     "create",
			},
		},
		{
			Name:  "should_parse_double_values",
			Query: `custom.item.value = 199.99`,
			ExpectedOutput: parser.Assertion{
				Attribute: "custom.item.value",
				Operator:  "=",
				Value:     "199.99",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := parser.ParseAssertion(testCase.Query)

			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedOutput, output)

		})
	}
}
