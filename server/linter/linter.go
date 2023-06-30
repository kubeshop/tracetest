package linter

import (
	"context"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/plugins"
	"github.com/kubeshop/tracetest/server/linter/rules"
	"github.com/kubeshop/tracetest/server/model"
)

var (
	commonPlugin = plugins.NewPlugin(
		analyzer.CommonId,
		rules.NewRegistry().
			Register(rules.NewEnforceDnsUsageRule()),
	)

	securityPlugin = plugins.NewPlugin(
		analyzer.SecurityId,
		rules.NewRegistry().
			Register(rules.NewEnsuresNoApiKeyLeakRule()).
			Register(rules.NewEnforceHttpsProtocolRule()),
	)

	standardsPlugin = plugins.NewPlugin(
		analyzer.StandardsId,
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
	Run(context.Context, model.Trace) (analyzer.LinterResult, error)
}

type linter struct {
	pluginsRegistry *plugins.Registry
	config          analyzer.Linter
}

func NewLinter(config analyzer.Linter, registry *plugins.Registry) Linter {
	return linter{
		pluginsRegistry: registry,
		config:          config,
	}
}

var _ Linter = &linter{}

func (l linter) Run(ctx context.Context, trace model.Trace) (analyzer.LinterResult, error) {
	cfgPlugins := l.config.EnabledPlugins()
	pluginResults := make([]analyzer.PluginResult, len(cfgPlugins))

	totalScore := 0
	passed := true
	for i, cfgPlugin := range cfgPlugins {
		plugin, err := l.pluginsRegistry.Get(cfgPlugin.Id)
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
