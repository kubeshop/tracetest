package executor_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/tracing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/trace"
)

var (
	randomIDGenerator = id.NewRandGenerator()
)

func Test_PollerExecutor_ExecuteRequest_NoRootSpan_NoSpanCase(t *testing.T) {
	t.Parallel()

	// Scenario: Trace without any spans, even root span
	// Given the trigger execution returns 0 spans
	// And tracetest does not send the root span
	// When the server do the polling process
	// Then it will not send a finished flag
	// And it will return a connection error on every call

	// Given conditions

	// maxRetries=30 (inferred by the calculation: maxWaitTimeForTrace / retryDelay)
	retryDelay := 1 * time.Second
	maxWaitTimeForTrace := 30 * time.Second

	tracePerIteration := []model.Trace{
		{},
		{},
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, retryDelay, maxWaitTimeForTrace, pollerExecutor, []iterationExpectedValues{
		{finished: false, expectNoTraceError: true},
		{finished: false, expectNoTraceError: true},
	})

	// it will return errors on repeated calls.
	// on this case, the trace polling process will be finished by TracePoller itself
}

func Test_PollerExecutor_ExecuteRequest_NoRootSpan_OneSpanCase(t *testing.T) {
	t.Parallel()

	// Scenario: Trace with only 1 span, without root span
	// Given the trigger execution returns 1 span on the second iteration
	// And find no trace on the first iteration
	// And tracetest does not send the root span
	// When the server do the polling process
	// Then it should stop at the third iteration
	// And it should handle the trace error on first iteration
	// And a root span should be added to it

	// Given conditions

	// maxRetries=30 (inferred by the calculation: maxWaitTimeForTrace / retryDelay)
	retryDelay := 1 * time.Second
	maxWaitTimeForTrace := 30 * time.Second

	trace := model.NewTrace(randomIDGenerator.TraceID().String(), []model.Span{
		{
			ID:        randomIDGenerator.SpanID(),
			Name:      "HTTP API",
			StartTime: time.Now(),
			EndTime:   time.Now().Add(retryDelay),
			Attributes: map[string]string{
				"testSpan": "true",
			},
			Children: []*model.Span{},
		},
	})

	// test
	tracePerIteration := []model.Trace{
		{},
		trace,
		trace,
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, retryDelay, maxWaitTimeForTrace, pollerExecutor, []iterationExpectedValues{
		{finished: false, expectNoTraceError: true},
		{finished: false, expectNoTraceError: false, expectRootSpan: false},
		{finished: true, expectNoTraceError: false, expectRootSpan: true},
	})
}

func Test_PollerExecutor_ExecuteRequest_NoRootSpan_TwoSpansCase(t *testing.T) {
	t.Parallel()

	// Scenario: Trace with 2 span, without root span
	// Given the trigger execution returns 1 span on second iteration and another one on third iteration
	// And find no trace on the first iteration
	// And tracetest does not send the root span
	// When the server do the polling process
	// Then it should stop at the fourth iteration
	// And it should handle the trace error on first iteration
	// And a root span should be added to it

	// Given conditions

	// maxRetries=30 (inferred by the calculation: maxWaitTimeForTrace / retryDelay)
	retryDelay := 1 * time.Second
	maxWaitTimeForTrace := 30 * time.Second

	traceID := randomIDGenerator.TraceID().String()

	firstSpan := model.Span{
		ID:        randomIDGenerator.SpanID(),
		Name:      "HTTP API",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(retryDelay),
		Attributes: map[string]string{
			"testSpan": "true",
		},
		Children: []*model.Span{},
	}

	secondSpan := model.Span{
		ID:        randomIDGenerator.SpanID(),
		Name:      "Database query",
		StartTime: firstSpan.EndTime,
		EndTime:   firstSpan.EndTime.Add(retryDelay),
		Attributes: map[string]string{
			"testSpan":                           "true",
			model.TracetestMetadataFieldParentID: firstSpan.ID.String(),
		},
		Children: []*model.Span{},
	}

	traceWithOneSpan := model.NewTrace(traceID, []model.Span{firstSpan})
	traceWithTwoSpans := model.NewTrace(traceID, []model.Span{firstSpan, secondSpan})

	// test
	tracePerIteration := []model.Trace{
		{},
		traceWithOneSpan,
		traceWithTwoSpans,
		traceWithTwoSpans,
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, retryDelay, maxWaitTimeForTrace, pollerExecutor, []iterationExpectedValues{
		{finished: false, expectNoTraceError: true},
		{finished: false, expectNoTraceError: false, expectRootSpan: false},
		{finished: false, expectNoTraceError: false, expectRootSpan: false},
		{finished: true, expectNoTraceError: false, expectRootSpan: true},
	})
}

