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
			Name:  "should parse equals operator",
			Query: "http.status_code = 200",
			ExpectedOutput: parser.Assertion{
				Attribute: "http.status_code",
				Operator:  "=",
				Value:     "200",
			},
		},
		{
			Name:  "should parse less than operator",
			Query: "tracetest.span.duration < 100",
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.span.duration",
				Operator:  "<",
				Value:     "100",
			},
		},
		{
			Name:  "should parse less or equal to operator",
			Query: "tracetest.span.duration <= 100",
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.span.duration",
				Operator:  "<=",
				Value:     "100",
			},
		},
		{
			Name:  "should parse not equal operator",
			Query: `tracetest.span.type != "http"`,
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.span.type",
				Operator:  "!=",
				Value:     "http",
			},
		},
		{
			Name:  "should parse greater than operator",
			Query: `tracetest.span.duration > 0`,
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.span.duration",
				Operator:  ">",
				Value:     "0",
			},
		},
		{
			Name:  "should parse greater or equal to operator",
			Query: `tracetest.span.duration >= 0`,
			ExpectedOutput: parser.Assertion{
				Attribute: "tracetest.span.duration",
				Operator:  ">=",
				Value:     "0",
			},
		},
		{
			Name:  "should parse contains operator",
			Query: `db.statement contains "INSERT INTO"`,
			ExpectedOutput: parser.Assertion{
				Attribute: "db.statement",
				Operator:  "contains",
				Value:     "INSERT INTO",
			},
		},
		{
			Name:  "should parse string values",
			Query: `db.statement = "create"`,
			ExpectedOutput: parser.Assertion{
				Attribute: "db.statement",
				Operator:  "=",
				Value:     "create",
			},
		},
		{
			Name:  "should parse double values",
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
