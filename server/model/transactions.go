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

	TransactionStep struct {
		ID      id.ID
		Name    string
		Trigger Trigger
	}

	TransactionStepRun struct {
		ID          int
		TestID      id.ID
		State       RunState
		Result      RunResults
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
		Steps       []TransactionStep
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
	tests := make([]TransactionStep, 0, len(transaction.Steps))

	for _, test := range transaction.Steps {
		tests = append(tests, TransactionStep{ID: test.ID, Name: test.Name, Trigger: test.ServiceUnderTest})
	}

	return TransactionRun{
		TransactionID:      transaction.ID,
		TransactionVersion: transaction.Version,
		CreatedAt:          time.Now(),
		State:              TransactionRunStateCreated,
		Steps:              tests,
		StepRuns:           make([]TransactionStepRun, 0, len(transaction.Steps)),
		CurrentTest:        0,
	}
}

func (run TransactionRun) InjectOutputsIntoEnvironment(env Environment) Environment {
	if run.CurrentTest == 0 {
		return env
	}

	lastExecutedTest := run.StepRuns[run.CurrentTest-1]
	lastEnvironment := lastExecutedTest.Environment
	newEnvVariables := make([]EnvironmentValue, 0)
	lastExecutedTest.Outputs.ForEach(func(key, val string) error {
		newEnvVariables = append(newEnvVariables, EnvironmentValue{
			Key:   key,
			Value: val,
		})

		return nil
	})

	newEnvironment := Environment{Values: newEnvVariables}

	return lastEnvironment.Merge(newEnvironment)
}

func (r TransactionRun) ResourceID() string {
	return fmt.Sprintf("transaction/%s/run/%d", r.TransactionID, r.ID)
}
