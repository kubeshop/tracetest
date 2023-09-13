package tracepollerworker

import (
	"context"
	"log"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type PollingStopStrategy interface {
	Evaluate(ctx context.Context, job *executor.Job, traceDB tracedb.TraceDB, trace *traces.Trace) (bool, string)
}

type tracePollerEvaluatorWorker struct {
	state        *workerState
	outputQueue  pipeline.Enqueuer[executor.Job]
	stopStrategy PollingStopStrategy
}

func NewEvaluatorWorker(
	eventEmitter executor.EventEmitter,
	newTraceDBFn tracedb.FactoryFunc,
	dsRepo resourcemanager.Current[datastore.DataStore],
	updater executor.RunUpdater,
	subscriptionManager *subscription.Manager,
	stopStrategy PollingStopStrategy,
	tracer trace.Tracer,
) *tracePollerEvaluatorWorker {
	state := &workerState{
		eventEmitter:        eventEmitter,
		newTraceDBFn:        newTraceDBFn,
		dsRepo:              dsRepo,
		updater:             updater,
		subscriptionManager: subscriptionManager,
		tracer:              tracer,
	}

	return &tracePollerEvaluatorWorker{state: state, stopStrategy: stopStrategy}
}

func (w *tracePollerEvaluatorWorker) SetInputQueue(queue pipeline.Enqueuer[executor.Job]) {
	w.state.inputQueue = queue
}

func (w *tracePollerEvaluatorWorker) SetOutputQueue(queue pipeline.Enqueuer[executor.Job]) {
	w.outputQueue = queue
}

func (w *tracePollerEvaluatorWorker) ProcessItem(ctx context.Context, job executor.Job) {
	ctx, span := w.state.tracer.Start(ctx, "Trace Evaluate")
	defer span.End()

	traceDB, err := getTraceDB(ctx, w.state)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: GetDataStore error: %s", job.Test.ID, job.Run.ID, err.Error())
		handleError(ctx, job, err, w.state)
		return
	}

	done, reason := w.stopStrategy.Evaluate(ctx, &job, traceDB, job.Run.Trace)

	spanCount := 0
	if job.Run.Trace != nil {
		spanCount = len(job.Run.Trace.Flat)
	}

	attrs := []attribute.KeyValue{
		attribute.String("tracetest.run.trace_poller.trace_id", job.Run.TraceID.String()),
		attribute.String("tracetest.run.trace_poller.span_id", job.Run.SpanID.String()),
		attribute.Bool("tracetest.run.trace_poller.succesful", done),
		attribute.String("tracetest.run.trace_poller.test_id", string(job.Test.ID)),
		attribute.Int("tracetest.run.trace_poller.amount_retrieved_spans", spanCount),
	}

	if reason != "" {
		attrs = append(attrs, attribute.String("tracetest.run.trace_poller.finish_reason", reason))
	}

	if err != nil {
		attrs = append(attrs, attribute.String("tracetest.run.trace_poller.error", err.Error()))
		span.RecordError(err)
	}

	span.SetAttributes(attrs...)

	if !done {
		err := w.state.eventEmitter.Emit(ctx, events.TracePollingIterationInfo(job.Test.ID, job.Run.ID, len(job.Run.Trace.Flat), job.EnqueueCount(), false, reason))
		if err != nil {
			log.Printf("[PollerExecutor] Test %s Run %d: failed to emit TracePollingIterationInfo event: error: %s", job.Test.ID, job.Run.ID, err.Error())
		}
		log.Printf("[PollerExecutor] Test %s Run %d: Not done polling. (%s)", job.Test.ID, job.Run.ID, reason)

		requeue(ctx, job, w.state)
		return
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
	err = w.state.updater.Update(ctx, job.Run)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: Update error: %s", job.Test.ID, job.Run.ID, err.Error())
		handleError(ctx, job, err, w.state)
		return
	}

	err = w.state.eventEmitter.Emit(ctx, events.TracePollingSuccess(job.Test.ID, job.Run.ID, reason))
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: failed to emit TracePollingSuccess event: error: %s\n", job.Test.ID, job.Run.ID, err.Error())
	}

	log.Printf("[TracePoller] Test %s Run %d: Done polling (reason: %s). Completed polling after %d iterations, number of spans collected %d\n", job.Test.ID, job.Run.ID, reason, job.EnqueueCount()+1, len(job.Run.Trace.Flat))

	err = w.state.eventEmitter.Emit(ctx, events.TraceFetchingSuccess(job.Test.ID, job.Run.ID))
	if err != nil {
		log.Printf("[TracePoller] Test %s Run %d: fail to emit TracePollingSuccess event: %s \n", job.Test.ID, job.Run.ID, err.Error())
	}

	handleDBError(w.state.updater.Update(ctx, job.Run))

	w.outputQueue.Enqueue(ctx, job)
}
