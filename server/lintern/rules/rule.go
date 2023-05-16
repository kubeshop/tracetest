package rules

import (
	"context"

	"github.com/kubeshop/tracetest/server/model"
)

type Rule interface {
	Run(context.Context, model.Trace) Result
	GetWeight() uint
}

type Result struct {
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
	Tips   []string `json:"tips,omitempty"`
}

type baseRule struct {
	name        string
	description string
	weight      uint
}

func (r *baseRule) GetWeight() uint {
	return r.weight
}
