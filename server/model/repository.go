package model

import (
	"context"

	"github.com/kubeshop/tracetest/server/id"
	"go.opentelemetry.io/otel/trace"
)

type TestRepository interface {
	CreateTest(context.Context, Test) (Test, error)
	UpdateTest(context.Context, Test) (Test, error)
	DeleteTest(context.Context, Test) error
	IDExists(context.Context, id.ID) (bool, error)
	GetLatestTestVersion(context.Context, id.ID) (Test, error)
	GetTestVersion(_ context.Context, _ id.ID, verson int) (Test, error)
	GetTests(_ context.Context, take, skip int32, query, sortBy, sortDirection string) ([]Test, error)
}

type RunRepository interface {
	CreateRun(context.Context, Test, Run) (Run, error)
	UpdateRun(context.Context, Run) error
	DeleteRun(context.Context, Run) error
	GetRun(_ context.Context, testID id.ID, runID int) (Run, error)
	GetTestRuns(_ context.Context, _ Test, take, skip int32) ([]Run, error)
	GetRunByTraceID(context.Context, trace.TraceID) (Run, error)
}

type Repository interface {
	TestRepository
	RunRepository

	ServerID() (id string, isNew bool, _ error)
	Drop() error
}
