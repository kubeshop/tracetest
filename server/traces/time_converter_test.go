package traces_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
)

func TestConversionBetweenNanoSecondsToDuration(t *testing.T) {
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

func TestConversionBetweenDurationToNanoSeconds(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          string
		ExpectedOutput int
	}{
		{
			Name:           "should convert ns",
			Input:          "120ns",
			ExpectedOutput: 120,
		},
		{
			Name:           "should floor ns with decimals",
			Input:          "120.5ns",
			ExpectedOutput: 120,
		},
		{
			Name:           "should convert μs",
			Input:          "35μs",
			ExpectedOutput: 35000,
		},
		{
			Name:           "should convert μs with decimal",
			Input:          "35.8μs",
			ExpectedOutput: 35800,
		},
		{
			Name:           "should convert ms",
			Input:          "68ms",
			ExpectedOutput: 68000000,
		},
		{
			Name:           "should convert ms with decimal",
			Input:          "68.35ms",
			ExpectedOutput: 68350000,
		},
		{
			Name:           "should convert s",
			Input:          "1s",
			ExpectedOutput: 1000000000,
		},
		{
			Name:           "should convert s with decimal",
			Input:          "1.23s",
			ExpectedOutput: 1230000000,
		},
		{
			Name:           "should convert m",
			Input:          "1m",
			ExpectedOutput: 60000000000,
		},
		{
			Name:           "should convert m with decimal",
			Input:          "1.5m",
			ExpectedOutput: 90000000000,
		},
		{
			Name:           "should convert h",
			Input:          "1h",
			ExpectedOutput: 3600000000000,
		},
		{
			Name:           "should convert h with decimal",
			Input:          "2.5m",
			ExpectedOutput: 8000000000000,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output := traces.ConvertTimeFieldIntoNanoSeconds(testCase.Input)
			assert.Equal(t, testCase.ExpectedOutput, output)
		})
	}
}
