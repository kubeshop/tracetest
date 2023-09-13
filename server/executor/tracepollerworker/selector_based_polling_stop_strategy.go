package tracepollerworker

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
)

type SelectorBasedPollingStopStrategy struct {
	eventEmitter    executor.EventEmitter
	wrappedStrategy PollingStopStrategy
}

const (
	selectorBasedPollerExecutorRetryHeader = "SelectorBasedPollerExecutor.retryCount"
)

func NewSelectorBasedPollingStopStrategy(eventEmitter executor.EventEmitter, strategy PollingStopStrategy) *SelectorBasedPollingStopStrategy {
	return &SelectorBasedPollingStopStrategy{
		eventEmitter:    eventEmitter,
		wrappedStrategy: strategy,
	}
}

// Evaluate implements PollingStopStrategy.
func (s *SelectorBasedPollingStopStrategy) Evaluate(ctx context.Context, job *executor.Job, traceDB tracedb.TraceDB, trace *traces.Trace) (bool, string) {
	if !traceDB.ShouldRetry() {
		return true, "TraceDB is not retryable"
	}

	finished, reason := s.wrappedStrategy.Evaluate(ctx, job, traceDB, trace)

	if !finished {
		job.Headers.SetInt(selectorBasedPollerExecutorRetryHeader, 0)
		return finished, reason
	}

	maxNumberRetries := 0
	if job.PollingProfile.Periodic != nil {
		maxNumberRetries = job.PollingProfile.Periodic.SelectorMatchRetries
	}

	currentNumberTries := job.Headers.GetInt(selectorBasedPollerExecutorRetryHeader)
	if currentNumberTries >= maxNumberRetries {

		s.eventEmitter.Emit(ctx, events.TracePollingIterationInfo(
			job.Test.ID,
			job.Run.ID,
			len(job.Run.Trace.Flat),
			currentNumberTries,
			true,
			fmt.Sprintf("Some selectors did not match any spans in the current trace, but after %d tries, the trace probably won't change", currentNumberTries),
		))

		return true, reason
	}

	allSelectorsMatchSpans := s.allSelectorsMatchSpans(job)
	if allSelectorsMatchSpans {
		s.eventEmitter.Emit(ctx, events.TracePollingIterationInfo(
			job.Test.ID,
			job.Run.ID,
			len(job.Run.Trace.Flat),
			currentNumberTries,
			true,
			"All selectors from the test matched at least one span in the current trace",
		))
		return true, reason
	}

	job.Headers.SetInt(selectorBasedPollerExecutorRetryHeader, currentNumberTries+1)

	s.eventEmitter.Emit(ctx, events.TracePollingIterationInfo(
		job.Test.ID,
		job.Run.ID,
		len(job.Run.Trace.Flat),
		job.Headers.GetInt(selectorBasedPollerExecutorRetryHeader),
		false,
		"All selectors from your test must match at least one span in the trace, some of them did not match any",
	))

	return false, "not all selectors got matching spans in the trace"
}

func (pe SelectorBasedPollingStopStrategy) allSelectorsMatchSpans(job *executor.Job) bool {
	allSelectorsHaveMatch := true
	for _, spec := range job.Test.Specs {
		spans := selector(spec.Selector).Filter(*job.Run.Trace)
		if len(spans) == 0 {
			allSelectorsHaveMatch = false
		}
	}

	return allSelectorsHaveMatch
}

func selector(sq test.SpanQuery) selectors.Selector {
	sel, _ := selectors.New(string(sq))
	return sel
}

var _ PollingStopStrategy = &SelectorBasedPollingStopStrategy{}
