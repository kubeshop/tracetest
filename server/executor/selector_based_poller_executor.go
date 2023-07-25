package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/model/events"
)

const (
	selectorBasedPollerExecutorRetryHeader = "SelectorBasedPollerExecutor.retryCount"
	selectorBasedPollerExecutorMaxTries    = 3
)

type selectorBasedPollerExecutor struct {
	pollerExecutor pollerExecutor
	eventEmitter   EventEmitter
}

func NewSelectorBasedPoller(innerPoller pollerExecutor, eventEmitter EventEmitter) selectorBasedPollerExecutor {
	return selectorBasedPollerExecutor{innerPoller, eventEmitter}
}

func (pe selectorBasedPollerExecutor) ExecuteRequest(ctx context.Context, job *Job) (PollResult, error) {
	res, err := pe.pollerExecutor.ExecuteRequest(ctx, job)
	if !res.finished {
		job.Headers.SetInt(selectorBasedPollerExecutorRetryHeader, 0)
		return res, err
	}

	maxNumberRetries := 0
	if job.PollingProfile.Periodic != nil {
		maxNumberRetries = job.PollingProfile.Periodic.SelectorMatchRetries
	}

	currentNumberTries := job.Headers.GetInt(selectorBasedPollerExecutorRetryHeader)
	if currentNumberTries >= maxNumberRetries {
		pe.eventEmitter.Emit(ctx, events.TracePollingIterationInfo(
			job.Test.ID,
			job.Run.ID,
			len(job.Run.Trace.Flat),
			currentNumberTries,
			true,
			fmt.Sprintf("Some selectors did not match any spans in the current trace, but after %d tries, the trace probably won't change", currentNumberTries),
		))
		res.finished = true
		return res, err
	}

	allSelectorsMatchSpans := pe.allSelectorsMatchSpans(job)
	if allSelectorsMatchSpans {
		pe.eventEmitter.Emit(ctx, events.TracePollingIterationInfo(
			job.Test.ID,
			job.Run.ID,
			len(job.Run.Trace.Flat),
			currentNumberTries,
			true,
			"All selectors from the test matched at least one span in the current trace",
		))
		res.finished = true
		return res, err
	}

	job.Headers.SetInt(selectorBasedPollerExecutorRetryHeader, currentNumberTries+1)

	pe.eventEmitter.Emit(ctx, events.TracePollingIterationInfo(
		job.Test.ID,
		job.Run.ID,
		len(job.Run.Trace.Flat),
		job.Headers.GetInt(selectorBasedPollerExecutorRetryHeader),
		false,
		"All selectors from your test must match at least one span in the trace, some of them did not match any",
	))

	res.finished = false
	res.reason = "not all selectors got matching spans in the trace"

	return res, err
}

func (pe selectorBasedPollerExecutor) allSelectorsMatchSpans(job *Job) bool {
	allSelectorsHaveMatch := true
	for _, spec := range job.Test.Specs {
		spans := selector(spec.Selector).Filter(*job.Run.Trace)
		if len(spans) == 0 {
			allSelectorsHaveMatch = false
		}
	}

	return allSelectorsHaveMatch
}
