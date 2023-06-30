package rules

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
)

type BaseRule struct {
	id string
}

var _ Rule = &BaseRule{}

func NewRule(id string) BaseRule {
	return BaseRule{id}
}

func (r BaseRule) Id() string {
	return r.id
}

func (r BaseRule) Evaluate(ctx context.Context, trace model.Trace, config analyzer.LinterRule) (analyzer.RuleResult, error) {
	return analyzer.RuleResult{}, fmt.Errorf("Rule Evaluation for %s is not implemented", r.id)
}
