package filters_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/expression/filters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegex(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          string
		Regex          string
		ExpectedOutput string
	}{
		{
			Name:           "should_extract_unique_field_from_JSON",
			Input:          `{ "id": 38, "name": "Tracetest" }`,
			Regex:          `"id": \d+`,
			ExpectedOutput: `"id": 38`,
		},
		{
			Name:           "should_extract_unique_field_from_JSON",
			Input:          `[{ "id": 38, "name": "Tracetest" }, { "id": 39, "name": "Kusk" }]`,
			Regex:          `"id": \d+`,
			ExpectedOutput: `["id": 38,"id": 39]`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := filters.Regex(testCase.Input, testCase.Regex)
			require.NoError(t, err)

			assert.Equal(t, testCase.ExpectedOutput, output)
		})
	}
}
