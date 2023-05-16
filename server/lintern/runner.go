package lintern

import (
	"context"

	"github.com/kubeshop/tracetest/server/model"
)

func NewRunner(ruleset ...Rule) Runner {
	return &runner{
		ruleset: ruleset,
	}
}

type runner struct {
	ruleset []Rule
}

func (r *runner) Run(ctx context.Context, trace model.Trace) Result {
	result := Result{
		Passed:      true,
		Score:       0,
		RuleResults: make([]RuleResult, 0, len(r.ruleset)),
	}

	totalWeight := r.getTotalWeight()

	for _, rule := range r.ruleset {
		ruleResult := rule.Run(ctx, trace)
		ruleResult.NormalizedWeight = float64(rule.GetWeight()) / float64(totalWeight)
		result.RuleResults = append(result.RuleResults, ruleResult)
	}

	return result
}

func (r *runner) getTotalWeight() uint {
	var totalWeight uint = 0
	for _, rule := range r.ruleset {
		totalWeight += rule.GetWeight()
	}

	return totalWeight
}
