package functions_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/kubeshop/tracetest/server/executor/functions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		Name        string
		Input       string
		ShouldFail  bool
		ShouldMatch string
	}{
		{
			Name:        "should_parse_uuid_function",
			Input:       "uuid()",
			ShouldMatch: `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, // uuid format
		},
		{
			Name:        "should_parse_first_name",
			Input:       "firstName()",
			ShouldMatch: `^[a-zA-Z]+$`,
		},
		{
			Name:        "should_parse_last_name",
			Input:       "lastName()",
			ShouldMatch: `^[a-zA-Z]+$`,
		},
		{
			Name:        "should_parse_full_name",
			Input:       "fullName()",
			ShouldMatch: `^[a-zA-Z]+ [a-zA-Z]+$`,
		},
		{
			Name:        "should_parse_email",
			Input:       "email()",
			ShouldMatch: `^(.+)@(.+)$`,
		},
		{
			Name:        "should_parse_phone",
			Input:       "phone()",
			ShouldMatch: `^\d{10}$`,
		},
		{
			Name:        "should_parse_credit_card",
			Input:       "creditCard()",
			ShouldMatch: `^[0-9]{14,19}$`,
		},
		{
			Name:        "should_parse_credit_card_ccv",
			Input:       "creditCardCvv()",
			ShouldMatch: `^[0-9]{3}$`,
		},
		{
			Name:        "should_parse_random_int_one_digit",
			Input:       "randomInt(1, 9)",
			ShouldMatch: `^[0-9]{1}$`,
		},
		{
			Name:        "should_parse_random_int_two_digit",
			Input:       "randomInt(10, 99)",
			ShouldMatch: `^[0-9]{2}$`,
		},
		{
			Name:        "should_parse_random_int_three_digit",
			Input:       "randomInt(100, 999)",
			ShouldMatch: `^[0-9]{3}$`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			function, args, err := functions.ParseFunction(testCase.Input)
			if testCase.ShouldFail {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				result, err := function.Invoke(args...)

				assert.NoError(t, err)
				assert.NotEmpty(t, result)

				regex, err := regexp.Compile(testCase.ShouldMatch)
				require.NoError(t, err, "provided regex is invalid")
				assert.True(t, regex.Match([]byte(result)), fmt.Sprintf("%s does not match provided regex", result))
			}
		})
	}
}
