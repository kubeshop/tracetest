package filters_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/expression/filters"
	"github.com/kubeshop/tracetest/server/expression/types"
	"github.com/kubeshop/tracetest/server/expression/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestType(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          value.Value
		ExpectedOutput string
	}{
		{
			Name:           "should_return_array_for_empty_array",
			Input:          value.NewArray([]types.TypedValue{}),
			ExpectedOutput: "array",
		},
		{
			Name: "should_return_array_for_an_array",
			Input: value.NewArray([]types.TypedValue{
				types.GetTypedValue("1"),
				types.GetTypedValue("2"),
			}),
			ExpectedOutput: "array",
		},
		{
			Name:           "should_return_number_for_number",
			Input:          value.New(types.GetTypedValue("1")),
			ExpectedOutput: "number",
		},
		{
			Name:           "should_return_duration_for_duration",
			Input:          value.New(types.GetTypedValue("25ms")),
			ExpectedOutput: "duration",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := filters.Type(testCase.Input)
			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedOutput, output.String())
		})
	}
}
