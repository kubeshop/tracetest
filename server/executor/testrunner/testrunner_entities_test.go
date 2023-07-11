package testrunner_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/stretchr/testify/assert"
)

func TestTestRunnerEntities(t *testing.T) {
	t.Run("should fail the gates result because on of the required gates failed", func(t *testing.T) {
		gates := []testrunner.RequiredGate{
			testrunner.RequiredGateTestSpecs,
			testrunner.RequiredGateAnalyzerScore,
		}

		result := testrunner.NewRequiredGatesResult(gates).
			OnFailed(testrunner.RequiredGateAnalyzerRules).
			OnFailed(testrunner.RequiredGateTestSpecs)

		assert.Equal(t, gates, result.Required)
		assert.Equal(t, []testrunner.RequiredGate{testrunner.RequiredGateAnalyzerRules, testrunner.RequiredGateTestSpecs}, result.Failed)
		assert.Equal(t, false, result.Passed)
	})

	t.Run("should pass the gates result even if a non required gate fails", func(t *testing.T) {
		gates := []testrunner.RequiredGate{
			testrunner.RequiredGateTestSpecs,
			testrunner.RequiredGateAnalyzerScore,
		}

		result := testrunner.NewRequiredGatesResult(gates).
			OnFailed(testrunner.RequiredGateAnalyzerRules)

		assert.Equal(t, gates, result.Required)
		assert.Equal(t, []testrunner.RequiredGate{testrunner.RequiredGateAnalyzerRules}, result.Failed)
		assert.Equal(t, true, result.Passed)
	})

	t.Run("should handle empty required gates", func(t *testing.T) {
		gates := []testrunner.RequiredGate{}

		result := testrunner.NewRequiredGatesResult(gates).
			OnFailed(testrunner.RequiredGateAnalyzerRules).
			OnFailed(testrunner.RequiredGateTestSpecs).
			OnFailed(testrunner.RequiredGateAnalyzerScore)

		assert.Equal(t, gates, result.Required)
		assert.Equal(t, true, result.Passed)
	})

	t.Run("should handle failing the same gate twice", func(t *testing.T) {
		gates := []testrunner.RequiredGate{
			testrunner.RequiredGateTestSpecs,
			testrunner.RequiredGateAnalyzerScore,
		}

		result := testrunner.NewRequiredGatesResult(gates).
			OnFailed(testrunner.RequiredGateAnalyzerScore).
			OnFailed(testrunner.RequiredGateAnalyzerScore)

		assert.Equal(t, gates, result.Required)
		assert.Equal(t, []testrunner.RequiredGate{testrunner.RequiredGateAnalyzerScore}, result.Failed)
		assert.Len(t, result.Failed, 1)
		assert.Equal(t, false, result.Passed)
	})
}
