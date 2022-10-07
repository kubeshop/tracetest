package expression_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
)

type executorTestCase struct {
	Name       string
	Query      string
	ShouldPass bool

	AttributeDataStore expression.AttributeDataStore
}

func TestBasicExpressions(t *testing.T) {
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

func TestBasicOperations(t *testing.T) {
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

func TestAttributes(t *testing.T) {
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

func TestStringInterpolations(t *testing.T) {
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

func executeTestCases(t *testing.T, testCases []executorTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			executor := expression.NewExecutor(testCase.AttributeDataStore)
			_, _, err := executor.ExecuteStatement(testCase.Query)
			if testCase.ShouldPass {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
