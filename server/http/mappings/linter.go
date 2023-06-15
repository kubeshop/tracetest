package mappings

import (
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

func (m OpenAPI) LinterResult(in model.LinterResult) openapi.LinterResult {
	plugins := make([]openapi.LinterResultPlugin, len(in.Plugins))
	for i, plugin := range in.Plugins {
		plugins[i] = m.LinterResultPlugin(plugin)
	}

	return openapi.LinterResult{
		Passed:       in.Passed,
		Score:        int32(in.Score),
		Plugins:      plugins,
		MinimumScore: int32(in.MinimumScore),
	}
}

func (m OpenAPI) LinterResultPlugin(in model.PluginResult) openapi.LinterResultPlugin {
	rules := make([]openapi.LinterResultPluginRule, len(in.Rules))
	for i, rule := range in.Rules {
		rules[i] = m.LinterResultPluginRule(rule)
	}

	return openapi.LinterResultPlugin{
		Passed:      in.Passed,
		Score:       int32(in.Score),
		Description: in.Description,
		Name:        in.Name,
		Rules:       rules,
	}
}

func (m OpenAPI) LinterResultPluginRule(in model.RuleResult) openapi.LinterResultPluginRule {
	results := make([]openapi.LinterResultPluginRuleResult, len(in.Results))
	for i, result := range in.Results {
		results[i] = m.LinterResultPluginRuleResult(result)
	}

	return openapi.LinterResultPluginRule{
		Passed:      in.Passed,
		Description: in.Description,
		Name:        in.Name,
		Weight:      int32(in.Weight),
		Tips:        in.Tips,
		Results:     results,
	}
}

func (m OpenAPI) LinterResultPluginRuleResult(in model.Result) openapi.LinterResultPluginRuleResult {
	groupedErrors := make([]openapi.LinterResultPluginRuleResultGroupedError, len(in.GroupedErrors))
	for i, groupedError := range in.GroupedErrors {
		groupedErrors[i] = m.LinterResultPluginRuleResultGroupedError(groupedError)
	}

	return openapi.LinterResultPluginRuleResult{
		SpanId:        in.SpanID,
		Passed:        in.Passed,
		Severity:      "",
		Errors:        in.Errors,
		GroupedErrors: groupedErrors,
	}
}

func (m OpenAPI) LinterResultPluginRuleResultGroupedError(in model.GroupedError) openapi.LinterResultPluginRuleResultGroupedError {
	return openapi.LinterResultPluginRuleResultGroupedError{
		Error:  in.Error,
		Values: in.Values,
	}
}
