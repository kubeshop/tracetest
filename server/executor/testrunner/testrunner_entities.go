package testrunner

import (
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"golang.org/x/exp/slices"
)

var (
	RequiredGateAnalyzerScore RequiredGate = "analyzer-score"
	RequiredGateAnalyzerRules RequiredGate = "analyzer-rules"
	RequiredGateTestSpecs     RequiredGate = "test-specs"
)

const (
	ResourceName       = "TestRunner"
	ResourceNamePlural = "TestRunners"
)

var Operations = []resourcemanager.Operation{
	resourcemanager.OperationGet,
	resourcemanager.OperationList,
	resourcemanager.OperationUpdate,
}

var DefaultTestRunner = TestRunner{
	ID:   id.ID("current"),
	Name: "default",
	RequiredGates: []RequiredGate{
		RequiredGateTestSpecs,
		RequiredGateAnalyzerRules,
		RequiredGateAnalyzerScore,
	},
}

type (
	RequiredGate  string
	RequiredGates []RequiredGate

	TestRunner struct {
		ID            id.ID         `json:"id"`
		Name          string        `json:"name"`
		RequiredGates RequiredGates `json:"requiredGates"`
	}

	RequiredGatesResult struct {
		Required RequiredGates `json:"required"`
		Failed   RequiredGates `json:"failed"`
		Passed   bool          `json:"passed"`
	}
)

func NewRequiredGatesResult(required []RequiredGate) RequiredGatesResult {
	return RequiredGatesResult{
		Required: required,
		Passed:   true,
	}
}

func (r RequiredGatesResult) OnFailed(failed RequiredGate) RequiredGatesResult {
	if r.isRequiredGate(failed) {
		r.Passed = false
	}

	index := slices.Index(r.Failed, failed)
	if index == -1 {
		r.Failed = append(r.Failed, failed)
	}

	return r
}

func (r RequiredGatesResult) isRequiredGate(gate RequiredGate) bool {
	return slices.Index(r.Required, gate) > -1
}

func (ppc TestRunner) Validate() error {
	return nil
}

func (pp TestRunner) HasID() bool {
	return pp.ID.String() != ""
}

func (pp TestRunner) GetID() id.ID {
	return pp.ID
}
