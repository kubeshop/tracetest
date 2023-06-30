package rules

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/linter/metadata"
	"github.com/kubeshop/tracetest/server/linter/results"
	"github.com/kubeshop/tracetest/server/model"
)

type BaseRule struct {
	metadata metadata.RuleMetadata
}

var _ Rule = &BaseRule{}

func NewRule(metadata metadata.RuleMetadata) BaseRule {
	return BaseRule{
		metadata: metadata,
	}
}

func (r BaseRule) Slug() string {
	return r.metadata.Slug
}

func (r BaseRule) Evaluate(ctx context.Context, trace model.Trace, config analyzer.LinterRule) (results.RuleResult, error) {
	return results.RuleResult{}, fmt.Errorf("Rule Evaluation for %s is not implemented", r.metadata.Slug)
}
