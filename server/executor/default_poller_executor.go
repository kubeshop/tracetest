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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type traceDBFactoryFn func(ds datastore.DataStore) (tracedb.TraceDB, error)

type DefaultPollerExecutor struct {
	ppGetter     PollingProfileGetter
	updater      RunUpdater
	newTraceDBFn traceDBFactoryFn
	dsRepo       resourcemanager.Current[datastore.DataStore]
	eventEmitter EventEmitter
}

type InstrumentedPollerExecutor struct {
	tracer         trace.Tracer
	pollerExecutor PollerExecutor
}

func (pe InstrumentedPollerExecutor) ExecuteRequest(request *PollingRequest) (bool, string, test.Run, error) {
	_, span := pe.tracer.Start(request.Context(), "Fetch trace")
	defer span.End()

	finished, finishReason, run, err := pe.pollerExecutor.ExecuteRequest(request)

	spanCount := 0
	if run.Trace != nil {
		spanCount = len(run.Trace.Flat)
	}

	attrs := []attribute.KeyValue{
		attribute.String("tracetest.run.trace_poller.trace_id", request.run.TraceID.String()),
		attribute.String("tracetest.run.trace_poller.span_id", request.run.SpanID.String()),
		attribute.Bool("tracetest.run.trace_poller.succesful", finished),
		attribute.String("tracetest.run.trace_poller.test_id", string(request.test.ID)),
		attribute.Int("tracetest.run.trace_poller.amount_retrieved_spans", spanCount),
	}

	if finishReason != "" {
		attrs = append(attrs, attribute.String("tracetest.run.trace_poller.finish_reason", finishReason))
	}

	if err != nil {
		attrs = append(attrs, attribute.String("tracetest.run.trace_poller.error", err.Error()))
		span.RecordError(err)
	}

	span.SetAttributes(attrs...)
	return finished, finishReason, run, err
}

