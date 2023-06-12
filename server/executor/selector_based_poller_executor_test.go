package executor_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type defaultPollerMock struct {
	mock.Mock
}

// ExecuteRequest implements executor.PollerExecutor
func (m *defaultPollerMock) ExecuteRequest(request *executor.PollingRequest) (bool, string, model.Run, error) {
	args := m.Called(request)
	return args.Bool(0), args.String(1), args.Get(2).(model.Run), args.Error(3)
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

	createRequest := func(test model.Test, run model.Run) *executor.PollingRequest {
		pollingProfile := pollingprofile.DefaultPollingProfile
		pollingProfile.Periodic.SelectorMatchRetries = 3

		request := executor.NewPollingRequest(context.Background(), test, run, 0, pollingProfile)
		return request
	}

	t.Run("should return false when default poller returns false", func(t *testing.T) {
		defaultPoller := new(defaultPollerMock)
		selectorBasedPoller := executor.NewSelectorBasedPoller(defaultPoller, eventEmitter)

		request := createRequest(model.Test{}, model.Run{})

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(false, "", model.Run{}, nil)
		ready, _, _, _ := selectorBasedPoller.ExecuteRequest(request)

		assert.False(t, ready)
	})

	t.Run("should return false when default poller returns true but not all selector match at least one span", func(t *testing.T) {
		defaultPoller := new(defaultPollerMock)
		selectorBasedPoller := executor.NewSelectorBasedPoller(defaultPoller, eventEmitter)

		specs := maps.Ordered[model.SpanQuery, model.NamedAssertions]{}.
			MustAdd(`span[name = "Tracetest trigger"]`, model.NamedAssertions{}).
			MustAdd(`span[name = "GET /api/tests"]`, model.NamedAssertions{})
		test := model.Test{Specs: specs}

		trace := model.NewTrace(randomIDGenerator.TraceID().String(), make([]model.Span, 0))
		run := model.Run{Trace: &trace}

		request := createRequest(test, run)

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(true, "all spans found", run, nil)
		ready, _, _, _ := selectorBasedPoller.ExecuteRequest(request)

		assert.False(t, ready)
		assert.Equal(t, "1", request.Header("SelectorBasedPollerExecutor.retryCount"))
	})

	t.Run("should return true if default poller returns true and selectors don't match spans 3 times in a row", func(t *testing.T) {
		defaultPoller := new(defaultPollerMock)
		selectorBasedPoller := executor.NewSelectorBasedPoller(defaultPoller, eventEmitter)

		specs := maps.Ordered[model.SpanQuery, model.NamedAssertions]{}.
			MustAdd(`span[name = "Tracetest trigger"]`, model.NamedAssertions{}).
			MustAdd(`span[name = "GET /api/tests"]`, model.NamedAssertions{})
		test := model.Test{Specs: specs}

		trace := model.NewTrace(randomIDGenerator.TraceID().String(), make([]model.Span, 0))
		run := model.Run{Trace: &trace}

		request := createRequest(test, run)

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(false, "trace not found", run, nil).Once()

		ready, _, _, _ := selectorBasedPoller.ExecuteRequest(request)
		assert.False(t, ready)

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(true, "all spans found", run, nil)

		ready, _, _, _ = selectorBasedPoller.ExecuteRequest(request)
		assert.False(t, ready)
		assert.Equal(t, "1", request.Header("SelectorBasedPollerExecutor.retryCount"))

		ready, _, _, _ = selectorBasedPoller.ExecuteRequest(request)
		assert.False(t, ready)
		assert.Equal(t, "2", request.Header("SelectorBasedPollerExecutor.retryCount"))

		ready, _, _, _ = selectorBasedPoller.ExecuteRequest(request)
		assert.False(t, ready)
		assert.Equal(t, "3", request.Header("SelectorBasedPollerExecutor.retryCount"))

		ready, _, _, _ = selectorBasedPoller.ExecuteRequest(request)
		assert.True(t, ready)
	})

	t.Run("should return true if default poller returns true and each selector match at least one span", func(t *testing.T) {
		defaultPoller := new(defaultPollerMock)
		selectorBasedPoller := executor.NewSelectorBasedPoller(defaultPoller, eventEmitter)

		specs := maps.Ordered[model.SpanQuery, model.NamedAssertions]{}.
			MustAdd(`span[name = "Tracetest trigger"]`, model.NamedAssertions{}).
			MustAdd(`span[name = "GET /api/tests"]`, model.NamedAssertions{})
		test := model.Test{Specs: specs}

		rootSpan := model.Span{ID: randomIDGenerator.SpanID(), Name: "Tracetest trigger", Attributes: make(model.Attributes)}
		trace := model.NewTrace(randomIDGenerator.TraceID().String(), []model.Span{
			rootSpan,
			{ID: randomIDGenerator.SpanID(), Name: "GET /api/tests", Attributes: model.Attributes{string(model.TracetestMetadataFieldParentID): rootSpan.ID.String()}},
		})
		run := model.Run{Trace: &trace}

		request := createRequest(test, run)

		defaultPoller.On("ExecuteRequest", mock.Anything).Return(true, "all spans found", run, nil)

		ready, _, _, _ := selectorBasedPoller.ExecuteRequest(request)
		assert.True(t, ready)
	})
}
