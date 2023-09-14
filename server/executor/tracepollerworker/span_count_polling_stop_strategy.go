package tracepollerworker

import (
	"context"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/tracedb"
)

type SpanCountPollingStopStrategy struct{}

func NewSpanCountPollingStopStrategy() *SpanCountPollingStopStrategy {
	return &SpanCountPollingStopStrategy{}
}

// Evaluate implements PollingStopStrategy.
func (s *SpanCountPollingStopStrategy) Evaluate(ctx context.Context, job *executor.Job, traceDB tracedb.TraceDB) (bool, string) {
	if !traceDB.ShouldRetry() {
		return true, "TraceDB is not retryable"
	}

	maxTracePollRetry := job.PollingProfile.Periodic.MaxTracePollRetry()

	// we're done if we have the same amount of spans after polling or `maxTracePollRetry` times
	log.Printf("[TracePoller] Test %s Run %d: Job count %d, max retries: %d", job.Test.ID, job.Run.ID, job.EnqueueCount(), maxTracePollRetry)
	if job.EnqueueCount() >= maxTracePollRetry {
		return true, fmt.Sprintf("Hit MaxRetry of %d", maxTracePollRetry)
	}

	trace := job.Run.Trace

	collectedSpans := job.Headers.GetInt("collectedSpans")

	haveCollectedSpansInThisIteration := collectedSpans > 0
	if haveCollectedSpansInThisIteration { // found spans on this iteration, continue the iterations
		return false, fmt.Sprintf("Found %d new spans.", collectedSpans)
	}

	totalSpans := len(trace.Flat)
	haveCollectedOnlyRootNode := totalSpans == 1 && trace.HasRootSpan()

	if haveCollectedOnlyRootNode { // found only "Tracetest trigger" span, continue the iterations
		return false, "Found only Tracetest trigger span."
	}

	haveCollectedSpansInPreviousIterations := len(trace.Flat) > 0

	if haveCollectedSpansInPreviousIterations && !haveCollectedSpansInThisIteration {
		// found not found spans on this iteration, but it found it on previous iterations
		return true, fmt.Sprintf("Trace has no new spans. Spans found: %d", len(trace.Flat))
	}

	return false, "" // given the previous conditions, this code should not be reached
}

var _ PollingStopStrategy = &SpanCountPollingStopStrategy{}
