package tracepollerworker

import (
	"context"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/tracedb"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type workerState struct {
	eventEmitter        executor.EventEmitter
	newTraceDBFn        tracedb.FactoryFunc
	dsRepo              resourcemanager.Current[datastore.DataStore]
	updater             executor.RunUpdater
	subscriptionManager *subscription.Manager
	tracer              trace.Tracer
	inputQueue          pipeline.Enqueuer[executor.Job]
}

func emitEvent(ctx context.Context, state *workerState, event model.TestRunEvent) {
	err := state.eventEmitter.Emit(ctx, event)
	if err != nil {
		log.Printf("[TracePoller] failed to emit %s event: error: %s", event.Type, err.Error())
	}
}

func getTraceDB(ctx context.Context, state *workerState) (tracedb.TraceDB, error) {
	ds, err := state.dsRepo.Current(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get default datastore: %w", err)
	}

	tdb, err := state.newTraceDBFn(ds)
	if err != nil {
		return nil, fmt.Errorf(`cannot get tracedb from DataStore config with ID "%s": %w`, ds.ID, err)
	}

	return tdb, nil
}

func handleError(ctx context.Context, job executor.Job, err error, state *workerState, span trace.Span) {
	log.Printf("[TracePoller] Test %s Run %d, Error: %s", job.Test.ID, job.Run.ID, err.Error())

	span.RecordError(err)
	span.SetAttributes(attribute.String("tracetest.run.trace_poller.error", err.Error()))
}

func handleDBError(err error) {
	if err == nil {
		return
	}

	log.Printf("[TracePoller] DB error when polling traces: %s\n", err.Error())
}

func populateSpan(span trace.Span, job executor.Job, reason string, done *bool) {
	spanCount := 0
	if job.Run.Trace != nil {
		spanCount = len(job.Run.Trace.Flat)
	}

	attrs := []attribute.KeyValue{
		attribute.String("tracetest.run.trace_poller.trace_id", job.Run.TraceID.String()),
		attribute.String("tracetest.run.trace_poller.span_id", job.Run.SpanID.String()),
		attribute.String("tracetest.run.trace_poller.test_id", string(job.Test.ID)),
		attribute.Int("tracetest.run.trace_poller.amount_retrieved_spans", spanCount),
	}

	if done != nil {
		attrs = append(attrs, attribute.Bool("tracetest.run.trace_poller.succesful", *done))
	}

	if reason != "" {
		attrs = append(attrs, attribute.String("tracetest.run.trace_poller.finish_reason", reason))
	}

	span.SetAttributes(attrs...)
}
