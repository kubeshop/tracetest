package testsuite

import (
	"encoding/json"
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

func (tr TestSuiteRun) RunMetadata(step int) test.RunMetadata {
	tr.Metadata["step"] = fmt.Sprintf("%d", step+1)
	tr.Metadata["testsuite_run_id"] = fmt.Sprintf("%d", tr.ID)
	tr.Metadata["testsuite_id"] = string(tr.TestSuiteID)
	tr.Metadata["testsuite_version"] = fmt.Sprintf("%d", tr.TestSuiteVersion)

	return tr.Metadata
}

func (r TestSuiteRun) MarshalJSON() ([]byte, error) {
	encoded, err := EncodeRun(r)
	if err != nil {
		return nil, err
	}

	return json.Marshal(encoded)
}

func (r *TestSuiteRun) UnmarshalJSON(data []byte) error {
	var encoded EncodedTestSuiteRun

	if err := json.Unmarshal(data, &encoded); err != nil {
		return err
	}

	decoded, err := encoded.ToTestSuiteRun()
	if err != nil {
		return err
	}

	*r = decoded

	return nil
}
