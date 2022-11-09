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

	transactionStep struct {
		ID   id.ID
		Name string
	}

	TransactionStepRun struct {
		ID          int
		TestID      id.ID
		State       RunState
		Environment Environment
		Outputs     OrderedMap[string, string]
	}

	TransactionRun struct {
		ID                 int
		TransactionID      id.ID
		TransactionVersion int

		// Timestamps
		CreatedAt   time.Time
		CompletedAt time.Time

		// trigger params
		State       TransactionRunState
		Steps       []transactionStep
		StepRuns    []TransactionStepRun
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
	testIds := make([]transactionStep, 0, len(transaction.Steps))

	for _, test := range transaction.Steps {
		testIds = append(testIds, transactionStep{ID: test.ID, Name: test.Name})
	}

	return TransactionRun{
		TransactionID:      transaction.ID,
		TransactionVersion: transaction.Version,
		CreatedAt:          time.Now(),
		State:              TransactionRunStateCreated,
		Steps:              testIds,
		StepRuns:           make([]TransactionStepRun, 0, len(transaction.Steps)),
		CurrentTest:        0,
	}
}