func Test_PollerExecutor_ExecuteRequest_WithRootSpan_NoSpanCase(t *testing.T) {
	t.Parallel()

	// Scenario: Trace without any spans, only root span
	// Given the trigger execution returns 0 spans
	// And tracetest sent the root span
	// When the server do the polling process
	// Then it should stop on third iteration
	// And it should handle the trace error on first iteration

	// Given conditions

	// maxRetries=3 (inferred by the calculation: maxWaitTimeForTrace / retryDelay)
	retryDelay := 1 * time.Second
	maxWaitTimeForTrace := 3 * time.Second

	rootSpan := model.Span{
		ID:        randomIDGenerator.SpanID(),
		Name:      model.TriggerSpanName,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(retryDelay),
		Attributes: map[string]string{
			"testSpan": "true",
		},
		Children: []*model.Span{},
	}

	trace := model.NewTrace(randomIDGenerator.TraceID().String(), []model.Span{rootSpan})

	tracePerIteration := []model.Trace{
		{},
		trace,
		trace,
		trace,
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, retryDelay, maxWaitTimeForTrace, pollerExecutor, []iterationExpectedValues{
		{finished: false, expectNoTraceError: true},
		{finished: false, expectNoTraceError: false, expectRootSpan: true},
		{finished: false, expectNoTraceError: false, expectRootSpan: true},
		{finished: true, expectNoTraceError: false, expectRootSpan: true},
	})
}

func Test_PollerExecutor_ExecuteRequest_WithRootSpan_OneSpanCase(t *testing.T) {
	t.Parallel()

	// Scenario: Trace with only 1 span, plus a root span
	// Given the trigger execution returns 1 span on second iteration
	// And find no trace on the first iteration
	// And tracetest sent the root span
	// When the server do the polling process
	// Then it should stop at the second iteration
	// And it should handle the trace error on first iteration
	// And a root span should be added to it

	// Given conditions

	// maxRetries=30 (inferred by the calculation: maxWaitTimeForTrace / retryDelay)
	retryDelay := 1 * time.Second
	maxWaitTimeForTrace := 30 * time.Second

	rootSpanID := randomIDGenerator.SpanID()

	trace := model.NewTrace(randomIDGenerator.TraceID().String(), []model.Span{
		{
			ID:        rootSpanID,
			Name:      model.TriggerSpanName,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(retryDelay),
			Attributes: map[string]string{
				"testSpan": "true",
			},
			Children: []*model.Span{},
		},
		{
			ID:        randomIDGenerator.SpanID(),
			Name:      "HTTP API",
			StartTime: time.Now(),
			EndTime:   time.Now().Add(retryDelay),
			Attributes: map[string]string{
				"testSpan":                           "true",
				model.TracetestMetadataFieldParentID: rootSpanID.String(),
			},
			Children: []*model.Span{},
		},
	})

	// test
	tracePerIteration := []model.Trace{
		{},
		trace,
		trace,
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, retryDelay, maxWaitTimeForTrace, pollerExecutor, []iterationExpectedValues{
		{finished: false, expectNoTraceError: true},
		{finished: false, expectNoTraceError: false, expectRootSpan: true},
		{finished: true, expectNoTraceError: false, expectRootSpan: true},
	})
}

func Test_PollerExecutor_ExecuteRequest_WithRootSpan_OneDelayedSpanCase(t *testing.T) {
	t.Parallel()

	// Scenario: Trace with only 1 delayed span, plus a root span
	// Given the trigger execution returns 1 span on fourth iteration
	// And find no trace on the first iteration
	// And tracetest sent the root span
	// When the server do the polling process
	// Then it should stop at the fifth iteration
	// And it should handle the trace error on first iteration
	// And a root span should be added to it

	// Given conditions

	// maxRetries=30 (inferred by the calculation: maxWaitTimeForTrace / retryDelay)
	retryDelay := 1 * time.Second
	maxWaitTimeForTrace := 30 * time.Second

	rootSpan := model.Span{
		ID:        randomIDGenerator.SpanID(),
		Name:      model.TriggerSpanName,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(retryDelay),
		Attributes: map[string]string{
			"testSpan": "true",
		},
		Children: []*model.Span{},
	}

	apiSpan := model.Span{
		ID:        randomIDGenerator.SpanID(),
		Name:      "HTTP API",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(retryDelay),
		Attributes: map[string]string{
			"testSpan":                           "true",
			model.TracetestMetadataFieldParentID: rootSpan.ID.String(),
		},
		Children: []*model.Span{},
	}

	traceWithOnlyRoot := model.NewTrace(randomIDGenerator.TraceID().String(), []model.Span{rootSpan})
	completeTrace := model.NewTrace(randomIDGenerator.TraceID().String(), []model.Span{rootSpan, apiSpan})

	// test
	tracePerIteration := []model.Trace{
		{},
		traceWithOnlyRoot,
		traceWithOnlyRoot,
		completeTrace,
		completeTrace,
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, retryDelay, maxWaitTimeForTrace, pollerExecutor, []iterationExpectedValues{
		{finished: false, expectNoTraceError: true},
		{finished: false, expectNoTraceError: false, expectRootSpan: true},
		{finished: false, expectNoTraceError: false, expectRootSpan: true},
		{finished: false, expectNoTraceError: false, expectRootSpan: true},
		{finished: true, expectNoTraceError: false, expectRootSpan: true},
	})
}

