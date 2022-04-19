package openapi_test

import (
	"context"
	"testing"
	"time"

	openapi "github.com/kubeshop/tracetest/server/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestPersistentRunner(t *testing.T) {
	t.Run("TestIsExecuted", func(t *testing.T) {
		test := openapi.Test{
			TestId: "test",
		}

		f := runnerSetup(t)
		f.expectSuccessExec(test)

		f.run([]openapi.Test{test}, 10*time.Millisecond)

		result := f.mockResultsDB.results[test.TestId]
		require.NotNil(t, result)
		assert.Greater(t, result.CompletedAt.UnixNano(), result.CreatedAt.UnixNano())

		f.assert(t)
	})

	t.Run("TestsCanBeExecutedConcurrently", func(t *testing.T) {
		test1 := openapi.Test{TestId: "test1"}
		test2 := openapi.Test{TestId: "test2"}

		f := runnerSetup(t)

		f.expectSuccessExecLong(test1)
		f.expectSuccessExec(test2)

		f.run([]openapi.Test{test1, test2}, 5100*time.Millisecond)

		result1 := f.mockResultsDB.results[test1.TestId]
		require.NotNil(t, result1)

		result2 := f.mockResultsDB.results[test2.TestId]
		require.NotNil(t, result2)

		assert.True(t, result1.CompletedAt.UnixNano() > result2.CompletedAt.UnixNano(), "test1 did not complete after test2")
	})

}

var (
	noError error = nil

	sampleResponse = openapi.HttpResponse{
		StatusCode: 200,
		Body:       "this is the body",
		Headers: []openapi.HttpResponseHeaders{
			{"Content-Type", "text/plain"},
		},
	}
)

type runnerFixture struct {
	runner          openapi.PersistentRunner
	mockExecutor    *mockExecutor
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
	f.mockResultsDB.expectCreateResult(test)
	f.mockResultsDB.expectUpdateResultState(test, openapi.TestRunStateExecuting)
	f.mockResultsDB.expectUpdateResultState(test, openapi.TestRunStateAwaitingTrace)
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

	mr := new(mockResultsDB)
	mr.t = t
	mr.Test(t)

	mtp := new(mockTracePoller)
	mtp.t = t
	mtp.Test(t)
	return runnerFixture{
		runner:          openapi.NewPersistentRunner(me, mr, mtp),
		mockExecutor:    me,
		mockResultsDB:   mr,
		mockTracePoller: mtp,
	}
}

type mockExecutor struct {
	mock.Mock
	t *testing.T
}

func (m *mockExecutor) Execute(test *openapi.Test, tid trace.TraceID, sid trace.SpanID) (openapi.HttpResponse, error) {
	args := m.Called(test.TestId)
	return args.Get(0).(openapi.HttpResponse), args.Error(1)
}

func (m *mockExecutor) expectExecuteTest(test openapi.Test) *mock.Call {
	return m.
		On("Execute", test.TestId).
		Return(sampleResponse, noError)
}

func (m *mockExecutor) expectExecuteTestLong(test openapi.Test) *mock.Call {
	return m.
		On("Execute", test.TestId).
		After(5*time.Second).
		Return(sampleResponse, noError)
}

type mockResultsDB struct {
	mock.Mock
	t *testing.T

	results map[string]openapi.TestRunResult
}

func (m *mockResultsDB) CreateResult(ctx context.Context, testID string, res *openapi.TestRunResult) error {
	args := m.Called(res.TestId)
	if m.results == nil {
		m.results = map[string]openapi.TestRunResult{}
	}

	m.results[res.TestId] = *res

	return args.Error(0)
}

func (m *mockResultsDB) UpdateResult(ctx context.Context, res *openapi.TestRunResult) error {
	args := m.Called(res.TestId, res.State)
	if m.results == nil {
		m.results = map[string]openapi.TestRunResult{}
	}

	m.results[res.TestId] = *res

	return args.Error(0)
}

func (m *mockResultsDB) expectCreateResult(test openapi.Test) *mock.Call {
	return m.
		On("CreateResult", test.TestId).
		Return(noError)
}

func (m *mockResultsDB) expectUpdateResultState(test openapi.Test, expectedState string) *mock.Call {
	return m.
		On("UpdateResult", test.TestId, expectedState).
		Return(noError)
}

type mockTracePoller struct {
	mock.Mock
	t *testing.T
}

func (m *mockTracePoller) Poll(_ context.Context, res openapi.TestRunResult) {
	m.Called(res.TestId)
}

func (m *mockTracePoller) expectPoll(test openapi.Test) *mock.Call {
	return m.
		On("Poll", test.TestId)
}
