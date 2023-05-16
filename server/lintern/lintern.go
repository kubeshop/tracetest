package lintern

import (
	"context"

	"github.com/kubeshop/tracetest/server/model"
)

type Runner interface {
	Run(context.Context, model.Trace) Result
}

type Rule interface {
	Run(context.Context, model.Trace) RuleResult
	GetWeight() uint
}

type Result struct {
	Passed      bool         `json:"passed"`
	Score       uint         `json:"score"`
	RuleResults []RuleResult `json:"rules"`
}

type RuleResult struct {
	Name             string       `json:"name"`
	Description      string       `json:"description"`
	Score            uint         `json:"score"`
	Passed           bool         `json:"passed"`
	NormalizedWeight float64      `json:"normalizedWeight"`
	SpansResults     []SpanResult `json:"spanResults"`
}

type SpanResult struct {
	SpanID string   `json:"spanId"`
	Passed bool     `json:"passed"`
	Score  uint     `json:"score"`
	Errors []string `json:"errors,omitempty"`
	Tipss  []string `json:"tips,omitempty"`
}