func Test_PollerExecutor_ExecuteRequest_WithRootSpan_TwoSpansCase(t *testing.T) {
	t.Parallel()

	// Scenario: Trace with 2 span, plus a root span
	// Given the trigger execution returns 1 span on second iteration and another one on third iteration
	// And find no trace on the first iteration
	// And tracetest sent the root span
	// When the server do the polling process
	// Then it should stop at the third iteration
	// And it should handle the trace error on first iteration
	// And a root span should be added to it

	// Given conditions

	// maxRetries=30 (inferred by the calculation: maxWaitTimeForTrace / retryDelay)
	retryDelay := 1 * time.Second
	maxWaitTimeForTrace := 30 * time.Second

	traceID := randomIDGenerator.TraceID().String()

	rootSpan := model.Span{
		ID:        randomIDGenerator.SpanID(),
		Name:      model.TriggerSpanName,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(retryDelay),
		Attributes: map[string]string{
			"testSpan": "true",
		},
		Children: []*model.Span{},
	}

	firstSpan := model.Span{
		ID:        randomIDGenerator.SpanID(),
		Name:      "HTTP API",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(retryDelay),
		Attributes: map[string]string{
			"testSpan":                           "true",
			model.TracetestMetadataFieldParentID: rootSpan.ID.String(),
		},
		Children: []*model.Span{},
	}

	secondSpan := model.Span{
		ID:        randomIDGenerator.SpanID(),
		Name:      "Database query",
		StartTime: firstSpan.EndTime,
		EndTime:   firstSpan.EndTime.Add(retryDelay),
		Attributes: map[string]string{
			"testSpan":                           "true",
			model.TracetestMetadataFieldParentID: firstSpan.ID.String(),
		},
		Children: []*model.Span{},
	}

	traceWithOneSpan := model.NewTrace(traceID, []model.Span{rootSpan, firstSpan})
	traceWithTwoSpans := model.NewTrace(traceID, []model.Span{rootSpan, firstSpan, secondSpan})

	// test
	tracePerIteration := []model.Trace{
		{},
		traceWithOneSpan,
		traceWithTwoSpans,
		traceWithTwoSpans,
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, retryDelay, maxWaitTimeForTrace, pollerExecutor, []iterationExpectedValues{
		{finished: false, expectNoTraceError: true},
		{finished: false, expectNoTraceError: false, expectRootSpan: true},
		{finished: false, expectNoTraceError: false, expectRootSpan: true},
		{finished: true, expectNoTraceError: false, expectRootSpan: true},
	})
}

// Helper structs / functions

type iterationExpectedValues struct {
	finished           bool
	expectNoTraceError bool
	expectRootSpan     bool
}

func executeAndValidatePollingRequests(t *testing.T, retryDelay, maxWaitTimeForTrace time.Duration, pollerExecutor *executor.InstrumentedPollerExecutor, expectedValues []iterationExpectedValues) {
	ctx := context.Background()
	run := test.NewRun()

	test := test.Test{
		ID: id.ID("some-test"),
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
		},
	}

	// using `pollingprofile.DefaultPollingProfile` and changing the periodic configs directly
	// causes a race condition because `DefaultPollingProfile.Periodic` is a reference, so it's shared by the copies.
	// The easiest solution is to create a new polling profile for each test.
	pp := pollingprofile.PollingProfile{
		Strategy: pollingprofile.Periodic,
		Periodic: &pollingprofile.PeriodicPollingConfig{
			RetryDelay: retryDelay.String(),
			Timeout:    maxWaitTimeForTrace.String(),
		},
	}

	job := executor.NewJob()
	job.Test = test
	job.Run = run
	job.PollingProfile = pp

	for i, expected := range expectedValues {
		res, err := pollerExecutor.ExecuteRequest(ctx, &job)
		run = res.Run() // should store a run to use in another iteration

		require.NotNilf(t, run, "The test run should not be nil on iteration %d", i)

		if expected.expectNoTraceError {
			require.Errorf(t, err, "An error should have happened on iteration %d", i)
			require.ErrorIsf(t, err, connection.ErrTraceNotFound, "An connection error should have happened on iteration %d", i)
		} else {
			require.NoErrorf(t, err, "An error should not have happened on iteration %d", i)
			require.NotEmptyf(t, res.Reason(), "The poller should have a reason on iteration %d", i)

			if expected.finished {
				require.Truef(t, res.Finished(), "The poller should have finished on iteration %d", i)
			} else {
				require.Falsef(t, res.Finished(), "The poller should have not finished on iteration %d", i)
			}

			// only validate root span if we have a root span
			if expected.expectRootSpan {
				require.Truef(t, run.Trace.HasRootSpan(), "The trace associated with the run on iteration %d should have a root span", i)
			} else {
				require.Falsef(t, run.Trace.HasRootSpan(), "The trace associated with the run on iteration %d should not have a root span", i)
			}
		}

		job.IncreaseEnqueueCount()
		job.Headers.SetBool("requeued", true)
	}
}

