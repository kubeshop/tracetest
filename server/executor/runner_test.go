package executor_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPersistentRunner(t *testing.T) {
	t.Run("TestIsTriggerd", func(t *testing.T) {
		t.Parallel()

		test := model.Test{
			ID:               id.ID("test1"),
			ServiceUnderTest: sampleTrigger,
		}

		f := runnerSetup(t)
		f.expectSuccessExec(test)

		f.run([]model.Test{test}, 10*time.Millisecond)

		result := f.mockDB.runs[test.ID]
		require.NotNil(t, result)
		assert.Greater(t, result.ServiceTriggerCompletedAt.UnixNano(), result.CreatedAt.UnixNano())

		f.assert(t)
	})

	t.Run("TestsCanBeTriggerdConcurrently", func(t *testing.T) {
		t.Parallel()

		test1 := model.Test{ID: id.ID("test1"), ServiceUnderTest: sampleTrigger}
		test2 := model.Test{ID: id.ID("test2"), ServiceUnderTest: sampleTrigger}

		f := runnerSetup(t)
		f.expectSuccessExecLong(test1)
		f.expectSuccessExec(test2)

		f.run([]model.Test{test1, test2}, 100*time.Millisecond)

		run1 := f.mockDB.runs[test1.ID]
		run2 := f.mockDB.runs[test2.ID]

		assert.Greater(t, run1.ServiceTriggerCompletedAt.UnixNano(), run2.ServiceTriggerCompletedAt.UnixNano(), "test1 did not complete after test2")
	})

}

var (
	noError error = nil

	sampleResponse = trigger.Response{
		SpanAttributes: map[string]string{
			"tracetest.run.trigger.http.response_code": "200",
		},
		Result: model.TriggerResult{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPResponse{
				StatusCode: 200,
				Body:       "this is the body",
				Headers: []model.HTTPHeader{
					{Key: "Content-Type", Value: "text/plain"},
				},
			},
		},
	}

	sampleTrigger = model.Trigger{
		Type: model.TriggerTypeHTTP,
	}
)

type runnerFixture struct {
	runner          executor.PersistentRunner
	mockExecutor    *mockTriggerer
	mockDB          *mockDB
	mockTracePoller *mockTracePoller
}

func (f runnerFixture) run(tests []model.Test, ttl time.Duration) {
	f.runner.Start(2)
	time.Sleep(10 * time.Millisecond)
	for _, test := range tests {
		f.runner.Run(context.TODO(), test, model.RunMetadata{}, environment.Environment{})
	}
	time.Sleep(ttl)
	f.runner.Stop()
}

func (f runnerFixture) expectSuccessExecLong(test model.Test) {
	f.mockExecutor.expectTriggerTestLong(test)
	f.expectSuccessResultPersist(test)
}

func (f runnerFixture) expectSuccessExec(test model.Test) {
	f.mockExecutor.expectTriggerTest(test)
	f.expectSuccessResultPersist(test)
}

func (f runnerFixture) expectSuccessResultPersist(test model.Test) {
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
	reg := trigger.NewRegsitry(tr, tr)

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
		mtp,
		tracer,
		subscription.NewManager(),
		tracedb.Factory(&testDB),
		getDataStoreRepositoryMock(t),
		eventEmitter,
		defaultProfileGetter{5 * time.Second, 30 * time.Second},
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

	runs map[id.ID]model.Run
}

func (m *mockDB) CreateRun(_ context.Context, test model.Test, run model.Run) (model.Run, error) {
	args := m.Called(test.ID)
	if m.runs == nil {
		m.runs = map[id.ID]model.Run{}
	}

	run.ID = rand.Intn(100)
	m.runs[test.ID] = run

	return run, args.Error(0)
}

func (m *mockDB) UpdateRun(_ context.Context, run model.Run) error {
	args := m.Called(run.ID)
	for k, v := range m.runs {
		if v.ID == run.ID {
			m.runs[k] = run
		}
	}

	return args.Error(0)
}

type mockTriggerer struct {
	mock.Mock
	t *testing.T
}

func (m *mockTriggerer) Type() model.TriggerType {
	return model.TriggerTypeHTTP
}

func (m *mockTriggerer) Trigger(_ context.Context, test model.Test, opts *trigger.TriggerOptions) (trigger.Response, error) {
	args := m.Called(test.ID)
	return args.Get(0).(trigger.Response), args.Error(1)
}

func (m *mockTriggerer) Resolve(_ context.Context, test model.Test, opts *trigger.TriggerOptions) (model.Test, error) {
	args := m.Called(test.ID)
	return args.Get(0).(model.Test), args.Error(1)
}

func (m *mockTriggerer) expectTriggerTest(test model.Test) *mock.Call {
	return m.
		On("Resolve", test.ID).
		Return(test, noError).
		On("Trigger", test.ID).
		Return(sampleResponse, noError)
}

func (m *mockTriggerer) expectTriggerTestLong(test model.Test) *mock.Call {
	return m.
		On("Trigger", test.ID).
		After(50*time.Millisecond).
		Return(sampleResponse, noError).
		On("Resolve", test.ID).
		Return(test, noError)
}

func expectCreateRun(m *mockDB, test model.Test) *mock.Call {
	return m.
		On("CreateRun", test.ID).
		Return(noError)
}

type mockTracePoller struct {
	mock.Mock
	t *testing.T
}

func (m *mockTracePoller) Poll(_ context.Context, test model.Test, run model.Run, pollingProfile pollingprofile.PollingProfile) {
	m.Called(test.ID)
}

func (m *mockTracePoller) expectPoll(test model.Test) *mock.Call {
	return m.
		On("Poll", test.ID)
}
