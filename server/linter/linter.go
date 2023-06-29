package linter

import (
	"context"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/metadata"
	"github.com/kubeshop/tracetest/server/linter/plugins"
	"github.com/kubeshop/tracetest/server/linter/results"
	"github.com/kubeshop/tracetest/server/linter/rules"
	"github.com/kubeshop/tracetest/server/model"
)

var (
	CommonPlugin = plugins.NewPlugin(
		metadata.CommonPlugin,
		rules.NewRegistry().Register(rules.NewEnforceDnsUsageRule(metadata.EnforceDnsRule)),
	)
	StandardsPlugin = plugins.NewPlugin(
		metadata.StandardsPlugin,
		rules.NewRegistry().
			Register(rules.NewEnsureSpanNamingRule(metadata.EnsureSpanNamingRule)).
			Register(rules.NewRequiredAttributesRule(metadata.RequiredAttributesRule)).
			Register(rules.NewEnsureAttributeNamingRule(metadata.EnsureAttributeNamingRule)).
			Register(rules.NewNotEmptyAttributesRule(metadata.NotEmptyAttributesRule)),
	)
	SecurityPlugin = plugins.NewPlugin(
		metadata.SecurityPlugin,
		rules.NewRegistry().
			Register(rules.NewEnforceHttpsProtocolRule(metadata.EnforceHttpsProtocolRule)).
			Register(rules.NewEnsuresNoApiKeyLeakRule(metadata.EnsuresNoApiKeyLeakRule)),
	)
	DefaultPluginRegistry = plugins.NewRegistry().
				Register(StandardsPlugin).
				Register(SecurityPlugin).
				Register(CommonPlugin)
)

type Linter interface {
	Run(context.Context, model.Trace) (results.LinterResult, error)
}

type linter struct {
	pluginsRegistry *plugins.Registry
	linterResource  analyzer.Linter
}

func NewLinter(linterResource analyzer.Linter, registry *plugins.Registry) Linter {
	return linter{
		pluginsRegistry: registry,
		linterResource:  linterResource,
	}
}

var _ Linter = &linter{}

func (l linter) Run(ctx context.Context, trace model.Trace) (results.LinterResult, error) {
	cfgPlugins := l.linterResource.EnabledPlugins()
	pluginResults := make([]results.PluginResult, len(cfgPlugins))

	totalScore := 0
	passed := true
	for i, cfgPlugin := range cfgPlugins {
		plugin, err := l.pluginsRegistry.Get(cfgPlugin.Slug)
		if err != nil {
			return results.LinterResult{}, err
		}

		pluginResult, err := plugin.Execute(ctx, trace, cfgPlugin)
		if err != nil {
			return results.LinterResult{}, err
		}

		passed = passed && pluginResult.Passed
		totalScore += pluginResult.Score
		pluginResults[i] = pluginResult
	}

	return results.NewLinterResult(pluginResults, totalScore, passed), nil
}