func getPollerExecutorWithMocks(t *testing.T, tracePerIteration []model.Trace) *executor.InstrumentedPollerExecutor {
	updater := getRunUpdaterMock(t)
	tracer := getTracerMock(t)
	testDB := getRunRepositoryMock(t)
	dataStoreRepo := getDataStoreRepositoryMock(t)
	traceDBFactory := getTraceDBMockFactory(t, tracePerIteration, &traceDBState{})
	eventEmitter := getEventEmitterMock(t, testDB)

	return executor.NewPollerExecutor(
		tracer,
		updater,
		traceDBFactory,
		dataStoreRepo,
		eventEmitter,
	)
}

// Mocks

// RunUpdater
type runUpdaterMock struct{}

func (m runUpdaterMock) Update(context.Context, test.Run) error { return nil }

func getRunUpdaterMock(_ *testing.T) executor.RunUpdater {
	return runUpdaterMock{}
}

// RunRepository
func getRunRepositoryMock(t *testing.T) *testdb.MockRepository {
	t.Helper()

	testDB := new(testdb.MockRepository)
	testDB.Test(t)
	testDB.On("CreateTestRunEvent", mock.Anything).Return(noError)

	return testDB
}

// DataStoreRepository
type dataStoreRepositoryMock struct {
	mock.Mock
}

func (m *dataStoreRepositoryMock) Current(context.Context) (datastore.DataStore, error) {
	return datastore.DataStore{Type: datastore.DataStoreTypeOTLP}, nil
}

func (m *dataStoreRepositoryMock) Get(_ context.Context, id id.ID) (datastore.DataStore, error) {
	args := m.Called(id)
	return args.Get(0).(datastore.DataStore), args.Error(1)
}

func getDataStoreRepositoryMock(t *testing.T) *dataStoreRepositoryMock {
	m := new(dataStoreRepositoryMock)
	m.Test(t)
	return m
}

// EventEmitter
func getEventEmitterMock(t *testing.T, db *testdb.MockRepository) executor.EventEmitter {
	t.Helper()

	return executor.NewEventEmitter(db, subscription.NewManager())
}

// Tracer
func getTracerMock(t *testing.T) trace.Tracer {
	t.Helper()

	tracer, err := tracing.NewTracer(context.TODO(), config.Must(config.New()))
	require.NoError(t, err)

	return tracer
}

// TraceDB
type traceDBMock struct {
	tracePerIteration []model.Trace
	state             *traceDBState
}

func (db *traceDBMock) GetTraceByID(_ context.Context, _ string) (t model.Trace, err error) {
	trace := db.tracePerIteration[db.state.currentIteration]
	db.state.currentIteration += 1

	if len(trace.Flat) == 0 {
		return trace, connection.ErrTraceNotFound
	}

	return trace, nil
}

func (db *traceDBMock) ShouldRetry() bool {
	return true // this provider should retry
}

func (db *traceDBMock) GetTraceID() trace.TraceID {
	return randomIDGenerator.TraceID()
}

// empty methods to respect TraceDB interface
func (db *traceDBMock) Connect(ctx context.Context) error { return nil }
func (db *traceDBMock) Close() error                      { return nil }
func (db *traceDBMock) Ready() bool                       { return true }
func (db *traceDBMock) GetEndpoints() string              { return "" }
func (db *traceDBMock) TestConnection(ctx context.Context) model.ConnectionResult {
	return model.ConnectionResult{}
}

type traceDBState struct {
	currentIteration int
}

func getTraceDBMockFactory(t *testing.T, tracePerIteration []model.Trace, state *traceDBState) func(datastore.DataStore) (tracedb.TraceDB, error) {
	t.Helper()

	return func(ds datastore.DataStore) (tracedb.TraceDB, error) {
		return &traceDBMock{
			tracePerIteration: tracePerIteration,
			state:             state,
		}, nil
	}
}
