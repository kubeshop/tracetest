package executor

import (
	"strconv"

	"github.com/kubeshop/tracetest/server/model"
)

const (
	selectorBasedPollerExecutorRetryHeader = "SelectorBasedPollerExecutor.retryCount"
	selectorBasedPollerExecutorMaxTries    = 3
)

type selectorBasedPollerExecutor struct {
	pollerExecutor PollerExecutor
}

func NewSelectorBasedPoller(innerPoller PollerExecutor) PollerExecutor {
	return selectorBasedPollerExecutor{innerPoller}
}

func (pe selectorBasedPollerExecutor) ExecuteRequest(request *PollingRequest) (bool, string, model.Run, error) {
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
		return true, "not all selectors matched, but trace haven't changed in a while", run, err
	}

	allSelectorsMatchSpans := pe.allSelectorsMatchSpans(request)
	if allSelectorsMatchSpans {
		return true, "all selectors have matched one or more spans", run, err
	}

	request.SetHeaderInt(selectorBasedPollerExecutorRetryHeader, currentNumberTries+1)
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
	request.test.Specs.ForEach(func(selectorQuery model.SpanQuery, _ model.NamedAssertions) error {
		spans := selector(selectorQuery).Filter(*request.run.Trace)
		if len(spans) == 0 {
			allSelectorsHaveMatch = false
		}

		return nil
	})

	return allSelectorsHaveMatch
}
