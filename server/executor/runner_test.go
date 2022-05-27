package executor_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestPersistentRunner(t *testing.T) {
	t.Run("TestIsExecuted", func(t *testing.T) {
		t.Parallel()

		test := model.Test{
			ID: id.NewRandGenerator().UUID(),
		}

		f := runnerSetup(t)
		f.expectSuccessExec(test)

		f.run([]model.Test{test}, 10*time.Millisecond)

		result := f.mockDB.runs[test.ID.String()]
		require.NotNil(t, result)
		assert.Greater(t, result.CompletedAt.UnixNano(), result.CreatedAt.UnixNano())

		f.assert(t)
	})

	t.Run("TestsCanBeExecutedConcurrently", func(t *testing.T) {
		t.Parallel()

		test1 := model.Test{ID: id.NewRandGenerator().UUID()}
		test2 := model.Test{ID: id.NewRandGenerator().UUID()}

		f := runnerSetup(t)

		f.expectSuccessExecLong(test1)
		f.expectSuccessExec(test2)

		f.run([]model.Test{test1, test2}, 100*time.Millisecond)

		run1 := f.mockDB.runs[test1.ID.String()]
		require.NotNil(t, run1)

		run2 := f.mockDB.runs[test2.ID.String()]
		require.NotNil(t, run2)
	})

}

var (
	noError error = nil

	sampleResponse = model.HTTPResponse{
		StatusCode: 200,
		Body:       "this is the body",
		Headers: []model.HTTPHeader{
			{Key: "Content-Type", Value: "text/plain"},
		},
	}
)

type runnerFixture struct {
	runner          executor.PersistentRunner
	mockExecutor    *mockExecutor
	mockDB          *mockDB
	mockTracePoller *mockTracePoller
}

func (f runnerFixture) run(tests []model.Test, ttl time.Duration) {
	f.runner.Start(2)
	time.Sleep(10 * time.Millisecond)
	for _, test := range tests {
		f.runner.Run(test)
	}
	time.Sleep(ttl)
	f.runner.Stop()
}

func (f runnerFixture) expectSuccessExecLong(test model.Test) {
	f.mockExecutor.expectExecuteTestLong(test)
	f.expectSuccessResultPersist(test)
}

func (f runnerFixture) expectSuccessExec(test model.Test) {
	f.mockExecutor.expectExecuteTest(test)
	f.expectSuccessResultPersist(test)
}

func (f runnerFixture) expectSuccessResultPersist(test model.Test) {
	expectCreateRun(f.mockDB, test)
	f.mockDB.On("UpdateRun", mock.Anything).Return(noError)
	f.mockDB.On("UpdateTestVersion", mock.Anything).Return(noError)
	f.mockDB.On("UpdateRun", mock.Anything).Return(noError)
	f.mockTracePoller.expectPoll(test)
}

func (f runnerFixture) assert(t *testing.T) {
	f.mockExecutor.AssertExpectations(t)
	f.mockDB.AssertExpectations(t)
}

func runnerSetup(t *testing.T) runnerFixture {
	me := new(mockExecutor)
	me.t = t
	me.Test(t)

	db := new(mockDB)
	db.T = t
	db.Test(t)

	mtp := new(mockTracePoller)
	mtp.t = t

	mtp.Test(t)
	return runnerFixture{
		runner:          executor.NewPersistentRunner(me, db, mtp),
		mockExecutor:    me,
		mockDB:          db,
		mockTracePoller: mtp,
	}
}

type mockDB struct {
	testdb.MockRepository

	runs map[string]model.Run
}

func (m *mockDB) CreateRun(_ context.Context, test model.Test, run model.Run) (model.Run, error) {
	args := m.Called(test.ID.String())
	if m.runs == nil {
		m.runs = map[string]model.Run{}
	}

	m.runs[test.ID.String()] = run

	return run, args.Error(0)
}

func (m *mockDB) UpdateRun(_ context.Context, run model.Run) error {
	args := m.Called(run.ID.String())
	for k, v := range m.runs {
		if v.ID.String() == run.ID.String() {
			m.runs[k] = run
		}
	}

	return args.Error(0)
}

type mockExecutor struct {
	mock.Mock
	t *testing.T
}

func (m *mockExecutor) Execute(test model.Test, tid trace.TraceID, sid trace.SpanID) (model.HTTPResponse, error) {
	args := m.Called(test.ID)
	return args.Get(0).(model.HTTPResponse), args.Error(1)
}

func (m *mockExecutor) expectExecuteTest(test model.Test) *mock.Call {
	return m.
		On("Execute", test.ID).
		Return(sampleResponse, noError)
}

func (m *mockExecutor) expectExecuteTestLong(test model.Test) *mock.Call {
	return m.
		On("Execute", test.ID).
		After(50*time.Millisecond).
		Return(sampleResponse, noError)
}

func expectCreateRun(m *mockDB, test model.Test) *mock.Call {
	return m.
		On("CreateRun", test.ID.String()).
		Return(noError)
}

type mockTracePoller struct {
	mock.Mock
	t *testing.T
}

func (m *mockTracePoller) Poll(_ context.Context, test model.Test, run model.Run) {
	m.Called(test.ID)
}

func (m *mockTracePoller) expectPoll(test model.Test) *mock.Call {
	return m.
		On("Poll", test.ID)
}
