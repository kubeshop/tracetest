package analyzer_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzerResults(t *testing.T) {
	ruleResult1 := analyzer.RuleResult{
		ID:     "rule1",
		Weight: 25,
		Level:  "error",
		Passed: true,
	}

	ruleResult2 := analyzer.RuleResult{
		ID:     "rule2",
		Weight: 30,
		Level:  "warning",
		Passed: false,
	}

	ruleResult3 := analyzer.RuleResult{
		ID:     "rule3",
		Weight: 30,
		Level:  "disabled",
		Passed: true,
	}

	ruleResult4 := analyzer.RuleResult{
		ID:     "rule4",
		Weight: 30,
		Level:  "error",
		Passed: false,
	}

	ruleResult5 := analyzer.RuleResult{
		ID:     "rule5",
		Weight: 80,
		Level:  "error",
		Passed: true,
	}

	t.Run("plugin result calculates the score and overall summary based on rules", func(t *testing.T) {
		pluginResult := analyzer.PluginResult{
			ID: "plugin",

			Rules: []analyzer.RuleResult{
				ruleResult1,
				ruleResult2,
				ruleResult3,
				ruleResult4,
				ruleResult5,
			},
		}.CalculateResults()

		// passed score * 100% / total score
		assert.Equal(t, (105*100)/135, pluginResult.Score)
		assert.Equal(t, false, pluginResult.Passed)
	})

	t.Run("plugin result handles calculation for empty rule result list", func(t *testing.T) {
		pluginResult := analyzer.PluginResult{
			ID: "plugin",

			Rules: []analyzer.RuleResult{},
		}.CalculateResults()

		assert.Equal(t, 0, pluginResult.Score)
		assert.Equal(t, true, pluginResult.Passed)
	})

	t.Run("plugin result handles calculation for only warning level rules", func(t *testing.T) {
		pluginResult := analyzer.PluginResult{
			ID: "plugin",

			Rules: []analyzer.RuleResult{
				ruleResult2,
				ruleResult2,
				ruleResult2,
			},
		}.CalculateResults()

		assert.Equal(t, 0, pluginResult.Score)
		assert.Equal(t, true, pluginResult.Passed)
	})

	t.Run("plugin result handles calculation for only disabled level rules", func(t *testing.T) {
		pluginResult := analyzer.PluginResult{
			ID: "plugin",

			Rules: []analyzer.RuleResult{
				ruleResult3,
				ruleResult3,
			},
		}.CalculateResults()

		assert.Equal(t, 0, pluginResult.Score)
		assert.Equal(t, true, pluginResult.Passed)
	})
}
