package executor

import (
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/test"
)

const (
	selectorBasedPollerExecutorRetryHeader = "SelectorBasedPollerExecutor.retryCount"
	selectorBasedPollerExecutorMaxTries    = 3
)

type selectorBasedPollerExecutor struct {
	pollerExecutor PollerExecutor
	eventEmitter   EventEmitter
}

func NewSelectorBasedPoller(innerPoller PollerExecutor, eventEmitter EventEmitter) PollerExecutor {
	return selectorBasedPollerExecutor{innerPoller, eventEmitter}
}

func (pe selectorBasedPollerExecutor) ExecuteRequest(request *PollingRequest) (bool, string, test.Run, error) {
	ready, reason, run, err := pe.pollerExecutor.ExecuteRequest(request)
	if !ready {
		request.SetHeaderInt(selectorBasedPollerExecutorRetryHeader, 0)
		return ready, reason, run, err
	}

	maxNumberRetries := 0
	if request.pollingProfile.Periodic != nil {
		maxNumberRetries = request.pollingProfile.Periodic.SelectorMatchRetries
	}

	currentNumberTries := pe.getNumberTries(request)
	if currentNumberTries >= maxNumberRetries {
		pe.eventEmitter.Emit(request.Context(), events.TracePollingIterationInfo(
			request.test.ID,
			request.run.ID,
			len(request.run.Trace.Flat),
			request.count,
			true,
			fmt.Sprintf("Some selectors did not match any spans in the current trace, but after %d tries, the trace probably won't change", currentNumberTries),
		))
		return true, "", run, err
	}

	allSelectorsMatchSpans := pe.allSelectorsMatchSpans(request)
	if allSelectorsMatchSpans {
		pe.eventEmitter.Emit(request.Context(), events.TracePollingIterationInfo(
			request.test.ID,
			request.run.ID,
			len(request.run.Trace.Flat),
			request.count,
			true,
			"All selectors from the test matched at least one span in the current trace",
		))
		return true, "", run, err
	}

	request.SetHeaderInt(selectorBasedPollerExecutorRetryHeader, currentNumberTries+1)

	pe.eventEmitter.Emit(request.Context(), events.TracePollingIterationInfo(
		request.test.ID,
		request.run.ID,
		len(request.run.Trace.Flat),
		request.count,
		false,
		"All selectors from your test must match at least one span in the trace, some of them did not match any",
	))

	return false, "not all selectors got matching spans in the trace", run, err
}

func (pe selectorBasedPollerExecutor) getNumberTries(request *PollingRequest) int {
	value := request.Header(selectorBasedPollerExecutorRetryHeader)
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}

	return 0
}

func (pe selectorBasedPollerExecutor) allSelectorsMatchSpans(request *PollingRequest) bool {
	allSelectorsHaveMatch := true
	for _, spec := range request.test.Specs {
		spans := selector(spec.Selector.Query).Filter(*request.run.Trace)
		if len(spans) == 0 {
			allSelectorsHaveMatch = false
		}
	}

	return allSelectorsHaveMatch
}
