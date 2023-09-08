package executor

import (
	"context"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DefaultPollerExecutor struct {
	updater      RunUpdater
	newTraceDBFn tracedb.FactoryFunc
	dsRepo       resourcemanager.Current[datastore.DataStore]
	eventEmitter EventEmitter
}

type InstrumentedPollerExecutor struct {
	tracer         trace.Tracer
	pollerExecutor pollerExecutor
}

func (pe InstrumentedPollerExecutor) ExecuteRequest(ctx context.Context, job *Job) (PollResult, error) {
	_, span := pe.tracer.Start(ctx, "Fetch trace")
	defer span.End()

	res, err := pe.pollerExecutor.ExecuteRequest(ctx, job)

	spanCount := 0
	if job.Run.Trace != nil {
		spanCount = len(job.Run.Trace.Flat)
	}

	attrs := []attribute.KeyValue{
		attribute.String("tracetest.run.trace_poller.trace_id", job.Run.TraceID.String()),
		attribute.String("tracetest.run.trace_poller.span_id", job.Run.SpanID.String()),
		attribute.Bool("tracetest.run.trace_poller.succesful", res.Finished()),
		attribute.String("tracetest.run.trace_poller.test_id", string(job.Test.ID)),
		attribute.Int("tracetest.run.trace_poller.amount_retrieved_spans", spanCount),
	}

	if res.reason != "" {
		attrs = append(attrs, attribute.String("tracetest.run.trace_poller.finish_reason", res.reason))
	}

	if err != nil {
		attrs = append(attrs, attribute.String("tracetest.run.trace_poller.error", err.Error()))
		span.RecordError(err)
	}

	span.SetAttributes(attrs...)
	return res, err
}

func NewPollerExecutor(
	tracer trace.Tracer,
	updater RunUpdater,
	newTraceDBFn tracedb.FactoryFunc,
	dsRepo resourcemanager.Current[datastore.DataStore],
	eventEmitter EventEmitter,
) *InstrumentedPollerExecutor {

	return &InstrumentedPollerExecutor{
		tracer: tracer,
		pollerExecutor: &DefaultPollerExecutor{
			updater:      updater,
			newTraceDBFn: newTraceDBFn,
			dsRepo:       dsRepo,
			eventEmitter: eventEmitter,
		},
	}
}

func (pe DefaultPollerExecutor) traceDB(ctx context.Context) (tracedb.TraceDB, error) {
	ds, err := pe.dsRepo.Current(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get default datastore: %w", err)
	}

	tdb, err := pe.newTraceDBFn(ds)
	if err != nil {
		return nil, fmt.Errorf(`cannot get tracedb from DataStore config with ID "%s": %w`, ds.ID, err)
	}

	return tdb, nil
}

func (pe DefaultPollerExecutor) ExecuteRequest(ctx context.Context, job *Job) (PollResult, error) {
	log.Printf("[PollerExecutor] Test %s Run %d: ExecuteRequest", job.Test.ID, job.Run.ID)

	traceDB, err := pe.traceDB(ctx)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: GetDataStore error: %s", job.Test.ID, job.Run.ID, err.Error())
		return PollResult{}, err
	}

	if isFirstRequest(job) {
		err := pe.testConnection(ctx, traceDB, job)
		if err != nil {
			return PollResult{}, err
		}
	}

	traceID := job.Run.TraceID.String()
	trace, err := traceDB.GetTraceByID(ctx, traceID)
	if err != nil {
		pe.emit(ctx, job, events.TracePollingIterationInfo(job.Test.ID, job.Run.ID, 0, job.EnqueueCount(), false, err.Error()))
		log.Printf("[PollerExecutor] Test %s Run %d: GetTraceByID (traceID %s) error: %s", job.Test.ID, job.Run.ID, traceID, err.Error())
		return PollResult{}, err
	}

	trace.ID = job.Run.TraceID
	done, reason := pe.donePollingTraces(job, traceDB, trace)
	// we need both values to be different to check for done, but after we want to have an updated job
	job.Run.Trace = &trace

	// we need to update at this point to persist the updated trace
	// otherwise we end up thinking every iteration is the first
	err = pe.updater.Update(ctx, job.Run)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: Update error: %s", job.Test.ID, job.Run.ID, err.Error())
		return PollResult{}, err
	}

	if !done {
		pe.emit(ctx, job, events.TracePollingIterationInfo(job.Test.ID, job.Run.ID, len(job.Run.Trace.Flat), job.EnqueueCount(), false, reason))
		log.Printf("[PollerExecutor] Test %s Run %d: Not done polling. (%s)", job.Test.ID, job.Run.ID, reason)

		return PollResult{
			finished: false,
			reason:   reason,
			run:      job.Run,
		}, nil
	}
	log.Printf("[PollerExecutor] Test %s Run %d: Done polling. (%s)", job.Test.ID, job.Run.ID, reason)

	log.Printf("[PollerExecutor] Test %s Run %d: Start Sorting", job.Test.ID, job.Run.ID)
	sorted := job.Run.Trace.Sort()
	job.Run.Trace = &sorted
	log.Printf("[PollerExecutor] Test %s Run %d: Sorting complete", job.Test.ID, job.Run.ID)

	if !job.Run.Trace.HasRootSpan() {
		newRoot := test.NewTracetestRootSpan(job.Run)
		job.Run.Trace = job.Run.Trace.InsertRootSpan(newRoot)
	} else {
		job.Run.Trace.RootSpan = traces.AugmentRootSpan(job.Run.Trace.RootSpan, job.Run.TriggerResult)
	}
	job.Run = job.Run.SuccessfullyPolledTraces(job.Run.Trace)

	log.Printf("[PollerExecutor] Completed polling process for Test Run %d after %d iterations, number of spans collected: %d ", job.Run.ID, job.EnqueueCount()+1, len(job.Run.Trace.Flat))

	log.Printf("[PollerExecutor] Test %s Run %d: Start updating", job.Test.ID, job.Run.ID)
	err = pe.updater.Update(ctx, job.Run)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: Update error: %s", job.Test.ID, job.Run.ID, err.Error())
		return PollResult{}, err
	}

	return PollResult{
		finished: true,
		reason:   reason,
		run:      job.Run,
	}, nil

}

