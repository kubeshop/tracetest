package executor_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/executor"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestPersistentRunner(t *testing.T) {
	t.Run("TestIsExecuted", func(t *testing.T) {
		t.Parallel()

		test := openapi.Test{
			Id: "test",
		}

		f := runnerSetup(t)
		f.expectSuccessExec(test)

		f.run([]openapi.Test{test}, 10*time.Millisecond)

		result := f.mockResultsDB.runs[test.Id]
		require.NotNil(t, result)
		assert.Greater(t, result.CompletedAt.UnixNano(), result.CreatedAt.UnixNano())

		f.assert(t)
	})

	t.Run("TestsCanBeExecutedConcurrently", func(t *testing.T) {
		t.Parallel()

		test1 := openapi.Test{Id: "test1"}
		test2 := openapi.Test{Id: "test2"}

		f := runnerSetup(t)

		f.expectSuccessExecLong(test1)
		f.expectSuccessExec(test2)

		f.run([]openapi.Test{test1, test2}, 100*time.Millisecond)

		run1 := f.mockResultsDB.runs[test1.Id]
		require.NotNil(t, run1)

		run2 := f.mockResultsDB.runs[test2.Id]
		require.NotNil(t, run2)

		assert.True(t, run1.CompletedAt.UnixNano() > run2.CompletedAt.UnixNano(), "test1 did not complete after test2")
	})

}

var (
	noError error = nil

	sampleResponse = openapi.HttpResponse{
		StatusCode: 200,
		Body:       "this is the body",
		Headers: []openapi.HttpHeader{
			{Key: "Content-Type", Value: "text/plain"},
		},
	}
)

type runnerFixture struct {
	runner          executor.PersistentRunner
	mockExecutor    *mockExecutor
	mockTestDB      *mockTestDB
	mockResultsDB   *mockResultsDB
	mockTracePoller *mockTracePoller
}

func (f runnerFixture) run(tests []openapi.Test, ttl time.Duration) {
	f.runner.Start(2)
	time.Sleep(10 * time.Millisecond)
	for _, test := range tests {
		f.runner.Run(test)
	}
	time.Sleep(ttl)
	f.runner.Stop()
}

func (f runnerFixture) expectSuccessExecLong(test openapi.Test) {
	f.mockExecutor.expectExecuteTestLong(test)
	f.expectSuccessResultPersist(test)
}

func (f runnerFixture) expectSuccessExec(test openapi.Test) {
	f.mockExecutor.expectExecuteTest(test)
	f.expectSuccessResultPersist(test)
}

func (f runnerFixture) expectSuccessResultPersist(test openapi.Test) {
	f.mockResultsDB.expectCreateRun(test)
	f.mockResultsDB.expectUpdateRunState(test, executor.TestRunStateExecuting)
	f.mockResultsDB.On("UpdateTest", test.Id).Return(noError)
	f.mockResultsDB.expectUpdateRunState(test, executor.TestRunStateAwaitingTrace)
	f.mockTracePoller.expectPoll(test)
}

func (f runnerFixture) assert(t *testing.T) {
	f.mockExecutor.AssertExpectations(t)
	f.mockResultsDB.AssertExpectations(t)
}

func runnerSetup(t *testing.T) runnerFixture {
	me := new(mockExecutor)
	me.t = t
	me.Test(t)

	mt := new(mockTestDB)
	mt.t = t
	mt.Test(t)

	mr := new(mockResultsDB)
	mr.t = t
	mr.Test(t)

	mtp := new(mockTracePoller)
	mtp.t = t

	mtp.Test(t)
	return runnerFixture{
		runner:          executor.NewPersistentRunner(me, mt, mr, mtp),
		mockExecutor:    me,
		mockTestDB:      mt,
		mockResultsDB:   mr,
		mockTracePoller: mtp,
	}
}

type mockExecutor struct {
	mock.Mock
	t *testing.T
}

