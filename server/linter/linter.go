package linter

import (
	"context"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/plugins"
	"github.com/kubeshop/tracetest/server/linter/rules"
	"github.com/kubeshop/tracetest/server/traces"
)

var (
	commonPlugin = plugins.NewPlugin(
		analyzer.CommonID,
		rules.NewRegistry().
			Register(rules.NewEnforceDnsUsageRule()),
	)

	securityPlugin = plugins.NewPlugin(
		analyzer.SecurityID,
		rules.NewRegistry().
			Register(rules.NewEnsuresNoApiKeyLeakRule()).
			Register(rules.NewEnforceHttpsProtocolRule()),
	)

	standardsPlugin = plugins.NewPlugin(
		analyzer.StandardsID,
		rules.NewRegistry().
			Register(rules.NewEnsureSpanNamingRule()).
			Register(rules.NewRequiredAttributesRule()).
			Register(rules.NewEnsureAttributeNamingRule()).
			Register(rules.NewNotEmptyAttributesRule()),
	)

	DefaultPluginRegistry = plugins.NewRegistry().
				Register(standardsPlugin).
				Register(securityPlugin).
				Register(commonPlugin)
)

type Linter interface {
	Run(context.Context, traces.Trace, analyzer.Linter) (analyzer.LinterResult, error)
}

type linter struct {
	pluginsRegistry *plugins.Registry
}

func NewLinter(registry *plugins.Registry) Linter {
	return linter{
		pluginsRegistry: registry,
	}
}

var _ Linter = &linter{}

func (l linter) Run(ctx context.Context, trace traces.Trace, config analyzer.Linter) (analyzer.LinterResult, error) {
	cfgPlugins := config.EnabledPlugins()
	pluginResults := make([]analyzer.PluginResult, len(cfgPlugins))

	totalScore := 0
	passed := true
	for i, cfgPlugin := range cfgPlugins {
		plugin, err := l.pluginsRegistry.Get(cfgPlugin.ID)
		if err != nil {
			return analyzer.LinterResult{}, err
		}

		pluginResult, err := plugin.Execute(ctx, trace, cfgPlugin)
		if err != nil {
			return analyzer.LinterResult{}, err
		}

		passed = passed && pluginResult.Passed
		totalScore += pluginResult.Score
		pluginResults[i] = pluginResult
	}

	return analyzer.NewLinterResult(pluginResults, totalScore, passed), nil
}
