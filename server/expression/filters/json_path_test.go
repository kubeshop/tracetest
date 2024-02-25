package filters_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/expression/filters"
	"github.com/kubeshop/tracetest/server/expression/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONPath(t *testing.T) {
	testCases := []struct {
		Name           string
		JSON           string
		Query          string
		ExpectedOutput string
	}{
		{
			Name:           "should_extract_unique_field_from_object",
			JSON:           `{ "id": 38, "name": "Tracetest" }`,
			Query:          `.id`,
			ExpectedOutput: `38`,
		},
		{
			Name:           "should_extract_unique_field_from_array",
			JSON:           `{ "array": [ { "id": 38, "name": "Tracetest" } ] }`,
			Query:          `$.array[0].id`,
			ExpectedOutput: `38`,
		},
		{
			Name:           "should_extract_multiple_values_from_array",
			JSON:           `{ "array": [ { "id": 38, "name": "Tracetest" }, {"id": 39, "name": "Kusk"} ] }`,
			Query:          `$.array[*].id`,
			ExpectedOutput: `[38, 39]`,
		},
		{
			Name:           "should_extract_multiple_fields_from_array",
			JSON:           `{ "array": [ { "id": 38, "name": "Tracetest" }, {"id": 39, "name": "Kusk"} ] }`,
			Query:          `$.array[*]..['id', 'name']`,
			ExpectedOutput: `[38, "Tracetest", 39, "Kusk"]`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := value.NewFromString(testCase.JSON)
			output, err := filters.JSON_path(input, testCase.Query)
			require.NoError(t, err)

			assert.Equal(t, testCase.ExpectedOutput, output.String())
		})
	}
}