func (m *mockExecutor) Execute(test *openapi.Test, tid trace.TraceID, sid trace.SpanID) (openapi.HttpResponse, error) {
	args := m.Called(test.Id)
	return args.Get(0).(openapi.HttpResponse), args.Error(1)
}

func (m *mockExecutor) expectExecuteTest(test openapi.Test) *mock.Call {
	return m.
		On("Execute", test.Id).
		Return(sampleResponse, noError)
}

func (m *mockExecutor) expectExecuteTestLong(test openapi.Test) *mock.Call {
	return m.
		On("Execute", test.Id).
		After(50*time.Millisecond).
		Return(sampleResponse, noError)
}

type mockTestDB struct {
	mock.Mock
	t *testing.T
}

var _ testdb.TestRepository = &mockTestDB{}

func (m *mockTestDB) CreateTest(ctx context.Context, test *openapi.Test) (string, error) {
	args := m.Called(ctx, test)
	return args.String(0), args.Error(1)
}

func (m *mockTestDB) UpdateTest(ctx context.Context, test *openapi.Test) error {
	args := m.Called(ctx, test)
	return args.Error(0)
}

func (m *mockTestDB) DeleteTest(ctx context.Context, test *openapi.Test) error {
	args := m.Called(ctx, test)
	return args.Error(0)
}

func (m *mockTestDB) GetTests(ctx context.Context, take, skip int32) ([]openapi.Test, error) {
	args := m.Called(ctx, take, skip)
	return args.Get(0).([]openapi.Test), args.Error(1)
}

func (m *mockTestDB) GetTest(ctx context.Context, id string) (*openapi.Test, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*openapi.Test), args.Error(1)
}

type mockResultsDB struct {
	mock.Mock
	t *testing.T

	runs map[string]openapi.TestRun
}

func (m *mockResultsDB) CreateRun(ctx context.Context, Id string, res *openapi.TestRun) error {
	args := m.Called(res.Id)
	if m.runs == nil {
		m.runs = map[string]openapi.TestRun{}
	}

	m.runs[res.Id] = *res

	return args.Error(0)
}

func (m *mockResultsDB) GetResult(ctx context.Context, id string) (*openapi.TestRunResult, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*openapi.TestRunResult), args.Error(1)
}

func (m *mockResultsDB) GetResultsByTestID(ctx context.Context, testid string, take, skip int32) ([]openapi.TestRunResult, error) {
	args := m.Called(ctx, testid, take, skip)
	return args.Get(0).([]openapi.TestRunResult), args.Error(1)
}

func (m *mockResultsDB) GetResultByTraceID(ctx context.Context, testid, traceid string) (openapi.TestRunResult, error) {
	args := m.Called(ctx, testid, traceid)
	return args.Get(0).(openapi.TestRunResult), args.Error(1)
}

func (m *mockResultsDB) UpdateTest(_ context.Context, test *openapi.Test) error {
	args := m.Called(test.Id)
	return args.Error(0)
}

func (m *mockResultsDB) UpdateRun(ctx context.Context, res *openapi.TestRun) error {
	args := m.Called(res.Id, res.State)
	if m.runs == nil {
		m.runs = map[string]openapi.TestRun{}
	}

	m.runs[res.Id] = *res

	return args.Error(0)
}

func (m *mockResultsDB) expectCreateRun(test openapi.Test) *mock.Call {
	return m.
		On("CreateRun", test.Id).
		Return(noError)
}

func (m *mockResultsDB) expectUpdateRunState(test openapi.Test, expectedState string) *mock.Call {
	return m.
		On("UpdateRun", test.Id, expectedState).
		Return(noError)
}

type mockTracePoller struct {
	mock.Mock
	t *testing.T
}

func (m *mockTracePoller) Poll(_ context.Context, res openapi.TestRun) {
	m.Called(res.Id)
}

func (m *mockTracePoller) expectPoll(test openapi.Test) *mock.Call {
	return m.
		On("Poll", test.Id)
}
