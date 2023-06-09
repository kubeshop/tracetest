package model

import (
	"context"
)

type PluginResult struct {
	BasePlugin
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

type Plugin interface {
	Execute(context.Context, Trace) (PluginResult, error)
	Name() string
}

type BasePlugin struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rules       []Rule `json:"rules"`
}

type Rule interface {
	Evaluate(context.Context, Trace) (RuleResult, error)
}

type LinterResult struct {
	Plugins      []PluginResult `json:"plugins"`
	Score        int            `json:"score"`
	MinimumScore int            `json:"minimumScore"`
	Passed       bool           `json:"passed"`
}

type BaseRule struct {
	Name        string
	Description string
	Tips        []string
	Weight      int
}

type RuleResult struct {
	BaseRule
	Passed  bool     `json:"passed"`
	Results []Result `json:"results"`
}

type Result struct {
	SpanID string   `json:"span_id"`
	Passed bool     `json:"passed"`
	Errors []string `json:"error"`
}
