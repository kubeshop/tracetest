package executor_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type defaultPollerMock struct {
	mock.Mock
}

func (m *defaultPollerMock) ExecuteRequest(_ context.Context, job *executor.Job) (executor.PollResult, error) {
	args := m.Called(job)
	return args.Get(0).(executor.PollResult), args.Error(1)
}

type eventEmitterMock struct {
	mock.Mock
}

// Emit implements executor.EventEmitter
func (m *eventEmitterMock) Emit(ctx context.Context, event model.TestRunEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

var _ executor.EventEmitter = &eventEmitterMock{}

func TestSelectorBasedPollerExecutor(t *testing.T) {

	eventEmitter := new(eventEmitterMock)
	eventEmitter.On("Emit", mock.Anything, mock.Anything).Return(nil)

	createRequest := func(test test.Test, run test.Run) *executor.Job {
		pollingProfile := pollingprofile.DefaultPollingProfile
		pollingProfile.Periodic.SelectorMatchRetries = 3

		job := executor.NewJob()
		job.Test = test
		job.Run = run
		job.PollingProfile = pollingProfile

		return &job
	}

	t.Run("should return false when default poller returns false", func(t *testing.T) {
		defaultPoller := new(defaultPollerMock)
		selectorBasedPoller := executor.NewSelectorBasedPoller(defaultPoller, eventEmitter)

		request := createRequest(test.Test{}, test.Run{})

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(executor.PollResult{}, nil)
		res, _ := selectorBasedPoller.ExecuteRequest(context.TODO(), request)

		assert.False(t, res.Finished())
	})

	t.Run("should return false when default poller returns true but not all selector match at least one span", func(t *testing.T) {
		defaultPoller := new(defaultPollerMock)
		selectorBasedPoller := executor.NewSelectorBasedPoller(defaultPoller, eventEmitter)

		specs := test.Specs{
			{Selector: test.SpanQuery(`span[name = "Tracetest trigger"]`), Assertions: []test.Assertion{}},
			{Selector: test.SpanQuery(`span[name = "GET /api/tests"]`), Assertions: []test.Assertion{}},
		}
		testObj := test.Test{Specs: specs}

		trace := traces.NewTrace(randomIDGenerator.TraceID().String(), make([]traces.Span, 0))
		run := test.Run{Trace: &trace}

		request := createRequest(testObj, run)

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(executor.NewPollResult(true, "all spans found", run), nil)
		res, _ := selectorBasedPoller.ExecuteRequest(context.TODO(), request)

		assert.False(t, res.Finished())
		assert.Equal(t, 1, request.Headers.GetInt("SelectorBasedPollerExecutor.retryCount"))
	})

	t.Run("should return true if default poller returns true and selectors don't match spans 3 times in a row", func(t *testing.T) {
		defaultPoller := new(defaultPollerMock)
		selectorBasedPoller := executor.NewSelectorBasedPoller(defaultPoller, eventEmitter)

		specs := test.Specs{
			{Selector: test.SpanQuery(`span[name = "Tracetest trigger"]`), Assertions: []test.Assertion{}},
			{Selector: test.SpanQuery(`span[name = "GET /api/tests"]`), Assertions: []test.Assertion{}},
		}
		testObj := test.Test{Specs: specs}

		trace := traces.NewTrace(randomIDGenerator.TraceID().String(), make([]traces.Span, 0))
		run := test.Run{Trace: &trace}

		request := createRequest(testObj, run)

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(executor.NewPollResult(false, "trace not found", run), nil).Once()

		res, _ := selectorBasedPoller.ExecuteRequest(context.TODO(), request)
		assert.False(t, res.Finished())

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(executor.NewPollResult(true, "all spans found", run), nil)

		res, _ = selectorBasedPoller.ExecuteRequest(context.TODO(), request)
		assert.False(t, res.Finished())
		assert.Equal(t, 1, request.Headers.GetInt("SelectorBasedPollerExecutor.retryCount"))

		res, _ = selectorBasedPoller.ExecuteRequest(context.TODO(), request)
		assert.False(t, res.Finished())
		assert.Equal(t, 2, request.Headers.GetInt("SelectorBasedPollerExecutor.retryCount"))

		res, _ = selectorBasedPoller.ExecuteRequest(context.TODO(), request)
		assert.False(t, res.Finished())
		assert.Equal(t, 3, request.Headers.GetInt("SelectorBasedPollerExecutor.retryCount"))

		res, _ = selectorBasedPoller.ExecuteRequest(context.TODO(), request)
		assert.True(t, res.Finished())
	})

	t.Run("should return true if default poller returns true and each selector match at least one span", func(t *testing.T) {
		defaultPoller := new(defaultPollerMock)
		selectorBasedPoller := executor.NewSelectorBasedPoller(defaultPoller, eventEmitter)

		specs := test.Specs{
			{Selector: test.SpanQuery(`span[name = "Tracetest trigger"]`), Assertions: []test.Assertion{}},
			{Selector: test.SpanQuery(`span[name = "GET /api/tests"]`), Assertions: []test.Assertion{}},
		}
		testObj := test.Test{Specs: specs}

		rootSpan := traces.Span{ID: randomIDGenerator.SpanID(), Name: "Tracetest trigger", Attributes: make(traces.Attributes)}
		trace := traces.NewTrace(randomIDGenerator.TraceID().String(), []traces.Span{
			rootSpan,
			{ID: randomIDGenerator.SpanID(), Name: "GET /api/tests", Attributes: traces.Attributes{traces.TracetestMetadataFieldParentID: rootSpan.ID.String()}},
		})
		run := test.Run{Trace: &trace}

		request := createRequest(testObj, run)

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(executor.NewPollResult(true, "all spans found", run), nil)

		res, _ := selectorBasedPoller.ExecuteRequest(context.TODO(), request)
		assert.True(t, res.Finished())
	})
}
