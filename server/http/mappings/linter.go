package mappings

import (
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/openapi"
)

func (m OpenAPI) LinterResult(in analyzer.LinterResult) openapi.LinterResult {
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

func (m OpenAPI) LinterResultPlugin(in analyzer.PluginResult) openapi.LinterResultPlugin {
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

func (m OpenAPI) LinterResultPluginRule(in analyzer.RuleResult) openapi.LinterResultPluginRule {
	results := make([]openapi.LinterResultPluginRuleResult, len(in.Results))
	for i, result := range in.Results {
		results[i] = m.LinterResultPluginRuleResult(result)
	}

	return openapi.LinterResultPluginRule{
		Id:               in.ID,
		Passed:           in.Passed,
		Description:      in.Description,
		ErrorDescription: in.ErrorDescription,
		Name:             in.Name,
		Weight:           int32(in.Weight),
		Tips:             in.Tips,
		Documentation:    in.Documentation,
		Results:          results,
		Level:            in.Level,
	}
}

func (m OpenAPI) LinterResultPluginRuleResult(in analyzer.Result) openapi.LinterResultPluginRuleResult {
	errors := make([]openapi.LinterResultPluginRuleResultError, len(in.Errors))
	for i, error := range in.Errors {
		errors[i] = m.LinterResultPluginRuleResultError(error)
	}

	return openapi.LinterResultPluginRuleResult{
		SpanId:   in.SpanID,
		Passed:   in.Passed,
		Severity: "",
		Errors:   errors,
	}
}

func (m OpenAPI) LinterResultPluginRuleResultError(in analyzer.Error) openapi.LinterResultPluginRuleResultError {
	return openapi.LinterResultPluginRuleResultError{
		Value:       in.Value,
		Expected:    in.Expected,
		Description: in.Description,
		Suggestions: in.Suggestions,
	}
}
