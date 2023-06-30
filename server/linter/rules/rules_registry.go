package rules

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
)

type Rule interface {
	Id() string
	Evaluate(context.Context, model.Trace, analyzer.LinterRule) (analyzer.RuleResult, error)
}

type RuleRegistry interface {
	List() []string
	Get(string) (Rule, error)
	Register(Rule) RuleRegistry
}

type ruleRegistry struct {
	rules map[string]Rule
}

func NewRegistry() RuleRegistry {
	return &ruleRegistry{
		rules: make(map[string]Rule),
	}
}

func (r *ruleRegistry) List() []string {
	keys := make([]string, 0, len(r.rules))
	for k := range r.rules {
		keys = append(keys, k)
	}

	return keys
}

func (r *ruleRegistry) Get(ruleName string) (Rule, error) {
	if rule, ok := r.rules[ruleName]; ok {
		return rule, nil
	}

	return nil, fmt.Errorf("plugin %s not found", ruleName)
}

func (r *ruleRegistry) Register(rule Rule) RuleRegistry {
	r.rules[rule.Id()] = rule
	return r
}
