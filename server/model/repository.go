package model

import (
	"context"

	"github.com/kubeshop/tracetest/server/id"
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
}

type EnvironmentRepository interface {
	CreateEnvironment(context.Context, Environment) (Environment, error)
	UpdateEnvironment(context.Context, Environment) (Environment, error)
	DeleteEnvironment(context.Context, Environment) error
	GetEnvironment(_ context.Context, id string) (Environment, error)
	EnvironmentIDExists(context.Context, string) (bool, error)
	GetEnvironments(_ context.Context, take, skip int32, query, sortBy, sortDirection string) (List[Environment], error)
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
	GetTransactionRun(context.Context, string, int) (TransactionRun, error)
	GetTransactionsRuns(context.Context, string, int32, int32) ([]TransactionRun, error)
}

type Repository interface {
	TestRepository
	RunRepository
	EnvironmentRepository

	TransactionRepository
	TransactionRunRepository

	ServerID() (id string, isNew bool, _ error)
	Drop() error
}
