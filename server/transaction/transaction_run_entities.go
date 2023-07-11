package transaction

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
)

type TransactionRun struct {
	ID                 int
	TransactionID      id.ID
	TransactionVersion int

	// Timestamps
	CreatedAt   time.Time
	CompletedAt time.Time

	// steps
	StepIDs []int
	Steps   []test.Run

	// trigger params
	State       TransactionRunState
	CurrentTest int

	// result info
	LastError error
	Pass      int
	Fail      int

	Metadata test.RunMetadata

	// environment
	Environment   environment.Environment
	RequiredGates *[]testrunner.RequiredGate
}

func (tr TransactionRun) ResourceID() string {
	return fmt.Sprintf("transaction/%s/run/%d", tr.TransactionID, tr.ID)
}

func (tr TransactionRun) ResultsCount() (pass, fail int) {
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
