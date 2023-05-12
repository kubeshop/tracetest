package model

import (
	"context"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"go.opentelemetry.io/otel/trace"
)

type List[T any] struct {
	Items      []T
	TotalCount int
}

type TestRepository interface {
	CreateTest(context.Context, Test) (Test, error)
	UpdateTest(context.Context, Test) (Test, error)
	DeleteTest(context.Context, Test) error
	TestIDExists(context.Context, id.ID) (bool, error)
	GetLatestTestVersion(context.Context, id.ID) (Test, error)
	GetTestVersion(_ context.Context, _ id.ID, version int) (Test, error)
	GetTests(_ context.Context, take, skip int32, query, sortBy, sortDirection string) (List[Test], error)
}

type RunRepository interface {
	CreateRun(context.Context, Test, Run) (Run, error)
	UpdateRun(context.Context, Run) error
	DeleteRun(context.Context, Run) error
	GetRun(_ context.Context, testID id.ID, runID int) (Run, error)
	GetTestRuns(_ context.Context, _ Test, take, skip int32) (List[Run], error)
	GetRunByTraceID(context.Context, trace.TraceID) (Run, error)
	GetLatestRunByTestVersion(context.Context, id.ID, int) (Run, error)
}

type TransactionRepository interface {
	CreateTransaction(context.Context, Transaction) (Transaction, error)
	UpdateTransaction(context.Context, Transaction) (Transaction, error)
	DeleteTransaction(context.Context, Transaction) error
	TransactionIDExists(context.Context, id.ID) (bool, error)
	GetLatestTransactionVersion(context.Context, id.ID) (Transaction, error)
	GetTransactionVersion(_ context.Context, _ id.ID, version int) (Transaction, error)
	GetTransactions(_ context.Context, take, skip int32, query, sortBy, sortDirection string) (List[Transaction], error)
}

type TransactionRunRepository interface {
	CreateTransactionRun(context.Context, TransactionRun) (TransactionRun, error)
	UpdateTransactionRun(context.Context, TransactionRun) error
	DeleteTransactionRun(context.Context, TransactionRun) error
	GetTransactionRun(context.Context, id.ID, int) (TransactionRun, error)
	GetTransactionsRuns(context.Context, id.ID, int32, int32) ([]TransactionRun, error)
	GetLatestRunByTransactionVersion(context.Context, id.ID, int) (TransactionRun, error)
}

type TestRunEventRepository interface {
	CreateTestRunEvent(context.Context, TestRunEvent) error
	GetTestRunEvents(context.Context, id.ID, int) ([]TestRunEvent, error)
}

type Repository interface {
	TestRepository
	RunRepository

	TransactionRepository
	TransactionRunRepository

	TestRunEventRepository

	ServerID() (id string, isNew bool, _ error)
	Close() error
	Drop() error
}
