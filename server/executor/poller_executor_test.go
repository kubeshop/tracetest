package executor_test

import (
	"context"
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

	// Scenario: Trace without spans
	// Given the trigger execution returns 0 spans
	// And tracetest does not send the root span
	// When the server do the polling process
	// Then it should stop at the second iteration
	// And it should have no error on the process
	// And a root span should be added to it

	// Given conditions

	// maxRetries=30 (inferred by the calculation: maxWaitTimeForTrace / retryDelay)
	retryDelay := 1 * time.Second
	maxWaitTimeForTrace := 30 * time.Second

	tracePerIteration := map[int]model.Trace{
		0: model.Trace{},
		1: model.Trace{},
	}

	// mock external dependencies
	updater := getRunUpdaterMock(t)
	tracer := getTracerMock(t)
	testDB := getDataStoreRepositoryMock(t)
	traceDBFactory := getTraceDBMockFactory(t, tracePerIteration, &traceDBState{})

	pollerExecutor := executor.NewPollerExecutor(
		retryDelay,
		maxWaitTimeForTrace,
		tracer,
		updater,
		traceDBFactory,
		testDB,
	)

	ctx := context.Background()
	test := model.Test{
		ID: id.ID("some-test"),
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
		},
	}

	// When doing polling process
	// Then validate outputs

	var finished bool
	var run model.Run
	var err error

	// first iteration
	firstRun := model.NewRun()
	requestForFirstIteration := executor.NewPollingRequest(ctx, test, firstRun, 0)

	finished, run, err = pollerExecutor.ExecuteRequest(requestForFirstIteration)
	require.False(t, finished)
	require.NoError(t, err)
	require.NotNil(t, run)
	require.False(t, run.Trace.HasRootSpan())

	// second iteration
	secondRun := run
	requestForSecondIteration := executor.NewPollingRequest(ctx, test, secondRun, 1)

	finished, run, err = pollerExecutor.ExecuteRequest(requestForSecondIteration)
	require.True(t, finished)
	require.NoError(t, err)
	require.NotNil(t, run)
	require.True(t, run.Trace.HasRootSpan())
}

func Test_PollerExecutor_ExecuteRequest_NoRootSpan_OneSpanCase(t *testing.T) {
	t.Parallel()

	// Scenario: Trace with only 1 span
	// Given the trigger execution returns 1 span on first iteration
	// And tracetest does not send the root span
	// When the server do the polling process
	// Then it should stop at the second iteration
	// And it should have no error on the process
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
	tracePerIteration := map[int]model.Trace{
		0: trace,
		1: trace,
	}

	// mock external dependencies
	updater := getRunUpdaterMock(t)
	tracer := getTracerMock(t)
	testDB := getDataStoreRepositoryMock(t)
	traceDBFactory := getTraceDBMockFactory(t, tracePerIteration, &traceDBState{})

	pollerExecutor := executor.NewPollerExecutor(
		retryDelay,
		maxWaitTimeForTrace,
		tracer,
		updater,
		traceDBFactory,
		testDB,
	)

	ctx := context.Background()
	test := model.Test{
		ID: id.ID("some-test"),
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
		},
	}

	// When doing polling process
	// Then validate outputs

	var finished bool
	var run model.Run
	var err error

	// first iteration
	firstRun := model.NewRun()
	requestForFirstIteration := executor.NewPollingRequest(ctx, test, firstRun, 0)

	finished, run, err = pollerExecutor.ExecuteRequest(requestForFirstIteration)
	require.False(t, finished)
	require.NoError(t, err)
	require.NotNil(t, run)
	require.False(t, run.Trace.HasRootSpan())

	// second iteration
	secondRun := run
	requestForSecondIteration := executor.NewPollingRequest(ctx, test, secondRun, 1)

	finished, run, err = pollerExecutor.ExecuteRequest(requestForSecondIteration)
	require.True(t, finished)
	require.NoError(t, err)
	require.NotNil(t, run)
	require.True(t, run.Trace.HasRootSpan())
}

func Test_PollerExecutor_ExecuteRequest_NoRootSpan_TwoSpansCase(t *testing.T) {
	t.Parallel()

	// Scenario: Trace with 2 span
	// Given the trigger execution returns 1 span on first iteration and another one on second iteration
	// And tracetest does not send the root span
	// When the server do the polling process
	// Then it should stop at the third iteration
	// And it should have no error on the process
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

	firstIterationTrace := model.NewTrace(traceID, []model.Span{firstSpan})
	secondIterationTrace := model.NewTrace(traceID, []model.Span{firstSpan, secondSpan})
	thirdIterationTrace := model.NewTrace(traceID, []model.Span{firstSpan, secondSpan})

	// test
	tracePerIteration := map[int]model.Trace{
		0: firstIterationTrace,
		1: secondIterationTrace,
		2: thirdIterationTrace,
	}

	// mock external dependencies
	updater := getRunUpdaterMock(t)
	tracer := getTracerMock(t)
	testDB := getDataStoreRepositoryMock(t)
	traceDBFactory := getTraceDBMockFactory(t, tracePerIteration, &traceDBState{})

	pollerExecutor := executor.NewPollerExecutor(
		retryDelay,
		maxWaitTimeForTrace,
		tracer,
		updater,
		traceDBFactory,
		testDB,
	)

	ctx := context.Background()
	test := model.Test{
		ID: id.ID("some-test"),
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
		},
	}

	// When doing polling process
	// Then validate outputs

	var finished bool
	var run model.Run
	var err error

	// first iteration
	firstRun := model.NewRun()
	requestForFirstIteration := executor.NewPollingRequest(ctx, test, firstRun, 0)

	finished, run, err = pollerExecutor.ExecuteRequest(requestForFirstIteration)
	require.False(t, finished)
	require.NoError(t, err)
	require.NotNil(t, run)
	require.False(t, run.Trace.HasRootSpan())

	// second iteration
	secondRun := run
	requestForSecondIteration := executor.NewPollingRequest(ctx, test, secondRun, 1)

	finished, run, err = pollerExecutor.ExecuteRequest(requestForSecondIteration)
	require.False(t, finished)
	require.NoError(t, err)
	require.NotNil(t, run)
	require.False(t, run.Trace.HasRootSpan())

	// third iteration
	thirdRun := run
	requestForThirdIteration := executor.NewPollingRequest(ctx, test, thirdRun, 1)

	finished, run, err = pollerExecutor.ExecuteRequest(requestForThirdIteration)
	require.True(t, finished)
	require.NoError(t, err)
	require.NotNil(t, run)
	require.True(t, run.Trace.HasRootSpan())
}

// TODO: add cases where Tracetest root came along

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

	tracer, err := tracing.NewTracer(context.TODO(), config.New())
	require.NoError(t, err)

	return tracer
}

// TraceDB
type traceDBMock struct {
	tracePerIteration map[int]model.Trace
	state             *traceDBState
}

func (db *traceDBMock) GetTraceByID(ctx context.Context, traceID string) (t model.Trace, err error) {
	trace := db.tracePerIteration[db.state.currentIteration]
	db.state.currentIteration += 1

	return trace, nil
}

func (db *traceDBMock) ShouldRetry() bool {
	return true // this provider should retry
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

func getTraceDBMockFactory(t *testing.T, tracePerIteration map[int]model.Trace, state *traceDBState) func(model.DataStore) (tracedb.TraceDB, error) {
	t.Helper()

	return func(ds model.DataStore) (tracedb.TraceDB, error) {
		return &traceDBMock{
			tracePerIteration: tracePerIteration,
			state:             state,
		}, nil
	}
}