func NewPollerExecutor(
	ppGetter PollingProfileGetter,
	tracer trace.Tracer,
	updater RunUpdater,
	newTraceDBFn traceDBFactoryFn,
	dsRepo resourcemanager.Current[datastore.DataStore],
	eventEmitter EventEmitter,
) PollerExecutor {

	defaultExecutor := &DefaultPollerExecutor{
		ppGetter:     ppGetter,
		updater:      updater,
		newTraceDBFn: newTraceDBFn,
		dsRepo:       dsRepo,
		eventEmitter: eventEmitter,
	}

	return &InstrumentedPollerExecutor{
		tracer:         tracer,
		pollerExecutor: defaultExecutor,
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

func (pe DefaultPollerExecutor) ExecuteRequest(request *PollingRequest) (bool, string, test.Run, error) {
	log.Printf("[PollerExecutor] Test %s Run %d: ExecuteRequest\n", request.test.ID, request.run.ID)
	run := request.run
	ctx := request.Context()

	traceDB, err := pe.traceDB(ctx)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: GetDataStore error: %s\n", request.test.ID, request.run.ID, err.Error())
		return false, "", test.Run{}, err
	}

	if request.IsFirstRequest() {
		if testableTraceDB, ok := traceDB.(tracedb.TestableTraceDB); ok {
			connectionResult := testableTraceDB.TestConnection(ctx)

			err = pe.eventEmitter.Emit(ctx, events.TraceDataStoreConnectionInfo(request.test.ID, request.run.ID, connectionResult))
			if err != nil {
				log.Printf("[PollerExecutor] Test %s Run %d: failed to emit TraceDataStoreConnectionInfo event: error: %s\n", request.test.ID, request.run.ID, err.Error())
			}
		}

		endpoints := traceDB.GetEndpoints()
		ds, err := pe.dsRepo.Current(ctx)
		if err != nil {
			return false, "", test.Run{}, fmt.Errorf("could not get current datastore: %w", err)
		}

		err = pe.eventEmitter.Emit(ctx, events.TracePollingStart(request.test.ID, request.run.ID, string(ds.Type), endpoints))
		if err != nil {
			log.Printf("[PollerExecutor] Test %s Run %d: failed to emit TracePollingStart event: error: %s\n", request.test.ID, request.run.ID, err.Error())
		}
	}

	traceID := run.TraceID.String()
	trace, err := traceDB.GetTraceByID(ctx, traceID)
	if err != nil {
		anotherErr := pe.eventEmitter.Emit(ctx, events.TracePollingIterationInfo(request.test.ID, request.run.ID, 0, request.count, false, err.Error()))
		if anotherErr != nil {
			log.Printf("[PollerExecutor] Test %s Run %d: failed to emit TracePollingIterationInfo event: error: %s\n", request.test.ID, request.run.ID, anotherErr.Error())
		}

		log.Printf("[PollerExecutor] Test %s Run %d: GetTraceByID (traceID %s) error: %s\n", request.test.ID, request.run.ID, traceID, err.Error())
		return false, "", test.Run{}, err
	}

	trace.ID = run.TraceID
	done, reason := pe.donePollingTraces(request, traceDB, trace)
	if !done {
		err := pe.eventEmitter.Emit(ctx, events.TracePollingIterationInfo(request.test.ID, request.run.ID, len(trace.Flat), request.count, false, reason))
		if err != nil {
			log.Printf("[PollerExecutor] Test %s Run %d: failed to emit TracePollingIterationInfo event: error: %s\n", request.test.ID, request.run.ID, err.Error())
		}

		log.Printf("[PollerExecutor] Test %s Run %d: Not done polling. (%s)\n", request.test.ID, request.run.ID, reason)
		run.Trace = &trace
		request.run = run
		return false, "", run, nil
	}

	log.Printf("[PollerExecutor] Test %s Run %d: Done polling. (%s)\n", request.test.ID, request.run.ID, reason)

	log.Printf("[PollerExecutor] Test %s Run %d: Start Sorting\n", request.test.ID, request.run.ID)
	trace = trace.Sort()
	log.Printf("[PollerExecutor] Test %s Run %d: Sorting complete\n", request.test.ID, request.run.ID)
	run.Trace = &trace
	request.run = run

	if !trace.HasRootSpan() {
		newRoot := test.NewTracetestRootSpan(run)
		run.Trace = run.Trace.InsertRootSpan(newRoot)
	} else {
		run.Trace.RootSpan = model.AugmentRootSpan(run.Trace.RootSpan, run.TriggerResult)
	}
	run = run.SuccessfullyPolledTraces(run.Trace)

	fmt.Printf("[PollerExecutor] Completed polling process for Test Run %d after %d iterations, number of spans collected: %d \n", run.ID, request.count+1, len(run.Trace.Flat))

	log.Printf("[PollerExecutor] Test %s Run %d: Start updating\n", request.test.ID, request.run.ID)
	err = pe.updater.Update(ctx, run)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: Update error: %s\n", request.test.ID, request.run.ID, err.Error())
		return false, "", test.Run{}, err
	}

	return true, reason, run, nil
}

func (pe DefaultPollerExecutor) donePollingTraces(job *PollingRequest, traceDB tracedb.TraceDB, trace model.Trace) (bool, string) {
	if !traceDB.ShouldRetry() {
		return true, "TraceDB is not retryable"
	}
	pp := pe.ppGetter.GetDefault(job.Context())
	if pp.Periodic == nil {
		return false, "Polling profile not configured"
	}

	maxTracePollRetry := pp.Periodic.MaxTracePollRetry()
	// we're done if we have the same amount of spans after polling or `maxTracePollRetry` times
	if job.count == maxTracePollRetry {
		return true, fmt.Sprintf("Hit MaxRetry of %d", maxTracePollRetry)
	}

	if job.run.Trace == nil {
		return false, "First iteration"
	}

	haveNotCollectedSpansSinceLastPoll := len(trace.Flat) == len(job.run.Trace.Flat)
	haveCollectedSpansInTestRun := len(trace.Flat) > 0
	haveCollectedOnlyRootNode := len(trace.Flat) == 1 && trace.HasRootSpan()

	// Today we consider that we finished collecting traces
	// if we haven't collected any new spans since our last poll
	// and we have collected at least one span for this test run
	// and we have not collected only the root span

	if haveNotCollectedSpansSinceLastPoll && haveCollectedSpansInTestRun && !haveCollectedOnlyRootNode {
		return true, fmt.Sprintf("Trace has no new spans. Spans found: %d", len(trace.Flat))
	}

	return false, fmt.Sprintf("New spans found. Before: %d After: %d", len(job.run.Trace.Flat), len(trace.Flat))
}
