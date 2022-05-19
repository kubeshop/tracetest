package model

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type TestRepository interface {
	CreateTest(context.Context, Test) (Test, error)
	UpdateTest(context.Context, Test) error
	DeleteTest(context.Context, Test) error
	GetTest(context.Context, uuid.UUID) (Test, error)
	GetTests(_ context.Context, take, skip int32) ([]Test, error)
}

type DefinitionRepository interface {
	GetDefiniton(context.Context, Test) (Definition, error)
	SetDefiniton(context.Context, Test, Definition) error
}

type RunRepository interface {
	CreateRun(context.Context, Test, Run) (Run, error)
	UpdateRun(context.Context, Run) error
	GetRun(context.Context, uuid.UUID) (Run, error)
	GetTestRuns(_ context.Context, _ Test, take, skip int32) ([]Run, error)
	GetRunByTraceID(context.Context, Test, trace.TraceID) (Run, error)
}

type Repository interface {
	TestRepository
	DefinitionRepository
	RunRepository

	Drop() error
}
