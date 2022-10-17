package filters_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/expression/filters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLength(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          filters.Value
		ExpectedOutput string
	}{
		{
			Name:           "should_return_zero_for_empty_string",
			Input:          filters.NewValueFromString(""),
			ExpectedOutput: "0",
		},
		{
			Name:           "should_return_length_for_nonempty_string",
			Input:          filters.NewValueFromString("result should be 19"),
			ExpectedOutput: "19",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := filters.NewArrayValue(testCase.Input)
			output, err := filters.Length(input)
			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedOutput, output.String())
		})
	}
}
