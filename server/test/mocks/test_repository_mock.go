package mocks

import (
	"context"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/stretchr/testify/mock"
)

type TestRepository struct {
	mock.Mock
	T mock.TestingT
}

func (m *TestRepository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]test.Test, error) {
	args := m.Called(ctx, take, skip, query, sortBy, sortDirection)
	return args.Get(0).([]test.Test), args.Error(1)
}

func (m *TestRepository) ListAugmented(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]test.Test, error) {
	args := m.Called(ctx, take, skip, query, sortBy, sortDirection)
	return args.Get(0).([]test.Test), args.Error(1)
}

func (m *TestRepository) Count(ctx context.Context, query string) (int, error) {
	args := m.Called(ctx, query)
	return args.Int(0), args.Error(1)
}

func (m *TestRepository) SortingFields() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *TestRepository) Provision(ctx context.Context, t test.Test) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *TestRepository) SetID(t test.Test, id id.ID) test.Test {
	args := m.Called(t, id)
	return args.Get(0).(test.Test)
}

func (m *TestRepository) Get(ctx context.Context, id id.ID) (test.Test, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(test.Test), args.Error(1)
}

func (m *TestRepository) GetAugmented(ctx context.Context, id id.ID) (test.Test, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(test.Test), args.Error(1)
}

func (m *TestRepository) Exists(ctx context.Context, id id.ID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *TestRepository) GetVersion(ctx context.Context, id id.ID, version int) (test.Test, error) {
	args := m.Called(ctx, id, version)
	return args.Get(0).(test.Test), args.Error(1)
}

func (m *TestRepository) Create(ctx context.Context, t test.Test) (test.Test, error) {
	args := m.Called(ctx, t)
	return args.Get(0).(test.Test), args.Error(1)
}

func (m *TestRepository) Update(ctx context.Context, t test.Test) (test.Test, error) {
	args := m.Called(ctx, t)
	return args.Get(0).(test.Test), args.Error(1)
}

func (m *TestRepository) Delete(ctx context.Context, id id.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
