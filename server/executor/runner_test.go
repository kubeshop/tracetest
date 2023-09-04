package executor_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	triggerer "github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/kubeshop/tracetest/server/tracing"
	"github.com/kubeshop/tracetest/server/variableset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

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

		result := f.runsMock.runs[testObj.ID]
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

		// the worker pool lib reverses the order of this test for some reason
		// but this doesn't matter, as long as both tests run,
		// we only care about the completion time in this test
		f.run([]test.Test{test2, test1}, 100*time.Millisecond)

		run1 := f.runsMock.runs[test1.ID]
		run2 := f.runsMock.runs[test2.ID]

		assert.Greater(t, run1.ServiceTriggerCompletedAt.UnixNano(), run2.ServiceTriggerCompletedAt.UnixNano(), "test1 did not complete after test2")
		f.assert(t)
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
	runner        *executor.TestPipeline
	dsMock        *datastoreGetterMock
	ppMock        *pollingprofileGetterMock
	trMock        *testrunnerGetterMock
	testMock      *testGetterMock
	runsMock      *runsRepoMock
	triggererMock *mockTriggerer
	processorMock *mockProcessor
}

func (f runnerFixture) run(tests []test.Test, ttl time.Duration) {
	// TODO - fix this test
	f.runner.Start()
	time.Sleep(10 * time.Millisecond)
	for _, testObj := range tests {
		newRun := f.runner.Run(context.TODO(), testObj, test.RunMetadata{}, variableset.VariableSet{}, nil)
		// readd this when using not-in-memory queues
		// f.runsMock.
		// 	On("GetRun", testObj.ID, newRun.ID).
		// 	Return(newRun, noError)
		f.processorMock.
			On("ProcessItem", testObj.ID, newRun.ID, datastore.DataStoreSingleID, pollingprofile.DefaultPollingProfile.ID).
			Return(newRun, noError)
	}
	time.Sleep(ttl)
	f.runner.Stop()
}

func (f runnerFixture) assert(t *testing.T) {
	f.dsMock.AssertExpectations(t)
	f.ppMock.AssertExpectations(t)
	f.trMock.AssertExpectations(t)
	f.testMock.AssertExpectations(t)
	f.runsMock.AssertExpectations(t)
	f.triggererMock.AssertExpectations(t)
}

func (f runnerFixture) expectSuccessExecLong(test test.Test) {
	f.triggererMock.expectTriggerTestLong(test)
	f.expectSuccessResultPersist(test)
}

func (f runnerFixture) expectSuccessExec(test test.Test) {
	f.testMock.On("GetAugmented", test.ID).Return(test, noError)
	f.triggererMock.expectTriggerTest(test)
	f.expectSuccessResultPersist(test)
}

func (f runnerFixture) expectSuccessResultPersist(test test.Test) {
	f.testMock.On("GetAugmented", test.ID).Return(test, noError)
	expectCreateRun(f.runsMock, test)
	f.runsMock.On("UpdateRun", mock.Anything).Return(noError)
	f.runsMock.On("UpdateRun", mock.Anything).Return(noError)
}

func runnerSetup(t *testing.T) runnerFixture {

	dsMock := new(datastoreGetterMock)
	dsMock.Test(t)

	ppMock := new(pollingprofileGetterMock)
	ppMock.Test(t)

	trMock := new(testrunnerGetterMock)
	trMock.Test(t)

	testMock := new(testGetterMock)
	testMock.Test(t)

	runsMock := new(runsRepoMock)
	runsMock.Test(t)

	triggererMock := new(mockTriggerer)
	runsMock.Test(t)

	processorMock := new(mockProcessor)
	processorMock.Test(t)

	tracesMock := new(mockTraces)
	tracesMock.Test(t)

	sm := subscription.NewManager()
	tracer, _ := tracing.NewTracer(context.Background(), config.Must(config.New()))
	eventEmitter := executor.NewEventEmitter(getTestRunEventRepositoryMock(t, false), sm)

	registry := triggerer.NewRegistry(tracer, tracer)
	registry.Add(triggererMock)

	runner := executor.NewTriggerResultProcessorWorker(
		tracer,
		sm,
		eventEmitter,
	)

	queueBuilder := executor.NewQueueBuilder().
		WithDataStoreGetter(dsMock).
		WithPollingProfileGetter(ppMock).
		WithTestGetter(testMock).
		WithRunGetter(runsMock)

	pipeline := executor.NewPipeline(queueBuilder,
		executor.PipelineStep{Processor: runner, Driver: executor.NewInMemoryQueueDriver("runner")},
		executor.PipelineStep{Processor: processorMock, Driver: executor.NewInMemoryQueueDriver("runner")},
	)

	testPipeline := executor.NewTestPipeline(
		pipeline,
		nil,
		pipeline.GetQueueForStep(1), // processorMock queue
		runsMock,
		trMock,
		ppMock,
		dsMock,
	)

	return runnerFixture{
		runner:        testPipeline,
		dsMock:        dsMock,
		ppMock:        ppMock,
		trMock:        trMock,
		testMock:      testMock,
		runsMock:      runsMock,
		triggererMock: triggererMock,
		processorMock: processorMock,
	}
}

