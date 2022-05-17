package testdb

import (
	"context"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/model"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace"
)

var _ model.Repository = &MockRepository{}

type MockRepository struct {
	mock.Mock
	T mock.TestingT
}

func (m MockRepository) CreateTest(_ context.Context, test model.Test) (model.Test, error) {
	args := m.Called(test)
	return args.Get(0).(model.Test), args.Error(1)
}

func (m MockRepository) UpdateTest(_ context.Context, test model.Test) error {
	args := m.Called(test)
	return args.Error(0)
}

func (m MockRepository) DeleteTest(_ context.Context, test model.Test) error {
	args := m.Called(test)
	return args.Error(0)
}

func (m MockRepository) GetTest(_ context.Context, id uuid.UUID) (model.Test, error) {
	args := m.Called(id)
	return args.Get(0).(model.Test), args.Error(1)
}

func (m MockRepository) GetTests(_ context.Context, take int32, skip int32) ([]model.Test, error) {
	args := m.Called(take, skip)
	return args.Get(0).([]model.Test), args.Error(1)
}

func (m MockRepository) GetDefiniton(_ context.Context, test model.Test) (model.Definition, error) {
	args := m.Called(test)
	return args.Get(0).(model.Definition), args.Error(1)
}

func (m MockRepository) SetDefiniton(_ context.Context, test model.Test, def model.Definition) error {
	args := m.Called(test, def)
	return args.Error(0)
}

func (m MockRepository) CreateRun(_ context.Context, test model.Test, run model.Run) (model.Run, error) {
	args := m.Called(test, run)
	return args.Get(0).(model.Run), args.Error(1)
}

func (m MockRepository) UpdateRun(_ context.Context, run model.Run) error {
	args := m.Called(run)
	return args.Error(0)
}

func (m MockRepository) GetRun(_ context.Context, id uuid.UUID) (model.Run, error) {
	args := m.Called(id)
	return args.Get(0).(model.Run), args.Error(1)
}

func (m MockRepository) GetTestRuns(_ context.Context, test model.Test, take int32, skip int32) ([]model.Run, error) {
	args := m.Called(test, take, skip)
	return args.Get(0).([]model.Run), args.Error(1)
}

func (m MockRepository) GetRunByTraceID(_ context.Context, test model.Test, tid trace.TraceID) (model.Run, error) {
	args := m.Called(test, tid)
	return args.Get(0).(model.Run), args.Error(1)
}

func (m MockRepository) Drop() error {
	args := m.Called()
	return args.Error(0)
}
