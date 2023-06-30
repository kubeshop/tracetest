package results

import (
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/metadata"
)

type Result struct {
	SpanID string  `json:"span_id"`
	Passed bool    `json:"passed"`
	Errors []Error `json:"errors"`
}

type Error struct {
	Value       string   `json:"value"`
	Expected    string   `json:"expected"`
	Description string   `json:"description"`
	Suggestions []string `json:"suggestions"`
}

type EvalRuleResult struct {
	Passed  bool     `json:"passed"`
	Results []Result `json:"results"`
}

type RuleResult struct {
	// metadata
	Slug             string   `json:"slug"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	ErrorDescription string   `json:"errorDescription"`
	Tips             []string `json:"tips"`

	// config
	Weight int    `json:"weight"`
	Level  string `json:"level"`

	// results
	Passed  bool     `json:"passed"`
	Results []Result `json:"results"`
}

func NewRuleResult(metadata metadata.RuleMetadata, config analyzer.LinterRule, results EvalRuleResult) RuleResult {
	return RuleResult{
		Slug:             metadata.Slug,
		Name:             metadata.Name,
		Description:      metadata.Description,
		ErrorDescription: metadata.ErrorDescription,
		Tips:             metadata.Tips,
		Weight:           config.Weight,
		Level:            config.ErrorLevel,
		Passed:           results.Passed,
		Results:          results.Results,
	}
}

type LinterResult struct {
	Plugins      []PluginResult `json:"plugins"`
	Score        int            `json:"score"`
	MinimumScore int            `json:"minimumScore"`
	Passed       bool           `json:"passed"`
}

func NewLinterResult(pluginResults []PluginResult, totalScore int, passed bool) LinterResult {
	return LinterResult{
		Plugins: pluginResults,
		Score:   totalScore / len(pluginResults),
		Passed:  passed,
	}
}

type PluginResult struct {
	// metadata
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`

	//results
	Passed bool         `json:"passed"`
	Score  int          `json:"score"`
	Rules  []RuleResult `json:"rules"`
}

func (pr PluginResult) CalculateResults() PluginResult {
	failedScore := 0
	passed := true

	for _, result := range pr.Rules {
		if !result.Passed {
			passed = false
			failedScore += result.Weight
		}
	}

	pr.Score = 100 - failedScore
	pr.Passed = passed
	return pr
}