type mockTraces struct {
	mock.Mock
}

func (r *mockTraces) Get(ctx context.Context, id trace.TraceID) (traces.Trace, error) {
	args := r.Called(id)
	return args.Get(0).(traces.Trace), args.Error(1)
}

func (r *mockTraces) SaveTrace(ctx context.Context, trace *traces.Trace) error {
	args := r.Called(trace)
	return args.Error(0)
}

type mockProcessor struct {
	mock.Mock
}

func (m *mockProcessor) ProcessItem(_ context.Context, job executor.Job) {
	m.Called(job.Test.ID, job.Run.ID, job.DataStore.ID, job.PollingProfile.ID)
}

func (m *mockProcessor) SetOutputQueue(_ executor.Enqueuer) {}

type datastoreGetterMock struct {
	mock.Mock
}

func (r *datastoreGetterMock) Get(ctx context.Context, id id.ID) (datastore.DataStore, error) {
	return r.Current(ctx)
}

func (r *datastoreGetterMock) Current(context.Context) (datastore.DataStore, error) {
	return datastore.DataStore{
		ID:      datastore.DataStoreSingleID,
		Name:    "test",
		Type:    datastore.DataStoreTypeOTLP,
		Default: true,
	}, nil
}

type pollingprofileGetterMock struct {
	mock.Mock
}

func (r *pollingprofileGetterMock) Get(ctx context.Context, _ id.ID) (pollingprofile.PollingProfile, error) {
	return r.GetDefault(ctx), nil
}

func (r *pollingprofileGetterMock) GetDefault(context.Context) pollingprofile.PollingProfile {
	return pollingprofile.DefaultPollingProfile
}

type testGetterMock struct {
	mock.Mock
}

func (r *testGetterMock) GetAugmented(_ context.Context, id id.ID) (test.Test, error) {
	args := r.Called(id)
	return args.Get(0).(test.Test), args.Error(1)
}

type runsRepoMock struct {
	testdb.MockRepository

	runs map[id.ID]test.Run
}

func (m *runsRepoMock) CreateRun(_ context.Context, testObj test.Test, run test.Run) (test.Run, error) {
	args := m.Called(testObj.ID)
	if m.runs == nil {
		m.runs = map[id.ID]test.Run{}
	}

	run.ID = rand.Intn(100)
	m.runs[testObj.ID] = run

	return run, args.Error(0)
}

func (m *runsRepoMock) UpdateRun(_ context.Context, run test.Run) error {
	args := m.Called(run.ID)
	for k, v := range m.runs {
		if v.ID == run.ID {
			m.runs[k] = run
		}
	}

	return args.Error(0)
}

func (r *runsRepoMock) GetRun(_ context.Context, testID id.ID, runID int) (test.Run, error) {
	if run, ok := r.runs[testID]; ok && run.ID == runID {
		return run, nil
	}

	args := r.Called(testID, runID)
	return args.Get(0).(test.Run), args.Error(1)
}

func (r *runsRepoMock) GetRunByTraceID(_ context.Context, id trace.TraceID) (test.Run, error) {
	args := r.Called(id)
	return args.Get(0).(test.Run), args.Error(1)
}

type testrunnerGetterMock struct {
	mock.Mock
}

func (r *testrunnerGetterMock) GetDefault(context.Context) testrunner.TestRunner {
	return testrunner.DefaultTestRunner
}

type mockTriggerer struct {
	mock.Mock
}

func (m *mockTriggerer) Type() trigger.TriggerType {
	return trigger.TriggerTypeHTTP
}

func (m *mockTriggerer) Trigger(_ context.Context, test test.Test) (triggerer.Response, error) {
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

func expectCreateRun(m *runsRepoMock, test test.Test) *mock.Call {
	return m.
		On("CreateRun", test.ID).
		Return(noError)
}
