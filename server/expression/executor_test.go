package expression_test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
)

type executorTestCase struct {
	Name       string
	Query      string
	ShouldPass bool

	AttributeDataStore      expression.DataStore
	MetaAttributesDataStore expression.DataStore
	VariableDataStore       expression.DataStore
}

func TestBasicExpressionExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:       "should_compare_equal_integers",
			Query:      `1 = 1`,
			ShouldPass: true,
		},
		{
			Name:       "should_fail_when_comparing_two_different_integers",
			Query:      `1 = 2`,
			ShouldPass: false,
		},
		{
			Name:       "should_detect_string_changes",
			Query:      `"matheus" != "jorge"`,
			ShouldPass: true,
		},
		{
			Name:       "should_be_able_to_detect_lower_numbers",
			Query:      `999 < 1000`,
			ShouldPass: true,
		},
		{
			Name:       "should_be_able_to_detect_lower_numbers",
			Query:      `13 > 12`,
			ShouldPass: true,
		},
	}

	executeTestCases(t, testCases)
}

func TestBasicOperationExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:       "should_allow_addition",
			Query:      "1 + 1 = 2",
			ShouldPass: true,
		},
		{
			Name:       "should_allow_subtraction",
			Query:      "8 - 3 > 0",
			ShouldPass: true,
		},
		{
			Name:       "should_allow_multiplication",
			Query:      "15 * 10 = 150",
			ShouldPass: true,
		},
		{
			Name:       "should_allow_multiplication",
			Query:      "8 / 2 = 4",
			ShouldPass: true,
		},
		{
			Name:       "should_add_durations",
			Query:      "100ms + 200ms = 300ms",
			ShouldPass: true,
		},
	}

	executeTestCases(t, testCases)
}

func TestAttributeExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:       "should_get_values_from_attributes",
			Query:      "attr:my_attribute = 42",
			ShouldPass: true,

			AttributeDataStore: expression.AttributeDataStore{
				Span: traces.Span{
					Attributes: traces.Attributes{
						"my_attribute": "42",
					},
				},
			},
		},
	}

	executeTestCases(t, testCases)
}

func TestVariableExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:       "should_get_values_from_variables",
			Query:      "var:my_variable = var:other_variable + 1",
			ShouldPass: true,

			VariableDataStore: expression.VariableDataStore(map[string]string{
				"my_variable":    "42",
				"other_variable": "41",
			}),
		},
	}

	executeTestCases(t, testCases)
}

func TestStringInterpolationExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:       "should_interpolate_simple_values",
			Query:      `attr:text = 'this run took ${"25ms"}'`,
			ShouldPass: true,
			AttributeDataStore: expression.AttributeDataStore{
				Span: traces.Span{
					Attributes: traces.Attributes{
						"text": "this run took 25ms",
					},
				},
			},
		},
		{
			Name:       "should_interpolate_multiple_values",
			Query:      `'${1} is a number, ${"dog"} is a string, and ${1ms + 1ns} is a duration' = '1 is a number, dog is a string, and 1000001 is a duration'`,
			ShouldPass: true,
		},
	}

	executeTestCases(t, testCases)
}

func TestFilterExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:       "should_extract_id_from_json",
			Query:      `attr:tracetest.response.body | json_path '.id' = 8`,
			ShouldPass: true,
			AttributeDataStore: expression.AttributeDataStore{
				Span: traces.Span{
					Attributes: traces.Attributes{
						"tracetest.response.body": `{"id": 8, "name": "john doe"}`,
					},
				},
			},
		},
		{
			Name:       "should_support_filters_with_arguments_containing_spaces",
			Query:      `'{ "name": "john", "age": 37 }' | regex_group '"age": (\d+)' = 37`,
			ShouldPass: true,
		},
		{
			Name:       "should_support_multiple_filters",
			Query:      `'{ "array": [{ "name": "john", "age": 37 }, { "name": "jonas", "age": 38 }]}' | regex_group '"age": (\d+)' | get_index 1 = 38`,
			ShouldPass: true,
		},
		{
			Name:       "should_count_array_input",
			Query:      `'{ "array": [1, 2, 3] }' | json_path '$.array[*]' | count = 3`,
			ShouldPass: true,
		},
		{
			Name:       "should_get_last_item_from_list",
			Query:      `'{ "array": [1, 2, 5] }' | json_path '$.array[*]' | get_index 'last' = 5`,
			ShouldPass: true,
		},
	}

	executeTestCases(t, testCases)
}

func TestMetaAttributesExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:               "should_support_count_meta_attribute",
			Query:              `attr:tracetest.selected_spans.count = 3`,
			ShouldPass:         true,
			AttributeDataStore: expression.AttributeDataStore{},
			MetaAttributesDataStore: expression.MetaAttributesDataStore{
				SelectedSpans: []traces.Span{
					// We don't have to fill the spans details to make the meta attribute work
					{},
					{},
					{},
				},
			},
		},
		{
			Name:               "should_support_count_meta_attribute",
			Query:              `"Selected matched ${attr:tracetest.selected_spans.count} spans" = "Selected matched 2 spans"`,
			ShouldPass:         true,
			AttributeDataStore: expression.AttributeDataStore{},
			MetaAttributesDataStore: expression.MetaAttributesDataStore{
				SelectedSpans: []traces.Span{
					{},
					{},
				},
			},
		},
	}

	executeTestCases(t, testCases)
}

func executeTestCases(t *testing.T, testCases []executorTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			executor := expression.NewExecutor(
				testCase.AttributeDataStore,
				testCase.MetaAttributesDataStore,
				testCase.VariableDataStore,
			)
			left, right, err := executor.ExecuteStatement(testCase.Query)
			debugMessage := fmt.Sprintf("left value: %s; right value: %s", left, right)
			if testCase.ShouldPass {
				assert.NoError(t, err, debugMessage)
			} else {
				assert.Error(t, err, debugMessage)
			}
		})
	}
}
