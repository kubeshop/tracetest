package traces_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
)

func TestTimeConvertion(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          int
		ExpectedOutput string
	}{
		{
			Name:           "should convert nanoseconds",
			Input:          38,
			ExpectedOutput: "38ns",
		},
		{
			Name:           "should convert microseconds",
			Input:          93000,
			ExpectedOutput: "93μs",
		},
		{
			Name:           "should not render decimals when time is in microseconds",
			Input:          93400,
			ExpectedOutput: "93μs",
		},
		{
			Name:           "should convert milliseconds",
			Input:          29000000,
			ExpectedOutput: "29ms",
		},
		{
			Name:           "should not render decimals when time is in milliseconds",
			Input:          27356000,
			ExpectedOutput: "27ms",
		},
		{
			Name:           "should convert seconds",
			Input:          1000000000,
			ExpectedOutput: "1s",
		},
		{
			Name:           "should render decimals when time is in seconds",
			Input:          1200000000,
			ExpectedOutput: "1.2s",
		},
		{
			Name:           "should convert minutes",
			Input:          1620000000000, // 27m
			ExpectedOutput: "27m",
		},
		{
			Name:           "should render decimals when time is in minutes",
			Input:          1650000000000, // 27m and 30 seconds
			ExpectedOutput: "27.5m",
		},
		{
			Name:           "should convert hours",
			Input:          7200000000000, // 2 hours
			ExpectedOutput: "2h",
		},
		{
			Name:           "should convert hours",
			Input:          9000000000000, // 2 hours and 30 minutes
			ExpectedOutput: "2.5h",
		},
		{
			Name:           "hour should be the largest unit",
			Input:          38880000000000000, // 450 days
			ExpectedOutput: "10800h",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output := traces.ConvertNanoSecondsIntoProperTimeUnit(testCase.Input)
			assert.Equal(t, testCase.ExpectedOutput, output)
		})
	}
}
