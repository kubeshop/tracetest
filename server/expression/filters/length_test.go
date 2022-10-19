package filters_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/expression/filters"
	"github.com/kubeshop/tracetest/server/expression/types"
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
			Name:           "should_get_zero_from_empty_list",
			Input:          filters.NewArrayValue([]types.TypedValue{}),
			ExpectedOutput: "0",
		},
		{
			Name: "should_get_one_from_single_item_list",
			Input: filters.NewArrayValue([]types.TypedValue{
				types.GetTypedValue("a"),
			}),
			ExpectedOutput: "1",
		},
		{
			Name: "should_count_multiple_item_list",
			Input: filters.NewArrayValue([]types.TypedValue{
				types.GetTypedValue("a"),
				types.GetTypedValue("b"),
				types.GetTypedValue("c"),
			}),
			ExpectedOutput: "3",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := filters.NewArrayValue(testCase.Input.Items)
			output, err := filters.Length(input)
			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedOutput, output.String())
		})
	}
}
