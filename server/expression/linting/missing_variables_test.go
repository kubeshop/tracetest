package linting_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/expression/linting"
	"github.com/stretchr/testify/assert"
)

func TestMissingVariableDetection(t *testing.T) {
	type HTTPAuth struct {
		Token string `expr_enabled:"true"`
	}

	type Assertion struct {
		Name  string
		Query string `stmt_enabled:"true"`
	}

	type HTTPRequest struct {
		URL        string `expr_enabled:"true"`
		Method     string `expr_enabled:"true"`
		Auth       HTTPAuth
		Assertions []Assertion
	}

	testCases := []struct {
		name                     string
		availableVariables       []string
		object                   interface{}
		expectedMissingVariables []string
	}{
		{
			name:               "no_missing_variables_if_no_variables",
			availableVariables: []string{"SERVER_URL", "PORT"},
			object: HTTPRequest{
				URL:    "http://localhost",
				Method: "GET",
				Auth: HTTPAuth{
					Token: "abc",
				},
			},
			expectedMissingVariables: []string{},
		},
		{
			name:               "no_missing_variables_if_variable_exists",
			availableVariables: []string{"SERVER_URL", "PORT", "TOKEN"},
			object: HTTPRequest{
				URL:    `"${env:SERVER_URL}:${PORT}"`,
				Method: "GET",
				Auth: HTTPAuth{
					Token: "abc",
				},
			},
			expectedMissingVariables: []string{},
		},
		{
			name:               "missing_variables_if_variable_doesnt_exists",
			availableVariables: []string{"SERVER_URL"},
			object: HTTPRequest{
				URL:    `"${env:SERVER_URL}:${env:PORT}"`,
				Method: "GET",
				Auth: HTTPAuth{
					Token: "abc",
				},
			},
			expectedMissingVariables: []string{"PORT"},
		},
		{
			name:               "missing_variables_if_inner_variable_doesnt_exists",
			availableVariables: []string{"SERVER_URL", "PORT"},
			object: HTTPRequest{
				URL:    `"${env:SERVER_URL}:${env:PORT}"`,
				Method: "GET",
				Auth: HTTPAuth{
					Token: "env:TOKEN",
				},
			},
			expectedMissingVariables: []string{"TOKEN"},
		},
		{
			name:               "missing_variables_in_inner_struct",
			availableVariables: []string{"SERVER_URL", "PORT"},
			object: HTTPRequest{
				URL:    `"${env:SERVER_URL}:${env:PORT}"`,
				Method: "GET",
				Auth: HTTPAuth{
					Token: "env:TOKEN",
				},
			},
			expectedMissingVariables: []string{"TOKEN"},
		},
		{
			name:               "missing_variable_in_statements",
			availableVariables: []string{"SERVER_URL", "PORT"},
			object: HTTPRequest{
				URL:    `"${env:SERVER_URL}:${env:PORT}"`,
				Method: "GET",
				Auth: HTTPAuth{
					Token: "abc",
				},
				Assertions: []Assertion{
					{Name: "test", Query: "env:ABC = env:ABC"},
					{Name: "test2", Query: "env:CDE = env:CDE"},
				},
			},
			expectedMissingVariables: []string{"ABC", "CDE"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			missingVariables := linting.DetectMissingVariables(testCase.object, testCase.availableVariables)
			assert.Equal(t, testCase.expectedMissingVariables, missingVariables)
		})
	}
}
