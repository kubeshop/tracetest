package transaction

import (
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
)

const (
	TransactionResourceName       = "Transaction"
	TransactionResourceNamePlural = "Transactions"
)

type Transaction struct {
	ID          id.ID          `json:"id"`
	CreatedAt   *time.Time     `json:"createdAt,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Version     *int           `json:"version,omitempty"`
	StepIDs     []id.ID        `json:"steps"`
	Steps       []model.Test   `json:"fullSteps,omitempty"`
	Summary     *model.Summary `json:"summary,omitempty"`
}

func setVersion(t *Transaction, v int) {
	t.Version = &v
}

func (t Transaction) GetVersion() int {
	if t.Version == nil {
		return 0
	}
	return *t.Version
}

func setCreatedAt(t *Transaction, d time.Time) {
	t.CreatedAt = &d
}

func (t Transaction) GetCreatedAt() time.Time {
	if t.CreatedAt == nil {
		return time.Time{}
	}
	return *t.CreatedAt
}

func (t Transaction) HasID() bool {
	return t.ID != ""
}

func (t Transaction) GetID() id.ID {
	return t.ID
}

func (t Transaction) Validate() error {
	return nil
}

func (t Transaction) NewRun() TransactionRun {

	return TransactionRun{
		TransactionID:      t.ID,
		TransactionVersion: t.GetVersion(),
		CreatedAt:          time.Now().UTC(),
		State:              TransactionRunStateCreated,
		Steps:              make([]model.Run, 0, len(t.StepIDs)),
		CurrentTest:        0,
	}
}

type TransactionRunState string

const (
	TransactionRunStateCreated   TransactionRunState = "CREATED"
	TransactionRunStateExecuting TransactionRunState = "EXECUTING"
	TransactionRunStateFailed    TransactionRunState = "FAILED"
	TransactionRunStateFinished  TransactionRunState = "FINISHED"
)

func (rs TransactionRunState) IsFinal() bool {
	return rs == TransactionRunStateFailed || rs == TransactionRunStateFinished
}
