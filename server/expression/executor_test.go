package expression_test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type executorTestCase struct {
	Name                 string
	Query                string
	ShouldPass           bool
	ExpectedErrorMessage string

	AttributeDataStore      expression.DataStore
	MetaAttributesDataStore expression.DataStore
	EnvironmentDataStore    expression.DataStore
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
		{
			Name:       "should_enable_spaces_in_durations",
			Query:      `100ms < 200 ms`,
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
				Span: model.Span{
					Attributes: model.Attributes{
						"my_attribute": "42",
					},
				},
			},
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
				Span: model.Span{
					Attributes: model.Attributes{
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
				Span: model.Span{
					Attributes: model.Attributes{
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
			Query:      `'{ "array": [1, 2, 3] }' | json_path '$.array[*]' | length = 3`,
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
				SelectedSpans: []model.Span{
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
				SelectedSpans: []model.Span{
					{},
					{},
				},
			},
		},
	}

	executeTestCases(t, testCases)
}

func TestFunctionExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:       "should_generate_a_non_empty_first_name",
			Query:      `firstName() | length > 0`,
			ShouldPass: true,
		},
		{
			Name:       "should_generate_a_random_first_name",
			Query:      `firstName() != firstName()`,
			ShouldPass: true,
		},
		{
			Name:       "should_generate_a_random_int",
			Query:      `randomInt(0,10) <= 10`,
			ShouldPass: true,
		},
		{
			Name:       "should_generate_a_random_int_and_fail_comparison",
			Query:      `randomInt(10,20) < 10`,
			ShouldPass: false,
		},
	}

	executeTestCases(t, testCases)
}

func TestArrayExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:       "should_compare_two_empty_arrays",
			Query:      `[] = []`,
			ShouldPass: true,
		},
		{
			Name:       "should_compare_two_filled_arrays",
			Query:      `[1, 2, 3] = [1, 2, 3]`,
			ShouldPass: true,
		},
		{
			Name:       "should_fail_to_compare_different_arrays",
			Query:      `[1, 2, 3] = [3, 2, 1]`,
			ShouldPass: false,
		},
		{
			Name:       "should_be_able_to_run_filters_on_arrays",
			Query:      `[1, 2, 3] | length = 3`,
			ShouldPass: true,
		},
		{
			Name:       "should_be_able_to_chain_filters_on_arrays",
			Query:      `["this", "is", "sparta"] | get_index 2 | length = 6`,
			ShouldPass: true,
		},
		{
			Name:       "arrays_should_be_of_type_array",
			Query:      `["this", "is", "sparta"] | type = "array"`,
			ShouldPass: true,
		},
		{
			Name:       "incomplete_arrays_should_not_be_equal",
			Query:      `[1,2,3] = [1, 2]`,
			ShouldPass: false,
		},
		{
			Name:       "arrays_can_be_compared_with_other_arrays_generated_by_filters",
			Query:      `'{ "array": [{ "name": "john", "age": 37 }, { "name": "jonas", "age": 38 }]}' | json_path '$.array[*].age' = [37, 38]`,
			ShouldPass: true,
		},
		{
			Name:       "should_check_if_array_contains_value",
			Query:      `[31,35,39] contains 35`,
			ShouldPass: true,
		},
		{
			Name:       "should_check_if_array_contains_value_and_fail_if_not",
			Query:      `[31,35,39] contains 42`,
			ShouldPass: false,
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
			)
			left, right, err := executor.Statement(testCase.Query)
			debugMessage := fmt.Sprintf("left value: %s; right value: %s", left, right)
			if testCase.ShouldPass {
				assert.NoError(t, err, debugMessage)
			} else {
				assert.Error(t, err, debugMessage)
			}
		})
	}
}

func TestBasicStatementExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:       "should_parse_a_single_integer",
			Query:      `1`,
			ShouldPass: true,
		},
		{
			Name:       "should_parse_double_quoted_strings",
			Query:      `"matheus"`,
			ShouldPass: true,
		},
	}

	executeResolveStatementTestCases(t, testCases)
}

func TestResolveStatementAttributeExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:       "should_get_values_from_attributes",
			Query:      "attr:my_attribute",
			ShouldPass: true,

			AttributeDataStore: expression.AttributeDataStore{
				Span: model.Span{
					Attributes: model.Attributes{
						"my_attribute": "42",
					},
				},
			},
		},
	}

	executeResolveStatementTestCases(t, testCases)
}

func TestResolveStatementStringInterpolationExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:       "should_interpolate_simple_values",
			Query:      `'this run took ${"25ms"}'`,
			ShouldPass: true,
			AttributeDataStore: expression.AttributeDataStore{
				Span: model.Span{
					Attributes: model.Attributes{
						"text": "this run took 25ms",
					},
				},
			},
		},
		{
			Name:       "should_interpolate_multiple_values",
			Query:      `'${1} is a number, ${"dog"} is a string, and ${1ms + 1ns} is a duration'`,
			ShouldPass: true,
		},
	}

	executeResolveStatementTestCases(t, testCases)
}

func TestResolveStatementFilterExecution(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:       "should_extract_id_from_json",
			Query:      `attr:tracetest.response.body`,
			ShouldPass: true,
			AttributeDataStore: expression.AttributeDataStore{
				Span: model.Span{
					Attributes: model.Attributes{
						"tracetest.response.body": `{"id": 8, "name": "john doe"}`,
					},
				},
			},
		},
		{
			Name:       "should_support_filters_with_arguments_containing_spaces",
			Query:      `'{ "name": "john", "age": 37 }' | regex_group '"age": (\d+)'`,
			ShouldPass: true,
		},
		{
			Name:       "should_support_multiple_filters",
			Query:      `'{ "array": [{ "name": "john", "age": 37 }, { "name": "jonas", "age": 38 }]}' | regex_group '"age": (\d+)' | get_index 1`,
			ShouldPass: true,
		},
		{
			Name:       "should_count_array_input",
			Query:      `'{ "array": [1, 2, 3] }' | json_path '$.array[*]' | length`,
			ShouldPass: true,
		},
		{
			Name:       "should_get_last_item_from_list",
			Query:      `'{ "array": [1, 2, 5] }' | json_path '$.array[*]' | get_index 'last'`,
			ShouldPass: true,
		},
	}

	executeResolveStatementTestCases(t, testCases)
}

func TestFailureCases(t *testing.T) {
	testCases := []executorTestCase{
		{
			Name:                 "should_report_missing_environment_variable",
			Query:                `env:test = "abc"`,
			ShouldPass:           false,
			ExpectedErrorMessage: `resolution error: environment variable "test" not found`,

			EnvironmentDataStore: expression.EnvironmentDataStore{
				Values: []model.EnvironmentValue{},
			},
		},
		{
			Name:                 "should_report_missing_attribute",
			Query:                `attr:my_attribute = "abc"`,
			ShouldPass:           false,
			ExpectedErrorMessage: `resolution error: attribute "my_attribute" not found`,

			AttributeDataStore: expression.AttributeDataStore{
				Span: model.Span{
					Attributes: model.Attributes{
						"attr1": "1",
						"attr2": "2",
					},
				},
			},
		},
		{
			Name:                 "should_report_missing_filter",
			Query:                `"value" | missingFilter = "abc"`,
			ShouldPass:           false,
			ExpectedErrorMessage: `resolution error: filter "missingFilter" not found`,
		},
		{
			Name:                 "should_report_problem_resolving_array_item",
			Query:                `["value", env:test, "anotherValue"] | get_index 0`,
			ShouldPass:           false,
			ExpectedErrorMessage: `resolution error: at index 1 of array: environment variable "test" not found`,

			EnvironmentDataStore: expression.EnvironmentDataStore{
				Values: []model.EnvironmentValue{},
			},
		},
	}

	executeResolveStatementTestCases(t, testCases)
}

func executeResolveStatementTestCases(t *testing.T, testCases []executorTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			executor := expression.NewExecutor(
				testCase.AttributeDataStore,
				testCase.MetaAttributesDataStore,
				testCase.EnvironmentDataStore,
			)
			left, err := executor.ResolveStatement(testCase.Query)
			debugMessage := fmt.Sprintf("left value: %s", left)
			if testCase.ShouldPass {
				assert.NoError(t, err, debugMessage)
			} else {
				require.Error(t, err, debugMessage)
				if testCase.ExpectedErrorMessage != "" {
					assert.Equal(t, testCase.ExpectedErrorMessage, err.Error())
					// all validation erros should be ErrExpressionResolution errors
					assert.ErrorIs(t, err, expression.ErrExpressionResolution)
				}
			}
		})
	}
}
