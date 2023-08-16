package expression_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTokens(t *testing.T) {
	testCases := []struct {
		Statement      string
		ExpectedTokens []expression.Token
	}{
		{
			Statement: `"abc" = "abc"`,
			ExpectedTokens: []expression.Token{
				{Type: expression.StrType},
				{Type: expression.StrType},
			},
		},
		{
			Statement: `"abc" = 8`,
			ExpectedTokens: []expression.Token{
				{Type: expression.StrType},
				{Type: expression.NumberType},
			},
		},
		{
			Statement: `3 = [1, 2, 3]`,
			ExpectedTokens: []expression.Token{
				{Type: expression.NumberType},
				{Type: expression.ArrayType},
			},
		},
		{
			Statement: `env:url = "http://localhost"`,
			ExpectedTokens: []expression.Token{
				{Identifier: "url", Type: expression.EnvironmentType},
				{Type: expression.StrType},
			},
		},
		{
			Statement: `var:url = "http://localhost"`,
			ExpectedTokens: []expression.Token{
				{Identifier: "url", Type: expression.VariableType},
				{Type: expression.StrType},
			},
		},
		{
			Statement: `"the server is ${env:url}" = "http://localhost"`,
			ExpectedTokens: []expression.Token{
				{Type: expression.StrType},
				{Identifier: "url", Type: expression.EnvironmentType},
				{Type: expression.StrType},
			},
		},
		{
			Statement: `"the server is ${var:url}" = "http://localhost"`,
			ExpectedTokens: []expression.Token{
				{Type: expression.StrType},
				{Identifier: "url", Type: expression.VariableType},
				{Type: expression.StrType},
			},
		},
		{
			Statement: `"the url has ${env:url | count} characters" = "the url has 22 characters"`,
			ExpectedTokens: []expression.Token{
				{Type: expression.StrType},
				{Identifier: "url", Type: expression.EnvironmentType},
				{Identifier: "count", Type: expression.FunctionCallType},
				{Type: expression.StrType},
			},
		},
		{
			Statement: `"test ${env:names | get_index env:name_index}" = "John Doe"`,
			ExpectedTokens: []expression.Token{
				{Type: expression.StrType},
				{Identifier: "names", Type: expression.EnvironmentType},
				{Identifier: "get_index", Type: expression.FunctionCallType},
				{Identifier: "name_index", Type: expression.EnvironmentType},
				{Type: expression.StrType},
			},
		},
		{
			Statement: `"test ${var:names | get_index var:name_index}" = "John Doe"`,
			ExpectedTokens: []expression.Token{
				{Type: expression.StrType},
				{Identifier: "names", Type: expression.VariableType},
				{Identifier: "get_index", Type: expression.FunctionCallType},
				{Identifier: "name_index", Type: expression.VariableType},
				{Type: expression.StrType},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Statement, func(t *testing.T) {
			tokens, err := expression.GetTokens(testCase.Statement)
			require.NoError(t, err)

			assert.Equal(t, testCase.ExpectedTokens, tokens)
		})
	}
}
