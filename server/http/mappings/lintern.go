package mappings

import (
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

func (m OpenAPI) LinternResult(in model.LinternResult) openapi.LinternResult {
	plugins := make([]openapi.LinternResultPlugin, len(in.Plugins))
	for i, plugin := range in.Plugins {
		plugins[i] = m.LinternResultPlugin(plugin)
	}

	return openapi.LinternResult{
		Passed:  in.Passed,
		Score:   int32(in.Score),
		Plugins: plugins,
	}
}

func (m OpenAPI) LinternResultPlugin(in model.PluginResult) openapi.LinternResultPlugin {
	rules := make([]openapi.LinternResultPluginRule, len(in.Rules))
	for i, rule := range in.Rules {
		rules[i] = m.LinternResultPluginRule(rule)
	}

	return openapi.LinternResultPlugin{
		Passed:      in.Passed,
		Score:       int32(in.Score),
		Description: in.Description,
		Name:        in.Name,
		Rules:       rules,
	}
}

func (m OpenAPI) LinternResultPluginRule(in model.RuleResult) openapi.LinternResultPluginRule {
	results := make([]openapi.LinternResultPluginRuleResult, len(in.Results))
	for i, result := range in.Results {
		results[i] = m.LinternResultPluginRuleResult(result)
	}

	return openapi.LinternResultPluginRule{
		Passed:      in.Passed,
		Description: in.Description,
		Name:        in.Name,
		Weight:      int32(in.Weight),
		Tips:        in.Tips,
		Results:     results,
	}
}

func (m OpenAPI) LinternResultPluginRuleResult(in model.Result) openapi.LinternResultPluginRuleResult {
	return openapi.LinternResultPluginRuleResult{
		SpanId:   in.SpanID,
		Passed:   in.Passed,
		Severity: "",
		Errors:   in.Errors,
	}
}
