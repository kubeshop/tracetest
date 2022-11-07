package model

import (
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

		// trigger params
		State       TransactionRunState
		Steps       []Test
		StepRuns    []Run
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

func NewTransactionRun(transaction Transaction) TransactionRun {
	return TransactionRun{
		TransactionID:      transaction.ID,
		TransactionVersion: transaction.Version,
		CreatedAt:          time.Now(),
		State:              TransactionRunStateCreated,
		Steps:              transaction.Steps,
		StepRuns:           make([]Run, 0, len(transaction.Steps)),
		CurrentTest:        0,
	}
}
