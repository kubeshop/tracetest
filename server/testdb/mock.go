package testdb

import (
	"context"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/transaction"
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

func (m *MockRepository) GetSpec(_ context.Context, test model.Test) (maps.Ordered[model.SpanQuery, []model.Assertion], error) {
	args := m.Called(test)
	return args.Get(0).(maps.Ordered[model.SpanQuery, []model.Assertion]), args.Error(1)
}

func (m *MockRepository) SetSpec(_ context.Context, test model.Test, def maps.Ordered[model.SpanQuery, []model.Assertion]) error {
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

func (m *MockRepository) GetLatestRunByTestVersion(_ context.Context, testID id.ID, version int) (model.Run, error) {
	args := m.Called(testID, version)
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

func (m *MockRepository) CreateTransaction(_ context.Context, t transaction.Transaction) (transaction.Transaction, error) {
	args := m.Called(t)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockRepository) UpdateTransaction(_ context.Context, t transaction.Transaction) (transaction.Transaction, error) {
	args := m.Called(t)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockRepository) DeleteTransaction(_ context.Context, transaction transaction.Transaction) error {
	args := m.Called(transaction)
	return args.Error(1)
}

func (m *MockRepository) GetLatestTransactionVersion(_ context.Context, id id.ID) (transaction.Transaction, error) {
	args := m.Called(id)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockRepository) TransactionIDExists(_ context.Context, id id.ID) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetTransactionVersion(_ context.Context, id id.ID, version int) (transaction.Transaction, error) {
	args := m.Called(id, version)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockRepository) GetTransactions(_ context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[transaction.Transaction], error) {
	args := m.Called(take, skip, query, sortBy, sortDirection)
	transactions := args.Get(0).([]transaction.Transaction)
	list := model.List[transaction.Transaction]{
		Items:      transactions,
		TotalCount: len(transactions),
	}
	return list, args.Error(1)
}

// DeleteTransactionRun implements model.Repository
func (m *MockRepository) DeleteTransactionRun(ctx context.Context, run transaction.TransactionRun) error {
	args := m.Called(ctx, run)
	return args.Error(0)
}

// GetTransactionRun implements model.Repository
func (m *MockRepository) GetTransactionRun(ctx context.Context, transactionID id.ID, runID int) (transaction.TransactionRun, error) {
	args := m.Called(ctx, transactionID, runID)
	return args.Get(0).(transaction.TransactionRun), args.Error(1)
}

func (m *MockRepository) GetLatestRunByTransactionVersion(_ context.Context, transactionID id.ID, version int) (transaction.TransactionRun, error) {
	args := m.Called(transactionID, version)
	return args.Get(0).(transaction.TransactionRun), args.Error(1)
}

// GetTransactionsRuns implements model.Repository
func (m *MockRepository) GetTransactionsRuns(ctx context.Context, transactionID id.ID, take, skip int32) ([]transaction.TransactionRun, error) {
	args := m.Called(ctx, transactionID, take, skip)
	return args.Get(0).([]transaction.TransactionRun), args.Error(1)
}

// UpdateTransactionRun implements model.Repository
func (m *MockRepository) UpdateTransactionRun(ctx context.Context, run transaction.TransactionRun) error {
	args := m.Called(ctx, run)
	return args.Error(0)
}

func (m *MockRepository) CreateTransactionRun(ctx context.Context, run transaction.TransactionRun) error {
	args := m.Called(ctx, run)
	return args.Error(0)
}

func (m *MockRepository) CreateTestRunEvent(_ context.Context, event model.TestRunEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockRepository) GetTestRunEvents(_ context.Context, testID id.ID, runID int) ([]model.TestRunEvent, error) {
	args := m.Called(testID, runID)
	return args.Get(0).([]model.TestRunEvent), args.Error(1)
}
