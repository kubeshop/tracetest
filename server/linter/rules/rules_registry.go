package rules

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
)

type Rule interface {
	ID() string
	Evaluate(context.Context, model.Trace, analyzer.LinterRule) (analyzer.RuleResult, error)
}

type RuleRegistry struct {
	rules map[string]Rule
}

func NewRegistry() *RuleRegistry {
	return &RuleRegistry{
		rules: make(map[string]Rule),
	}
}

func (r *RuleRegistry) List() []string {
	keys := make([]string, 0, len(r.rules))
	for k := range r.rules {
		keys = append(keys, k)
	}

	return keys
}

func (r *RuleRegistry) Get(ruleName string) (Rule, error) {
	if rule, ok := r.rules[ruleName]; ok {
		return rule, nil
	}

	return nil, fmt.Errorf("rule %s not found", ruleName)
}

func (r *RuleRegistry) Register(rule Rule) *RuleRegistry {
	r.rules[rule.ID()] = rule
	return r
}
