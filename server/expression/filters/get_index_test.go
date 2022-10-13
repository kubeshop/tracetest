package filters_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/expression/filters"
	"github.com/kubeshop/tracetest/server/expression/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIndex(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          filters.Value
		Index          string
		ExpectedOutput string
		ShouldFail     bool
	}{
		{
			Name: "should_fail_with_invalid_argument",
			Input: filters.NewArrayValue([]types.TypedValue{
				types.GetTypedValue("28"),
				types.GetTypedValue("29"),
				types.GetTypedValue("30"),
			}),
			Index:      `abc`,
			ShouldFail: true,
		},
		{
			Name:       "should_fail_with_index_out_of_boundaries",
			Input:      filters.NewValue(types.GetTypedValue("abc")),
			Index:      `1`,
			ShouldFail: true,
		},
		{
			Name: "should_get_correct_item",
			Input: filters.NewArrayValue([]types.TypedValue{
				types.GetTypedValue("abc"),
				types.GetTypedValue("def"),
				types.GetTypedValue("ghi"),
			}),
			Index:          `1`,
			ShouldFail:     false,
			ExpectedOutput: "def",
		},
		{
			Name: "should_get_last_item",
			Input: filters.NewArrayValue([]types.TypedValue{
				types.GetTypedValue("abc"),
				types.GetTypedValue("def"),
				types.GetTypedValue("ghi"),
			}),
			Index:          `last`,
			ShouldFail:     false,
			ExpectedOutput: "ghi",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := filters.NewArrayValue(testCase.Input)
			output, err := filters.GetIndex(input, testCase.Index)
			if testCase.ShouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedOutput, output.String())
			}
		})
	}
}
