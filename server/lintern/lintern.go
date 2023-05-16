package lintern

import (
	"context"

	"github.com/kubeshop/tracetest/server/lintern/rules"
	"github.com/kubeshop/tracetest/server/model"
)

type Runner interface {
	Run(context.Context, model.Trace) Result
}

type Result struct {
	Passed      bool           `json:"passed"`
	Score       uint           `json:"score"`
	RuleResults []rules.Result `json:"rules"`
}
