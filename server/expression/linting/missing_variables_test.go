package linting_test

import (
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/expression/linting"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/stretchr/testify/assert"
)

type HTTPAuth struct {
	Token string `expr_enabled:"true"`
}

type Assertion struct {
	Name    string
	Queries []string `stmt_enabled:"true"`
}

type HTTPRequest struct {
	URL        string `expr_enabled:"true"`
	Method     string `expr_enabled:"true"`
	Auth       HTTPAuth
	Assertions []Assertion `stmt_enabled:"true"`
}

func TestMissingVariableDetection(t *testing.T) {
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
				URL:    `${env:SERVER_URL}:${var:PORT}`,
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
				URL:    `${env:SERVER_URL}:${var:PORT}`,
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
				URL:    `${env:SERVER_URL}:${env:PORT}`,
				Method: "GET",
				Auth: HTTPAuth{
					Token: "${env:TOKEN}",
				},
			},
			expectedMissingVariables: []string{"TOKEN"},
		},
		{
			name:               "missing_variables_in_inner_struct",
			availableVariables: []string{"SERVER_URL", "PORT"},
			object: HTTPRequest{
				URL:    `${env:SERVER_URL}:${var:PORT}`,
				Method: "GET",
				Auth: HTTPAuth{
					Token: "${env:TOKEN}",
				},
			},
			expectedMissingVariables: []string{"TOKEN"},
		},
		{
			name:               "missing_variable_in_statements",
			availableVariables: []string{"SERVER_URL", "PORT"},
			object: HTTPRequest{
				URL:    `${env:SERVER_URL}:${env:PORT}`,
				Method: "GET",
				Auth: HTTPAuth{
					Token: "abc",
				},
				Assertions: []Assertion{
					{Name: "test", Queries: []string{"env:ABC = var:ABC2"}},
					{Name: "test2", Queries: []string{"env:CDE = var:CDE"}},
				},
			},
			expectedMissingVariables: []string{"ABC", "ABC2", "CDE"},
		},
		{
			name:                     "should_work_with_time_objects",
			availableVariables:       []string{},
			object:                   time.Now(),
			expectedMissingVariables: []string{},
		},
		{
			name:               "should_work_with_pointers",
			availableVariables: []string{"SERVER_URL", "PORT"},
			object: &HTTPRequest{
				URL:    `"${env:SERVER_URL}:${env:PORT}"`,
				Method: "GET",
				Auth: HTTPAuth{
					Token: "abc",
				},
				Assertions: []Assertion{
					{Name: "test", Queries: []string{"env:ABC = env:ABC2"}},
					{Name: "test2", Queries: []string{"env:CDE = var:CDE"}},
				},
			},
			expectedMissingVariables: []string{"ABC", "ABC2", "CDE"},
		},
		{
			name:               "should_detect_missing_variables_in_test_http_body",
			availableVariables: []string{},
			object: test.Test{
				Trigger: trigger.Trigger{
					Type: trigger.TriggerTypeHTTP,
					HTTP: &trigger.HTTPRequest{
						Body: `{"id": ${env:pokemonId},"name": "${var:pokemonName}"}`,
					},
				},
			},
			expectedMissingVariables: []string{"pokemonId", "pokemonName"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			missingVariables := linting.DetectMissingVariables(testCase.object, testCase.availableVariables)
			assert.Equal(t, testCase.expectedMissingVariables, missingVariables)
		})
	}
}

func TestMissingVariableDetectionInOrderedMap(t *testing.T) {
	testCases := []struct {
		name                     string
		availableVariables       []string
		input                    maps.Ordered[string, HTTPRequest]
		expectedMissingVariables []string
	}{
		{
			name: "should_be_able_to_scan_ordered_map",
			input: maps.Ordered[string, HTTPRequest]{}.MustAdd("abc", HTTPRequest{
				URL:    "${env:URL}",
				Method: "GET",
				Auth: HTTPAuth{
					Token: "${env:TOKEN}",
				},
				Assertions: []Assertion{
					{Name: "abc", Queries: []string{"env:ABC = env:ABC"}},
				},
			}),
			availableVariables:       []string{},
			expectedMissingVariables: []string{"URL", "TOKEN", "ABC"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			missingVariables := linting.DetectMissingVariables(testCase.input, testCase.availableVariables)
			assert.Equal(t, testCase.expectedMissingVariables, missingVariables)
		})
	}
}
