package testsuite

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/variableset"
)

type TestSuiteRun struct {
	ID               int
	TestSuiteID      id.ID
	TestSuiteVersion int

	// Timestamps
	CreatedAt   time.Time
	CompletedAt time.Time

	// steps
	StepIDs []int
	Steps   []test.Run

	// trigger params
	State       TestSuiteRunState
	CurrentTest int

	// result info
	LastError                   error
	Pass                        int
	Fail                        int
	AllStepsRequiredGatesPassed bool

	Metadata test.RunMetadata

	// variable set
	VariableSet   variableset.VariableSet
	RequiredGates *[]testrunner.RequiredGate
}

func (tr TestSuiteRun) ResourceID() string {
	return fmt.Sprintf("testsuites/%s/run/%d", tr.TestSuiteID, tr.ID)
}

func (tr TestSuiteRun) ResultsCount() (pass, fail int) {
	if tr.Steps == nil {
		return
	}

	for _, step := range tr.Steps {
		stepPass, stepFail := step.ResultsCount()

		pass += stepPass
		fail += stepFail
	}

	return
}

func (tr TestSuiteRun) StepsGatesValidation() bool {
	for _, step := range tr.Steps {
		if !step.RequiredGatesResult.Passed {
			return false
		}
	}

	return true
}
