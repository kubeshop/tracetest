package test_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/pkg/timing"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/variableset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunExecutionTime(t *testing.T) {
	cases := []struct {
		name     string
		run      test.Run
		now      time.Time
		expected int
	}{
		{
			name: "CompletedOk",
			run: test.Run{
				CreatedAt:   timing.MustParse("2022-01-25T12:45:33.100000000Z"),
				CompletedAt: timing.MustParse("2022-01-25T12:45:36.400000000Z"),
			},
			expected: 4,
		},
		{
			name: "LessThan1Sec",
			run: test.Run{
				CreatedAt:   timing.MustParse("2022-01-25T12:45:33.100000000Z"),
				CompletedAt: timing.MustParse("2022-01-25T12:45:33.400000000Z"),
			},
			expected: 1,
		},
		{
			name: "StillRunning",
			run: test.Run{
				CreatedAt: timing.MustParse("2022-01-25T12:45:33.100000000Z"),
			},
			now:      timing.MustParse("2022-01-25T12:45:34.300000000Z"),
			expected: 2,
		},
		{
			name: "ZeroedDate",
			run: test.Run{
				CreatedAt:   timing.MustParse("2022-01-25T12:45:33.100000000Z"),
				CompletedAt: time.Unix(0, 0),
			},
			now:      timing.MustParse("2022-01-25T12:45:34.300000000Z"),
			expected: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			now := timing.Now
			if c.now.Unix() > 0 {
				timing.Now = func() time.Time {
					return c.now
				}
			}

			assert.Equal(t, c.expected, c.run.ExecutionTime())
			timing.Now = now
		})
	}
}

func TestRunTriggerTime(t *testing.T) {
	cases := []struct {
		name     string
		run      test.Run
		now      time.Time
		expected int
	}{
		{
			name: "CompletedOk",
			run: test.Run{
				ServiceTriggeredAt:        timing.MustParse("2022-01-25T12:45:33.100000000Z"),
				ServiceTriggerCompletedAt: timing.MustParse("2022-01-25T12:45:36.400000000Z"),
			},
			expected: 3300,
		},
		{
			name: "LessThan1Sec",
			run: test.Run{
				ServiceTriggeredAt:        timing.MustParse("2022-01-25T12:45:33.100000000Z"),
				ServiceTriggerCompletedAt: timing.MustParse("2022-01-25T12:45:33.400000000Z"),
			},
			expected: 300,
		},
		{
			name: "StillRunning",
			run: test.Run{
				ServiceTriggeredAt: timing.MustParse("2022-01-25T12:45:33.100000000Z"),
			},
			now:      timing.MustParse("2022-01-25T12:45:34.300000000Z"),
			expected: 1200,
		},
		{
			name: "ZeroedDate",
			run: test.Run{
				ServiceTriggeredAt:        timing.MustParse("2022-01-25T12:45:33.100000000Z"),
				ServiceTriggerCompletedAt: time.Unix(0, 0),
			},
			now:      timing.MustParse("2022-01-25T12:45:34.300000000Z"),
			expected: 1200,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			now := timing.Now
			if c.now.Unix() > 0 {
				timing.Now = func() time.Time {
					return c.now
				}
			}

			assert.Equal(t, c.expected, c.run.TriggerTime())
			timing.Now = now
		})
	}
}

