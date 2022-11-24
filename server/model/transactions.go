package model

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/id"
)

type (
	Transaction struct {
		ID          id.ID
		CreatedAt   time.Time
		Name        string
		Description string
		Version     int
		Steps       []Test
		Summary     Summary
	}

	TransactionRunState string

	TransactionRun struct {
		ID                 int
		TransactionID      id.ID
		TransactionVersion int

		// Timestamps
		CreatedAt   time.Time
		CompletedAt time.Time

		//
		Steps []Run

		// trigger params
		State       TransactionRunState
		CurrentTest int

		// result info
		LastError error

		Metadata RunMetadata

		// environment
		Environment Environment
	}
)

const (
	TransactionRunStateCreated   TransactionRunState = "CREATED"
	TransactionRunStateExecuting TransactionRunState = "EXECUTING"
	TransactionRunStateFailed    TransactionRunState = "FAILED"
	TransactionRunStateFinished  TransactionRunState = "FINISHED"
)

func (rs TransactionRunState) IsFinal() bool {
	return rs == TransactionRunStateFailed || rs == TransactionRunStateFinished
}

func (t Transaction) HasID() bool {
	return t.ID != ""
}

func (t Transaction) NewRun() TransactionRun {

	return TransactionRun{
		TransactionID:      t.ID,
		TransactionVersion: t.Version,
		CreatedAt:          time.Now(),
		State:              TransactionRunStateCreated,
		Steps:              make([]Run, 0, len(t.Steps)),
		CurrentTest:        0,
	}
}

func (tr TransactionRun) ResourceID() string {
	return fmt.Sprintf("transaction/%s/run/%d", tr.TransactionID, tr.ID)
}
