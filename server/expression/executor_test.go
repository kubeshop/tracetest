package expression_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/stretchr/testify/assert"
)

type executorTestCase struct {
	Name       string
	Query      string
	ShouldPass bool
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

func executeTestCases(t *testing.T, testCases []executorTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := expression.ExecuteStatement(testCase.Query)
			if testCase.ShouldPass {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
