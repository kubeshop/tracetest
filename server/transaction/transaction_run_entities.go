package transaction

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
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
	Steps   []model.Run

	// trigger params
	State       TransactionRunState
	CurrentTest int

	// result info
	LastError error
	Pass      int
	Fail      int

	Metadata model.RunMetadata

	// environment
	Environment environment.Environment
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
