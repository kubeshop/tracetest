package testdb

import (
	"context"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace"
)

var _ model.Repository = &MockRepository{}

type MockRepository struct {
	mock.Mock
	T mock.TestingT
}

// Close implements model.Repository
func (m *MockRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

// CreateTransactionRun implements model.Repository
func (*MockRepository) CreateTransactionRun(context.Context, model.TransactionRun) (model.TransactionRun, error) {
	panic("unimplemented")
}

func (m *MockRepository) ServerID() (string, bool, error) {
	args := m.Called()
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *MockRepository) TestIDExists(_ context.Context, id id.ID) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CreateTest(_ context.Context, test model.Test) (model.Test, error) {
	args := m.Called(test)
	return args.Get(0).(model.Test), args.Error(1)
}

func (m *MockRepository) UpdateTest(_ context.Context, test model.Test) (model.Test, error) {
	args := m.Called(test)
	return args.Get(0).(model.Test), args.Error(1)
}

func (m *MockRepository) UpdateTestVersion(_ context.Context, test model.Test) error {
	args := m.Called(test)
	return args.Error(0)
}

func (m *MockRepository) DeleteTest(_ context.Context, test model.Test) error {
	args := m.Called(test)
	return args.Error(0)
}

func (m *MockRepository) GetTestVersion(_ context.Context, id id.ID, version int) (model.Test, error) {
	args := m.Called(id, version)
	return args.Get(0).(model.Test), args.Error(1)
}

func (m *MockRepository) GetLatestTestVersion(_ context.Context, id id.ID) (model.Test, error) {
	args := m.Called(id)
	return args.Get(0).(model.Test), args.Error(1)
}

func (m *MockRepository) GetTests(_ context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.Test], error) {
	args := m.Called(take, skip, query, sortBy, sortDirection)
	tests := args.Get(0).([]model.Test)
	list := model.List[model.Test]{
		Items:      tests,
		TotalCount: len(tests),
	}
	return list, args.Error(1)
}

func (m *MockRepository) GetSpec(_ context.Context, test model.Test) (model.OrderedMap[model.SpanQuery, []model.Assertion], error) {
	args := m.Called(test)
	return args.Get(0).(model.OrderedMap[model.SpanQuery, []model.Assertion]), args.Error(1)
}

func (m *MockRepository) SetSpec(_ context.Context, test model.Test, def model.OrderedMap[model.SpanQuery, []model.Assertion]) error {
	args := m.Called(test, def)
	return args.Error(0)
}

func (m *MockRepository) CreateRun(_ context.Context, test model.Test, run model.Run) (model.Run, error) {
	args := m.Called(test, run)
	return args.Get(0).(model.Run), args.Error(1)
}

func (m *MockRepository) UpdateRun(_ context.Context, run model.Run) error {
	args := m.Called(run)
	return args.Error(0)
}

func (m *MockRepository) DeleteRun(_ context.Context, run model.Run) error {
	args := m.Called(run)
	return args.Error(0)
}

func (m *MockRepository) GetRun(_ context.Context, testID id.ID, id int) (model.Run, error) {
	args := m.Called(testID, id)
	return args.Get(0).(model.Run), args.Error(1)
}

func (m *MockRepository) GetTestRuns(_ context.Context, test model.Test, take int32, skip int32) (model.List[model.Run], error) {
	args := m.Called(test, take, skip)
	runs := args.Get(0).([]model.Run)
	list := model.List[model.Run]{
		Items:      runs,
		TotalCount: len(runs),
	}
	return list, args.Error(1)
}

func (m *MockRepository) GetRunByTraceID(_ context.Context, tid trace.TraceID) (model.Run, error) {
	args := m.Called(tid)
	return args.Get(0).(model.Run), args.Error(1)
}

func (m *MockRepository) Drop() error {
	args := m.Called()
	return args.Error(0)
}

// environments

func (m *MockRepository) CreateEnvironment(_ context.Context, environment model.Environment) (model.Environment, error) {
	args := m.Called(environment)
	return args.Get(0).(model.Environment), args.Error(1)
}

func (m *MockRepository) UpdateEnvironment(_ context.Context, environment model.Environment) (model.Environment, error) {
	args := m.Called(environment)
	return args.Get(0).(model.Environment), args.Error(1)
}