func (pe DefaultPollerExecutor) emit(ctx context.Context, job *Job, event model.TestRunEvent) {
	err := pe.eventEmitter.Emit(ctx, event)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: failed to emit TracePollingIterationInfo event: error: %s", job.Test.ID, job.Run.ID, err.Error())
	}
}

func (pe DefaultPollerExecutor) testConnection(ctx context.Context, traceDB tracedb.TraceDB, job *Job) error {
	if testableTraceDB, ok := traceDB.(tracedb.TestableTraceDB); ok {
		connectionResult := testableTraceDB.TestConnection(ctx)

		err := pe.eventEmitter.Emit(ctx, events.TraceDataStoreConnectionInfo(job.Test.ID, job.Run.ID, connectionResult))
		if err != nil {
			log.Printf("[PollerExecutor] Test %s Run %d: failed to emit TraceDataStoreConnectionInfo event: error: %s", job.Test.ID, job.Run.ID, err.Error())
		}
	}

	endpoints := traceDB.GetEndpoints()
	ds, err := pe.dsRepo.Current(ctx)
	if err != nil {
		return fmt.Errorf("could not get current datastore: %w", err)
	}

	err = pe.eventEmitter.Emit(ctx, events.TracePollingStart(job.Test.ID, job.Run.ID, string(ds.Type), endpoints))
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: failed to emit TracePollingStart event: error: %s", job.Test.ID, job.Run.ID, err.Error())
	}

	return nil
}

func (pe DefaultPollerExecutor) donePollingTraces(job *Job, traceDB tracedb.TraceDB, trace traces.Trace) (bool, string) {
	if !traceDB.ShouldRetry() {
		return true, "TraceDB is not retryable"
	}

	maxTracePollRetry := job.PollingProfile.Periodic.MaxTracePollRetry()
	// we're done if we have the same amount of spans after polling or `maxTracePollRetry` times
	log.Printf("[PollerExecutor] Test %s Run %d: Job count %d, max retries: %d", job.Test.ID, job.Run.ID, job.EnqueueCount(), maxTracePollRetry)
	if job.EnqueueCount() >= maxTracePollRetry {
		return true, fmt.Sprintf("Hit MaxRetry of %d", maxTracePollRetry)
	}

	if job.Run.Trace == nil {
		return false, "First iteration"
	}

	haveNotCollectedSpansSinceLastPoll := len(trace.Flat) == len(job.Run.Trace.Flat)
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

func isFirstRequest(job *Job) bool {
	return !job.Headers.GetBool("requeued")
}
