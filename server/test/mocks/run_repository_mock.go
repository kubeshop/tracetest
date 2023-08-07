package mocks

import (
	"context"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace"
)

type RunRepository struct {
	mock.Mock
	T mock.TestingT
}

func (m *RunRepository) CreateRun(ctx context.Context, t test.Test, r test.Run) (test.Run, error) {
	args := m.Called(ctx, t, r)
	return args.Get(0).(test.Run), args.Error(1)
}

func (m *RunRepository) UpdateRun(ctx context.Context, r test.Run) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *RunRepository) DeleteRun(ctx context.Context, r test.Run) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *RunRepository) GetRun(ctx context.Context, testID id.ID, runID int) (test.Run, error) {
	args := m.Called(ctx, testID, runID)
	return args.Get(0).(test.Run), args.Error(1)
}

func (m *RunRepository) GetTestRuns(ctx context.Context, t test.Test, take, skip int32) ([]test.Run, error) {
	args := m.Called(ctx, t, take, skip)
	return args.Get(0).([]test.Run), args.Error(1)
}

func (m *RunRepository) GetRunByTraceID(ctx context.Context, traceID trace.TraceID) (test.Run, error) {
	args := m.Called(ctx, traceID)
	return args.Get(0).(test.Run), args.Error(1)
}

func (m *RunRepository) GetLatestRunByTestVersion(ctx context.Context, id id.ID, version int) (test.Run, error) {
	args := m.Called(ctx, id, version)
	return args.Get(0).(test.Run), args.Error(1)
}

func (m *RunRepository) GetTestSuiteRunSteps(ctx context.Context, id id.ID, runID int) ([]test.Run, error) {
	args := m.Called(ctx, id, runID)
	return args.Get(0).([]test.Run), args.Error(1)
}