func (m *MockRepository) DeleteEnvironment(_ context.Context, environment model.Environment) error {
	args := m.Called(environment)
	return args.Error(0)
}

func (m *MockRepository) EnvironmentIDExists(_ context.Context, id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetEnvironment(_ context.Context, id string) (model.Environment, error) {
	args := m.Called(id)
	return args.Get(0).(model.Environment), args.Error(1)
}

func (m *MockRepository) GetEnvironments(_ context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.Environment], error) {
	args := m.Called(take, skip, query, sortBy, sortDirection)
	environments := args.Get(0).([]model.Environment)

	list := model.List[model.Environment]{
		Items:      environments,
		TotalCount: len(environments),
	}
	return list, args.Error(1)
}

func (m *MockRepository) CreateTransaction(_ context.Context, transaction model.Transaction) (model.Transaction, error) {
	args := m.Called(transaction)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *MockRepository) UpdateTransaction(_ context.Context, transaction model.Transaction) (model.Transaction, error) {
	args := m.Called(transaction)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *MockRepository) DeleteTransaction(_ context.Context, transaction model.Transaction) error {
	args := m.Called(transaction)
	return args.Error(1)
}

func (m *MockRepository) GetLatestTransactionVersion(_ context.Context, id id.ID) (model.Transaction, error) {
	args := m.Called(id)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *MockRepository) TransactionIDExists(_ context.Context, id id.ID) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetTransactionVersion(_ context.Context, id id.ID, version int) (model.Transaction, error) {
	args := m.Called(id, version)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *MockRepository) GetTransactions(_ context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.Transaction], error) {
	args := m.Called(take, skip, query, sortBy, sortDirection)
	transactions := args.Get(0).([]model.Transaction)
	list := model.List[model.Transaction]{
		Items:      transactions,
		TotalCount: len(transactions),
	}
	return list, args.Error(1)
}

// DeleteTransactionRun implements model.Repository
func (m *MockRepository) DeleteTransactionRun(ctx context.Context, run model.TransactionRun) error {
	args := m.Called(ctx, run)
	return args.Error(0)
}

// GetTransactionRun implements model.Repository
func (m *MockRepository) GetTransactionRun(ctx context.Context, transactionID id.ID, runID int) (model.TransactionRun, error) {
	args := m.Called(ctx, transactionID, runID)
	return args.Get(0).(model.TransactionRun), args.Error(1)
}

// GetTransactionsRuns implements model.Repository
func (m *MockRepository) GetTransactionsRuns(ctx context.Context, transactionID id.ID, take, skip int32) ([]model.TransactionRun, error) {
	args := m.Called(ctx, transactionID, take, skip)
	return args.Get(0).([]model.TransactionRun), args.Error(1)
}

// UpdateTransactionRun implements model.Repository
func (m *MockRepository) UpdateTransactionRun(ctx context.Context, run model.TransactionRun) error {
	args := m.Called(ctx, run)
	return args.Error(0)
}

// data stores

func (m *MockRepository) CreateDataStore(_ context.Context, dataStore model.DataStore) (model.DataStore, error) {
	args := m.Called(dataStore)
	return args.Get(0).(model.DataStore), args.Error(1)
}

func (m *MockRepository) UpdateDataStore(_ context.Context, dataStore model.DataStore) (model.DataStore, error) {
	args := m.Called(dataStore)
	return args.Get(0).(model.DataStore), args.Error(1)
}

func (m *MockRepository) DeleteDataStore(_ context.Context, dataStore model.DataStore) error {
	args := m.Called(dataStore)
	return args.Error(0)
}

func (m *MockRepository) DataStoreIDExists(_ context.Context, id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetDataStore(_ context.Context, id string) (model.DataStore, error) {
	args := m.Called(id)
	return args.Get(0).(model.DataStore), args.Error(1)
}

func (m *MockRepository) DefaultDataStore(_ context.Context) (model.DataStore, error) {
	args := m.Called()
	return args.Get(0).(model.DataStore), args.Error(1)
}

func (m *MockRepository) GetDataStores(_ context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.DataStore], error) {
	args := m.Called(take, skip, query, sortBy, sortDirection)
	dataStores := args.Get(0).([]model.DataStore)

	list := model.List[model.DataStore]{
		Items:      dataStores,
		TotalCount: len(dataStores),
	}
	return list, args.Error(1)
}
