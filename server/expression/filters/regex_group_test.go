package filters_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/expression/filters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexGroup(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          string
		Regex          string
		ExpectedOutput string
	}{
		{
			Name:           "should_be_able_to_extract_one_group",
			Input:          `{ "id": 38, "name": "Tracetest" }`,
			Regex:          `"id": (\d+)`,
			ExpectedOutput: `38`,
		},
		{
			Name:           "should_be_able_to_extract_one_group_multiple_times",
			Input:          `[{ "id": 38, "name": "Tracetest" }, { "id": 39, "name": "Kusk" }]`,
			Regex:          `"id": (\d+)`,
			ExpectedOutput: `[38, 39]`,
		},
		{
			Name:           "should_be_able_to_extract_multiple_groups_multiple_times",
			Input:          `[{ "id": 38, "name": "Tracetest" }, { "id": 39, "name": "Kusk" }]`,
			Regex:          `"id": (\d+), "name": "(\w+)"`,
			ExpectedOutput: `[38, "Tracetest", 39, "Kusk"]`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := filters.NewValueFromString(testCase.Input)
			output, err := filters.RegexGroup(input, testCase.Regex)
			require.NoError(t, err)

			assert.Equal(t, testCase.ExpectedOutput, output.String())
		})
	}
}
