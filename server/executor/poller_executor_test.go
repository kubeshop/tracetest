package executor_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/tracing"
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
		model.Trace{},
		model.Trace{},
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, retryDelay, maxWaitTimeForTrace, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, pollerExecutor, []iterationExpectedValues{
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
		model.Trace{},
		trace,
		trace,
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, retryDelay, maxWaitTimeForTrace, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, pollerExecutor, []iterationExpectedValues{
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
			"testSpan":  "true",
			"parent_id": firstSpan.ID.String(),
		},
		Children: []*model.Span{},
	}

	traceWithOneSpan := model.NewTrace(traceID, []model.Span{firstSpan})
	traceWithTwoSpans := model.NewTrace(traceID, []model.Span{firstSpan, secondSpan})

	// test
	tracePerIteration := []model.Trace{
		model.Trace{},
		traceWithOneSpan,
		traceWithTwoSpans,
		traceWithTwoSpans,
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, retryDelay, maxWaitTimeForTrace, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, pollerExecutor, []iterationExpectedValues{
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
		model.Trace{},
		trace,
		trace,
		trace,
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, retryDelay, maxWaitTimeForTrace, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, pollerExecutor, []iterationExpectedValues{
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
		model.Span{
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
				"testSpan":  "true",
				"parent_id": rootSpanID.String(),
			},
			Children: []*model.Span{},
		},
	})

	// test
	tracePerIteration := []model.Trace{
		model.Trace{},
		trace,
		trace,
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, retryDelay, maxWaitTimeForTrace, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, pollerExecutor, []iterationExpectedValues{
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
			"testSpan":  "true",
			"parent_id": rootSpan.ID.String(),
		},
		Children: []*model.Span{},
	}

	traceWithOnlyRoot := model.NewTrace(randomIDGenerator.TraceID().String(), []model.Span{rootSpan})
	completeTrace := model.NewTrace(randomIDGenerator.TraceID().String(), []model.Span{rootSpan, apiSpan})

	// test
	tracePerIteration := []model.Trace{
		model.Trace{},
		traceWithOnlyRoot,
		traceWithOnlyRoot,
		completeTrace,
		completeTrace,
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, retryDelay, maxWaitTimeForTrace, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, pollerExecutor, []iterationExpectedValues{
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
			"testSpan":  "true",
			"parent_id": rootSpan.ID.String(),
		},
		Children: []*model.Span{},
	}

	secondSpan := model.Span{
		ID:        randomIDGenerator.SpanID(),
		Name:      "Database query",
		StartTime: firstSpan.EndTime,
		EndTime:   firstSpan.EndTime.Add(retryDelay),
		Attributes: map[string]string{
			"testSpan":  "true",
			"parent_id": firstSpan.ID.String(),
		},
		Children: []*model.Span{},
	}

	traceWithOneSpan := model.NewTrace(traceID, []model.Span{rootSpan, firstSpan})
	traceWithTwoSpans := model.NewTrace(traceID, []model.Span{rootSpan, firstSpan, secondSpan})

	// test
	tracePerIteration := []model.Trace{
		model.Trace{},
		traceWithOneSpan,
		traceWithTwoSpans,
		traceWithTwoSpans,
	}

	// mock external dependencies
	pollerExecutor := getPollerExecutorWithMocks(t, retryDelay, maxWaitTimeForTrace, tracePerIteration)

	// When doing polling process
	// Then validate outputs
	executeAndValidatePollingRequests(t, pollerExecutor, []iterationExpectedValues{
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

func executeAndValidatePollingRequests(t *testing.T, pollerExecutor executor.PollerExecutor, expectedValues []iterationExpectedValues) {
	ctx := context.Background()
	run := model.NewRun()

	test := model.Test{
		ID: id.ID("some-test"),
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
		},
	}

	for i, value := range expectedValues {
		request := executor.NewPollingRequest(ctx, test, run, i)

		finished, anotherRun, err := pollerExecutor.ExecuteRequest(request)
		run = anotherRun // should store a run to use in another iteration

		require.NotNilf(t, run, "The test run should not be nil on iteration %d", i)

		if value.finished {
			require.Truef(t, finished, "The poller should have finished on iteration %d", i)
		} else {
			require.Falsef(t, finished, "The poller should have not finished on iteration %d", i)
		}

		if value.expectNoTraceError {
			require.Errorf(t, err, "An error should have happened on iteration %d", i)
			require.ErrorIsf(t, err, connection.ErrTraceNotFound, "An connection error should have happened on iteration %d", i)
		} else {
			require.NoErrorf(t, err, "An error should not have happened on iteration %d", i)

			// only validate root span if we have a root span
			if value.expectRootSpan {
				require.Truef(t, run.Trace.HasRootSpan(), "The trace associated with the run on iteration %d should have a root span", i)
			} else {
				require.Falsef(t, run.Trace.HasRootSpan(), "The trace associated with the run on iteration %d should not have a root span", i)
			}
		}
	}
}

func getPollerExecutorWithMocks(t *testing.T, retryDelay, maxWaitTimeForTrace time.Duration, tracePerIteration []model.Trace) executor.PollerExecutor {
	updater := getRunUpdaterMock(t)
	tracer := getTracerMock(t)
	testDB := getDataStoreRepositoryMock(t)
	traceDBFactory := getTraceDBMockFactory(t, tracePerIteration, &traceDBState{})

	return executor.NewPollerExecutor(
		retryDelay,
		maxWaitTimeForTrace,
		tracer,
		updater,
		traceDBFactory,
		testDB,
	)
}

// Mocks

// RunUpdater
type runUpdaterMock struct{}

func (m runUpdaterMock) Update(context.Context, model.Run) error { return nil }

func getRunUpdaterMock(t *testing.T) executor.RunUpdater {
	return runUpdaterMock{}
}

// DataStoreRepository
type dataStoreRepositoryMock struct {
	testdb.MockRepository
	// ...
}

func (m dataStoreRepositoryMock) DefaultDataStore(_ context.Context) (model.DataStore, error) {
	return model.DataStore{}, nil
}

func getDataStoreRepositoryMock(t *testing.T) model.Repository {
	t.Helper()

	mock := new(dataStoreRepositoryMock)
	mock.T = t
	mock.Test(t)

	return mock
}

// Tracer
func getTracerMock(t *testing.T) trace.Tracer {
	t.Helper()

	tracer, err := tracing.NewTracer(context.TODO(), config.Must(config.New(nil, log.Default())))
	require.NoError(t, err)

	return tracer
}

// TraceDB
type traceDBMock struct {
	tracePerIteration []model.Trace
	state             *traceDBState
}

func (db *traceDBMock) GetTraceByID(ctx context.Context, traceID string) (t model.Trace, err error) {
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
func (db *traceDBMock) TestConnection(ctx context.Context) connection.ConnectionTestResult {
	return connection.ConnectionTestResult{}
}

type traceDBState struct {
	currentIteration int
}

func getTraceDBMockFactory(t *testing.T, tracePerIteration []model.Trace, state *traceDBState) func(model.DataStore) (tracedb.TraceDB, error) {
	t.Helper()

	return func(ds model.DataStore) (tracedb.TraceDB, error) {
		return &traceDBMock{
			tracePerIteration: tracePerIteration,
			state:             state,
		}, nil
	}
}
