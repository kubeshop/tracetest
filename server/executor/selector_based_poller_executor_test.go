package executor_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type defaultPollerMock struct {
	mock.Mock
}

// ExecuteRequest implements executor.PollerExecutor
func (m *defaultPollerMock) ExecuteRequest(ctx context.Context, request *executor.PollingRequest) (bool, string, test.Run, error) {
	args := m.Called(request)
	return args.Bool(0), args.String(1), args.Get(2).(test.Run), args.Error(3)
}

var _ executor.PollerExecutor = &defaultPollerMock{}

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

	createRequest := func(test test.Test, run test.Run) *executor.PollingRequest {
		pollingProfile := pollingprofile.DefaultPollingProfile
		pollingProfile.Periodic.SelectorMatchRetries = 3

		request := executor.NewPollingRequest(context.Background(), test, run, 0, pollingProfile)
		return request
	}

	t.Run("should return false when default poller returns false", func(t *testing.T) {
		defaultPoller := new(defaultPollerMock)
		selectorBasedPoller := executor.NewSelectorBasedPoller(defaultPoller, eventEmitter)

		request := createRequest(test.Test{}, test.Run{})

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(false, "", test.Run{}, nil)
		ready, _, _, _ := selectorBasedPoller.ExecuteRequest(context.Background(), request)

		assert.False(t, ready)
	})

	t.Run("should return false when default poller returns true but not all selector match at least one span", func(t *testing.T) {
		defaultPoller := new(defaultPollerMock)
		selectorBasedPoller := executor.NewSelectorBasedPoller(defaultPoller, eventEmitter)

		specs := test.Specs{
			{Selector: test.SpanQuery(`span[name = "Tracetest trigger"]`), Assertions: []test.Assertion{}},
			{Selector: test.SpanQuery(`span[name = "GET /api/tests"]`), Assertions: []test.Assertion{}},
		}
		testObj := test.Test{Specs: specs}

		trace := model.NewTrace(randomIDGenerator.TraceID().String(), make([]model.Span, 0))
		run := test.Run{Trace: &trace}

		request := createRequest(testObj, run)

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(true, "all spans found", run, nil)
		ready, _, _, _ := selectorBasedPoller.ExecuteRequest(context.Background(), request)

		assert.False(t, ready)
		assert.Equal(t, "1", request.Header("SelectorBasedPollerExecutor.retryCount"))
	})

	t.Run("should return true if default poller returns true and selectors don't match spans 3 times in a row", func(t *testing.T) {
		defaultPoller := new(defaultPollerMock)
		selectorBasedPoller := executor.NewSelectorBasedPoller(defaultPoller, eventEmitter)

		specs := test.Specs{
			{Selector: test.SpanQuery(`span[name = "Tracetest trigger"]`), Assertions: []test.Assertion{}},
			{Selector: test.SpanQuery(`span[name = "GET /api/tests"]`), Assertions: []test.Assertion{}},
		}
		testObj := test.Test{Specs: specs}

		trace := model.NewTrace(randomIDGenerator.TraceID().String(), make([]model.Span, 0))
		run := test.Run{Trace: &trace}

		request := createRequest(testObj, run)

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(false, "trace not found", run, nil).Once()

		ready, _, _, _ := selectorBasedPoller.ExecuteRequest(context.Background(), request)
		assert.False(t, ready)

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(true, "all spans found", run, nil)

		ready, _, _, _ = selectorBasedPoller.ExecuteRequest(context.Background(), request)
		assert.False(t, ready)
		assert.Equal(t, "1", request.Header("SelectorBasedPollerExecutor.retryCount"))

		ready, _, _, _ = selectorBasedPoller.ExecuteRequest(context.Background(), request)
		assert.False(t, ready)
		assert.Equal(t, "2", request.Header("SelectorBasedPollerExecutor.retryCount"))

		ready, _, _, _ = selectorBasedPoller.ExecuteRequest(context.Background(), request)
		assert.False(t, ready)
		assert.Equal(t, "3", request.Header("SelectorBasedPollerExecutor.retryCount"))

		ready, _, _, _ = selectorBasedPoller.ExecuteRequest(context.Background(), request)
		assert.True(t, ready)
	})

	t.Run("should return true if default poller returns true and each selector match at least one span", func(t *testing.T) {
		defaultPoller := new(defaultPollerMock)
		selectorBasedPoller := executor.NewSelectorBasedPoller(defaultPoller, eventEmitter)

		specs := test.Specs{
			{Selector: test.SpanQuery(`span[name = "Tracetest trigger"]`), Assertions: []test.Assertion{}},
			{Selector: test.SpanQuery(`span[name = "GET /api/tests"]`), Assertions: []test.Assertion{}},
		}
		testObj := test.Test{Specs: specs}

		rootSpan := model.Span{ID: randomIDGenerator.SpanID(), Name: "Tracetest trigger", Attributes: make(model.Attributes)}
		trace := model.NewTrace(randomIDGenerator.TraceID().String(), []model.Span{
			rootSpan,
			{ID: randomIDGenerator.SpanID(), Name: "GET /api/tests", Attributes: model.Attributes{model.TracetestMetadataFieldParentID: rootSpan.ID.String()}},
		})
		run := test.Run{Trace: &trace}

		request := createRequest(testObj, run)

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(true, "all spans found", run, nil)

		ready, _, _, _ := selectorBasedPoller.ExecuteRequest(context.Background(), request)
		assert.True(t, ready)
	})
}
