package executor_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	triggerer "github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type defaultTestRunnerGetter struct{}

func (dpc defaultTestRunnerGetter) GetDefault(context.Context) testrunner.TestRunner {
	testRunner := testrunner.DefaultTestRunner

	return testRunner
}

func TestPersistentRunner(t *testing.T) {
	t.Run("TestIsTriggerd", func(t *testing.T) {
		t.Parallel()

		testObj := test.Test{
			ID:      id.ID("test1"),
			Trigger: sampleTrigger,
		}

		f := runnerSetup(t)
		f.expectSuccessExec(testObj)

		f.run([]test.Test{testObj}, 10*time.Millisecond)

		result := f.mockDB.runs[testObj.ID]
		require.NotNil(t, result)
		assert.Greater(t, result.ServiceTriggerCompletedAt.UnixNano(), result.CreatedAt.UnixNano())

		f.assert(t)
	})

	t.Run("TestsCanBeTriggerdConcurrently", func(t *testing.T) {
		t.Parallel()

		test1 := test.Test{ID: id.ID("test1"), Trigger: sampleTrigger}
		test2 := test.Test{ID: id.ID("test2"), Trigger: sampleTrigger}

		f := runnerSetup(t)
		f.expectSuccessExecLong(test1)
		f.expectSuccessExec(test2)

		f.run([]test.Test{test1, test2}, 100*time.Millisecond)

		run1 := f.mockDB.runs[test1.ID]
		run2 := f.mockDB.runs[test2.ID]

		assert.Greater(t, run1.ServiceTriggerCompletedAt.UnixNano(), run2.ServiceTriggerCompletedAt.UnixNano(), "test1 did not complete after test2")
	})

}

var (
	noError error = nil

	sampleResponse = triggerer.Response{
		SpanAttributes: map[string]string{
			"tracetest.run.trigger.http.response_code": "200",
		},
		Result: trigger.TriggerResult{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPResponse{
				StatusCode: 200,
				Body:       "this is the body",
				Headers: []trigger.HTTPHeader{
					{Key: "Content-Type", Value: "text/plain"},
				},
			},
		},
	}

	sampleTrigger = trigger.Trigger{
		Type: trigger.TriggerTypeHTTP,
	}
)

type runnerFixture struct {
	runner          executor.PersistentRunner
	mockExecutor    *mockTriggerer
	mockDB          *mockDB
	mockTracePoller *mockTracePoller
}

func (f runnerFixture) run(tests []test.Test, ttl time.Duration) {
	// TODO - fix this test
	// f.runner.Start(2)
	// time.Sleep(10 * time.Millisecond)
	// for _, testObj := range tests {
	// 	f.runner.Run(context.TODO(), testObj, test.RunMetadata{}, environment.Environment{}, nil)
	// }
	// time.Sleep(ttl)
	// f.runner.Stop()
}

func (f runnerFixture) expectSuccessExecLong(test test.Test) {
	f.mockExecutor.expectTriggerTestLong(test)
	f.expectSuccessResultPersist(test)
}

func (f runnerFixture) expectSuccessExec(test test.Test) {
	f.mockExecutor.expectTriggerTest(test)
	f.expectSuccessResultPersist(test)
}

func (f runnerFixture) expectSuccessResultPersist(test test.Test) {
	expectCreateRun(f.mockDB, test)
	f.mockDB.On("UpdateRun", mock.Anything).Return(noError)
	f.mockDB.On("UpdateRun", mock.Anything).Return(noError)
	f.mockTracePoller.expectPoll(test)
}

func (f runnerFixture) assert(t *testing.T) {
	f.mockExecutor.AssertExpectations(t)
	f.mockDB.AssertExpectations(t)
}

func runnerSetup(t *testing.T) runnerFixture {
	tr, _ := tracing.NewTracer(context.TODO(), config.Must(config.New()))
	reg := triggerer.NewRegsitry(tr, tr)

	me := new(mockTriggerer)
	me.t = t
	me.Test(t)
	reg.Add(me)

	db := new(mockDB)
	db.T = t
	db.Test(t)

	mtp := new(mockTracePoller)
	mtp.t = t

	tracer, _ := tracing.NewTracer(context.Background(), config.Must(config.New()))

	testDB := testdb.MockRepository{}
	testDB.Mock.On("CreateTestRunEvent", mock.Anything).Return(noError)

	eventEmitter := executor.NewEventEmitter(&testDB, subscription.NewManager())

	persistentRunner := executor.NewPersistentRunner(
		reg,
		db,
		executor.NewDBUpdater(db),
		tracer,
		subscription.NewManager(),
		tracedb.Factory(db),
		getDataStoreRepositoryMock(t),
		eventEmitter,
	)

	mtp.Test(t)
	return runnerFixture{
		runner:          persistentRunner,
		mockExecutor:    me,
		mockDB:          db,
		mockTracePoller: mtp,
	}
}

type mockDB struct {
	testdb.MockRepository

	runs map[id.ID]test.Run
}

var _ test.RunRepository = &mockDB{}

func (m *mockDB) CreateRun(_ context.Context, testObj test.Test, run test.Run) (test.Run, error) {
	args := m.Called(testObj.ID)
	if m.runs == nil {
		m.runs = map[id.ID]test.Run{}
	}

	run.ID = rand.Intn(100)
	m.runs[testObj.ID] = run

	return run, args.Error(0)
}

func (m *mockDB) UpdateRun(_ context.Context, run test.Run) error {
	args := m.Called(run.ID)
	for k, v := range m.runs {
		if v.ID == run.ID {
			m.runs[k] = run
		}
	}

	return args.Error(0)
}

func (m *mockDB) GetTransactionRunSteps(ctx context.Context, id id.ID, runID int) ([]test.Run, error) {
	args := m.Called(ctx, id, runID)
	return args.Get(0).([]test.Run), args.Error(1)
}

type mockRunRepository struct {
	mock.Mock
}

type mockTriggerer struct {
	mock.Mock
	t *testing.T
}

func (m *mockTriggerer) Type() trigger.TriggerType {
	return trigger.TriggerTypeHTTP
}

func (m *mockTriggerer) Trigger(_ context.Context, test test.Test, opts *triggerer.TriggerOptions) (triggerer.Response, error) {
	args := m.Called(test.ID)
	return args.Get(0).(triggerer.Response), args.Error(1)
}

func (m *mockTriggerer) Resolve(_ context.Context, testObj test.Test, opts *triggerer.TriggerOptions) (test.Test, error) {
	args := m.Called(testObj.ID)
	return args.Get(0).(test.Test), args.Error(1)
}

func (m *mockTriggerer) expectTriggerTest(test test.Test) *mock.Call {
	return m.
		On("Resolve", test.ID).
		Return(test, noError).
		On("Trigger", test.ID).
		Return(sampleResponse, noError)
}

func (m *mockTriggerer) expectTriggerTestLong(test test.Test) *mock.Call {
	return m.
		On("Trigger", test.ID).
		After(50*time.Millisecond).
		Return(sampleResponse, noError).
		On("Resolve", test.ID).
		Return(test, noError)
}

func expectCreateRun(m *mockDB, test test.Test) *mock.Call {
	return m.
		On("CreateRun", test.ID).
		Return(noError)
}

type mockTracePoller struct {
	mock.Mock
	t *testing.T
}

func (m *mockTracePoller) Poll(_ context.Context, test test.Test, run test.Run, pollingProfile pollingprofile.PollingProfile) {
	m.Called(test.ID)
}

func (m *mockTracePoller) expectPoll(test test.Test) *mock.Call {
	return m.
		On("Poll", test.ID)
}