func TestRunRequiredGates(t *testing.T) {
	t.Run("NoGates", func(t *testing.T) {
		run := test.Run{}.ConfigureRequiredGates([]testrunner.RequiredGate{})
		linterResult := analyzer.LinterResult{
			Passed:       false,
			MinimumScore: 80,
			Score:        60,
		}

		run = run.SuccessfulLinterExecution(linterResult)

		assert.Equal(t, true, run.RequiredGatesResult.Passed)
		assert.Len(t, run.RequiredGatesResult.Failed, 2)

		failed := []testrunner.RequiredGate{testrunner.RequiredGateAnalyzerRules, testrunner.RequiredGateAnalyzerScore}
		assert.Equal(t, failed, run.RequiredGatesResult.Failed)
	})

	t.Run("AnalyzerGates", func(t *testing.T) {
		gates := []testrunner.RequiredGate{testrunner.RequiredGateAnalyzerRules, testrunner.RequiredGateAnalyzerScore}
		run := test.Run{}.ConfigureRequiredGates(gates)
		linterResult := analyzer.LinterResult{
			Passed:       true,
			MinimumScore: 80,
			Score:        60,
		}

		run = run.SuccessfulLinterExecution(linterResult)

		assert.Len(t, run.RequiredGatesResult.Required, 2)
		assert.Equal(t, false, run.RequiredGatesResult.Passed)
		assert.Len(t, run.RequiredGatesResult.Failed, 1)
		failed := []testrunner.RequiredGate{testrunner.RequiredGateAnalyzerScore}
		assert.Equal(t, failed, run.RequiredGatesResult.Failed)
	})

	t.Run("TestSpecGate", func(t *testing.T) {
		gates := []testrunner.RequiredGate{testrunner.RequiredGateTestSpecs}
		outputs := maps.Ordered[string, test.RunOutput]{}
		variableSet := variableset.VariableSet{}
		res := maps.Ordered[test.SpanQuery, []test.AssertionResult]{}
		allPassed := false
		run := test.Run{}.ConfigureRequiredGates(gates)

		run = run.SuccessfullyAsserted(outputs, variableSet, res, allPassed)

		assert.Len(t, run.RequiredGatesResult.Required, 1)
		assert.Equal(t, false, run.RequiredGatesResult.Passed)
		assert.Len(t, run.RequiredGatesResult.Failed, 1)
		failed := []testrunner.RequiredGate{testrunner.RequiredGateTestSpecs}
		assert.Equal(t, failed, run.RequiredGatesResult.Failed)
	})

	t.Run("GenerateRequiredGateResult", func(t *testing.T) {
		gates := []testrunner.RequiredGate{testrunner.RequiredGateTestSpecs, testrunner.RequiredGateAnalyzerRules, testrunner.RequiredGateAnalyzerScore}
		result := test.Run{
			Results: &test.RunResults{
				AllPassed: false,
			},
			Linter: analyzer.LinterResult{
				Passed:       false,
				MinimumScore: 80,
				Score:        60,
			},
		}.GenerateRequiredGateResult(gates)

		fmt.Println(result)

		assert.Len(t, result.Required, 3)
		assert.Equal(t, false, result.Passed)
		assert.Len(t, result.Failed, 3)
		assert.Equal(t, gates, result.Failed)
	})
}

func TestRunOutput(t *testing.T) {
	jsonContent := []byte(`[{"Key": "MY_OUTPUT", "Value": {"Name": "MY_OUTPUT", "Error": {}, "Value": "", "SpanID": "", "Resolved": false}}]`)

	var output test.Run
	err := json.Unmarshal(jsonContent, &output.Outputs)
	require.NoError(t, err)
}

func TestRunOutputJSON(t *testing.T) {
	jsonContent := []byte(`{"Name": "MY_OUTPUT", "Error": {}, "Value": "1", "SpanID": "2", "Resolved": true}`)

	var output test.RunOutput
	err := json.Unmarshal(jsonContent, &output)
	require.NoError(t, err)

	assert.Equal(t, "MY_OUTPUT", output.Name)
	assert.Equal(t, "1", output.Value)
	assert.Equal(t, "2", output.SpanID)
	assert.Equal(t, true, output.Resolved)
	assert.Nil(t, output.Error)
}

func TestRunOutputMarshalJSON(t *testing.T) {
	var runOutput = test.RunOutput{
		Name:     "MY_OUTPUT",
		Value:    "1",
		SpanID:   "2",
		Resolved: true,
		Error:    errors.New("not empty"),
	}
	jsonContent, _ := json.Marshal(runOutput)

	var output test.RunOutput
	err := json.Unmarshal(jsonContent, &output)
	require.NoError(t, err)

	assert.Equal(t, "MY_OUTPUT", output.Name)
	assert.Equal(t, "1", output.Value)
	assert.Equal(t, "2", output.SpanID)
	assert.Equal(t, true, output.Resolved)
	assert.Equal(t, "not empty", output.Error.Error())
}
