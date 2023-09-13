package tracepollerworker

import (
	"context"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
)

type SpanCountPollingStopStrategy struct{}

func NewSpanCountPollingStopStrategy() *SpanCountPollingStopStrategy {
	return &SpanCountPollingStopStrategy{}
}

// Evaluate implements PollingStopStrategy.
func (s *SpanCountPollingStopStrategy) Evaluate(ctx context.Context, job *executor.Job, traceDB tracedb.TraceDB, trace *traces.Trace) (bool, string) {
	if !traceDB.ShouldRetry() {
		return true, "TraceDB is not retryable"
	}

	maxTracePollRetry := job.PollingProfile.Periodic.MaxTracePollRetry()
	// we're done if we have the same amount of spans after polling or `maxTracePollRetry` times
	log.Printf("[TracePoller] Test %s Run %d: Job count %d, max retries: %d", job.Test.ID, job.Run.ID, job.EnqueueCount(), maxTracePollRetry)
	if job.EnqueueCount() >= maxTracePollRetry {
		return true, fmt.Sprintf("Hit MaxRetry of %d", maxTracePollRetry)
	}

	if job.Run.Trace == nil {
		return false, "First iteration"
	}

	collectedSpans := job.Headers.GetInt("collectedSpans")

	haveNotCollectedSpansSinceLastPoll := collectedSpans == 0
	haveCollectedSpansInTestRun := len(trace.Flat) > 0
	haveCollectedOnlyRootNode := len(trace.Flat) == 1 && trace.HasRootSpan()

	// Today we consider that we finished collecting traces
	// if we haven't collected any new spans since our last poll
	// and we have collected at least one span for this test run
	// and we have not collected only the root span

	if haveNotCollectedSpansSinceLastPoll && haveCollectedSpansInTestRun && !haveCollectedOnlyRootNode {
		return true, fmt.Sprintf("Trace has no new spans. Spans found: %d", len(trace.Flat))
	}

	return false, fmt.Sprintf("New spans found. Before: %d After: %d", len(job.Run.Trace.Flat), len(trace.Flat))
}

var _ PollingStopStrategy = &SpanCountPollingStopStrategy{}
