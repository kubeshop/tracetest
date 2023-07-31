package test_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/variableset"
	"github.com/stretchr/testify/assert"
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
				CreatedAt:   time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				CompletedAt: time.Date(2022, 01, 25, 12, 45, 36, int(400*time.Millisecond), time.UTC),
			},
			expected: 4,
		},
		{
			name: "LessThan1Sec",
			run: test.Run{
				CreatedAt:   time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				CompletedAt: time.Date(2022, 01, 25, 12, 45, 33, int(400*time.Millisecond), time.UTC),
			},
			expected: 1,
		},
		{
			name: "StillRunning",
			run: test.Run{
				CreatedAt: time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, int(300*time.Millisecond), time.UTC),
			expected: 2,
		},
		{
			name: "ZeroedDate",
			run: test.Run{
				CreatedAt:   time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				CompletedAt: time.Unix(0, 0),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, int(300*time.Millisecond), time.UTC),
			expected: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			now := test.Now
			if c.now.Unix() > 0 {
				test.Now = func() time.Time {
					return c.now
				}
			}

			assert.Equal(t, c.expected, c.run.ExecutionTime())
			test.Now = now
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
				ServiceTriggeredAt:        time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				ServiceTriggerCompletedAt: time.Date(2022, 01, 25, 12, 45, 36, int(400*time.Millisecond), time.UTC),
			},
			expected: 3300,
		},
		{
			name: "LessThan1Sec",
			run: test.Run{
				ServiceTriggeredAt:        time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				ServiceTriggerCompletedAt: time.Date(2022, 01, 25, 12, 45, 33, int(400*time.Millisecond), time.UTC),
			},
			expected: 300,
		},
		{
			name: "StillRunning",
			run: test.Run{
				ServiceTriggeredAt: time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, int(300*time.Millisecond), time.UTC),
			expected: 1200,
		},
		{
			name: "ZeroedDate",
			run: test.Run{
				ServiceTriggeredAt:        time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				ServiceTriggerCompletedAt: time.Unix(0, 0),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, int(300*time.Millisecond), time.UTC),
			expected: 1200,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			now := test.Now
			if c.now.Unix() > 0 {
				test.Now = func() time.Time {
					return c.now
				}
			}

			assert.Equal(t, c.expected, c.run.TriggerTime())
			test.Now = now
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
