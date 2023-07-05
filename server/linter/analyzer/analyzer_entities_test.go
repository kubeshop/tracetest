package analyzer_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/stretchr/testify/assert"
)

func getLinter() analyzer.Linter {
	return analyzer.Linter{
		ID:           id.ID("current"),
		Name:         "analyzer",
		Enabled:      true,
		Plugins:      []analyzer.LinterPlugin{},
		MinimumScore: 0,
	}
}

func TestAnalyzerEntities(t *testing.T) {
	t.Run("validate returns errors for unknown plugin", func(t *testing.T) {
		linter := getLinter()

		plugin := analyzer.StandardsPlugin
		plugin.ID = "unknown"
		linter.Plugins = append(linter.Plugins, plugin)
		err := linter.Validate()

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "plugin unknown not supported")
	})

	t.Run("validate returns errors for unknown rule", func(t *testing.T) {
		linter := getLinter()

		plugin := analyzer.LinterPlugin{
			ID: analyzer.StandardsID,
			Rules: []analyzer.LinterRule{{
				ID: "unknown1",
			}, {
				ID: "unknown2",
			}, {
				ID: "unknown3",
			}, {
				ID: "unknown4",
			}},
		}
		linter.Plugins = append(linter.Plugins, plugin)
		err := linter.Validate()

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "rule unknown1 not found for plugin standards")
	})

	t.Run("validate returns errors because missing rules", func(t *testing.T) {
		linter := getLinter()

		plugin := analyzer.LinterPlugin{
			ID: analyzer.StandardsID,
			Rules: []analyzer.LinterRule{{
				ID: "unknown1",
			}, {
				ID: "unknown2",
			}, {
				ID: "unknown3",
			}},
		}
		linter.Plugins = append(linter.Plugins, plugin)
		err := linter.Validate()

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "plugin standards requires 4 rules, but 3 provided, supported rules for plugin are")
	})

	t.Run("returns the valid list of enabled plugins", func(t *testing.T) {
		linter := getLinter()

		standards := analyzer.StandardsPlugin
		standards.Enabled = false

		linter.Plugins = []analyzer.LinterPlugin{standards, analyzer.SecurityPlugin}
		enabledPlugins := linter.EnabledPlugins()

		assert.Equal(t, 1, len(enabledPlugins))
		assert.Equal(t, enabledPlugins, []analyzer.LinterPlugin{analyzer.SecurityPlugin})
	})

	t.Run("returns the list of plugins with metadata even if missing", func(t *testing.T) {
		linter := getLinter()
		linter.Plugins = []analyzer.LinterPlugin{{
			ID: analyzer.StandardsID,
			Rules: []analyzer.LinterRule{{
				ID:         analyzer.EnsureSpanNamingRuleID,
				Weight:     20,
				ErrorLevel: analyzer.ErrorLevelWarning,
			}, {
				ID:         analyzer.EnsureAttributeNamingRuleID,
				Weight:     20,
				ErrorLevel: analyzer.ErrorLevelWarning,
			}},
		}}
		linter, err := linter.WithMetadata()

		plugin := linter.Plugins[0]
		assert.Nil(t, err, nil)
		assert.Equal(t, 2, len(plugin.Rules))
		assert.Equal(t, analyzer.StandardsPlugin.Name, plugin.Name)
		assert.Equal(t, analyzer.StandardsPlugin.Description, plugin.Description)

		spanNameRule := plugin.Rules[0]
		assert.Equal(t, analyzer.EnsureSpanNamingRule.Name, spanNameRule.Name)
		assert.Equal(t, analyzer.EnsureSpanNamingRule.Description, spanNameRule.Description)
		assert.Equal(t, analyzer.EnsureSpanNamingRule.ErrorDescription, spanNameRule.ErrorDescription)
		assert.Equal(t, analyzer.EnsureSpanNamingRule.Tips, spanNameRule.Tips)

		attributeNameRule := plugin.Rules[1]
		assert.Equal(t, analyzer.EnsureAttributeNamingRule.Name, attributeNameRule.Name)
		assert.Equal(t, analyzer.EnsureAttributeNamingRule.Description, attributeNameRule.Description)
		assert.Equal(t, analyzer.EnsureAttributeNamingRule.ErrorDescription, attributeNameRule.ErrorDescription)
		assert.Equal(t, analyzer.EnsureAttributeNamingRule.Tips, attributeNameRule.Tips)
	})
}
