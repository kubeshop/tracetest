package executor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type traceDBFactoryFn func(ds model.DataStore) (tracedb.TraceDB, error)

type DefaultPollerExecutor struct {
	pollingProfileRepo pollingprofile.Repository
	updater            RunUpdater
	newTraceDBFn       traceDBFactoryFn
	dsRepo             model.DataStoreRepository
}

type InstrumentedPollerExecutor struct {
	tracer         trace.Tracer
	pollerExecutor PollerExecutor
}

func (pe InstrumentedPollerExecutor) ExecuteRequest(request *PollingRequest) (bool, model.Run, error) {
	_, span := pe.tracer.Start(request.ctx, "Fetch trace")
	defer span.End()

	finished, run, err := pe.pollerExecutor.ExecuteRequest(request)

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

	if err != nil {
		attrs = append(attrs, attribute.String("tracetest.run.trace_poller.error", err.Error()))
		span.RecordError(err)
	}

	span.SetAttributes(attrs...)
	return finished, run, err
}

func NewPollerExecutor(
	pollingProfileRepo pollingprofile.Repository,
	tracer trace.Tracer,
	updater RunUpdater,
	newTraceDBFn traceDBFactoryFn,
	dsRepo model.DataStoreRepository,
) PollerExecutor {
	pollerExecutor := &DefaultPollerExecutor{
		pollingProfileRepo: pollingProfileRepo,
		updater:            updater,
		newTraceDBFn:       newTraceDBFn,
		dsRepo:             dsRepo,
	}

	return &InstrumentedPollerExecutor{
		tracer:         tracer,
		pollerExecutor: pollerExecutor,
	}
}

func (pe DefaultPollerExecutor) traceDB(ctx context.Context) (tracedb.TraceDB, error) {
	ds, err := pe.dsRepo.DefaultDataStore(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get default datastore: %w", err)
	}

	tdb, err := pe.newTraceDBFn(ds)
	if err != nil {
		return nil, fmt.Errorf(`cannot get tracedb from DataStore config with ID "%s": %w`, ds.ID, err)
	}

	return tdb, nil
}

func (pe DefaultPollerExecutor) ExecuteRequest(request *PollingRequest) (bool, model.Run, error) {
	log.Printf("[PollerExecutor] Test %s Run %d: ExecuteRequest\n", request.test.ID, request.run.ID)
	run := request.run

	traceDB, err := pe.traceDB(request.ctx)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: GetDataStore error: %s\n", request.test.ID, request.run.ID, err.Error())
		return false, model.Run{}, err
	}

	traceID := run.TraceID.String()
	trace, err := traceDB.GetTraceByID(request.ctx, traceID)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: GetTraceByID (traceID %s) error: %s\n", request.test.ID, request.run.ID, traceID, err.Error())
		return false, model.Run{}, err
	}

	trace.ID = run.TraceID

	done, reason := pe.donePollingTraces(request, traceDB, trace)

	if !done {
		log.Printf("[PollerExecutor] Test %s Run %d: Not done polling. (%s)\n", request.test.ID, request.run.ID, reason)
		run.Trace = &trace
		request.run = run
		return false, run, nil
	}

	log.Printf("[PollerExecutor] Test %s Run %d: Done polling. (%s)\n", request.test.ID, request.run.ID, reason)

	log.Printf("[PollerExecutor] Test %s Run %d: Start Sorting\n", request.test.ID, request.run.ID)
	trace = trace.Sort()
	log.Printf("[PollerExecutor] Test %s Run %d: Sorting complete\n", request.test.ID, request.run.ID)
	run.Trace = &trace
	request.run = run

	if !trace.HasRootSpan() {
		newRoot := model.NewTracetestRootSpan(run)
		run.Trace = run.Trace.InsertRootSpan(newRoot)
	} else {
		run.Trace.RootSpan = model.AugmentRootSpan(run.Trace.RootSpan, run.TriggerResult)
	}
	run = run.SuccessfullyPolledTraces(run.Trace)

	fmt.Printf("[PollerExecutor] Completed polling process for Test Run %d after %d iterations, number of spans collected: %d \n", run.ID, request.count+1, len(run.Trace.Flat))

	log.Printf("[PollerExecutor] Test %s Run %d: Start updating\n", request.test.ID, request.run.ID)
	err = pe.updater.Update(request.ctx, run)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: Update error: %s\n", request.test.ID, request.run.ID, err.Error())
		return false, model.Run{}, err
	}

	return true, run, nil
}

func (pe DefaultPollerExecutor) donePollingTraces(job *PollingRequest, traceDB tracedb.TraceDB, trace model.Trace) (bool, string) {
	if !traceDB.ShouldRetry() {
		return true, "TraceDB is not retryable"
	}

	pollingProfile := pe.getDefaultPollingProfile(job.ctx)

	// we're done if we have the same amount of spans after polling or `maxRetries` times
	maxRetries := pollingProfile.Periodic.GetMaxRetries()
	if job.count == maxRetries {
		return true, fmt.Sprintf("Hit MaxRetry of %d", maxRetries)
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

var defaultPollingProfile = pollingprofile.PollingProfile{
	ID:       id.ID("current"),
	Name:     "default",
	Default:  true,
	Strategy: pollingprofile.Periodic,
	Periodic: &pollingprofile.PeriodicPollingConfig{
		Timeout:    "1m",
		RetryDelay: "5s",
	},
}

func (pe DefaultPollerExecutor) getDefaultPollingProfile(ctx context.Context) pollingprofile.PollingProfile {
	profile, err := pe.pollingProfileRepo.GetDefault(ctx)
	if err != nil && errors.Unwrap(err) == sql.ErrNoRows {
		return defaultPollingProfile
	}

	return profile
}
