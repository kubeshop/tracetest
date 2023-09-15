package tracepollerworker

import (
	"context"
	"log"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

type traceFetcherWorker struct {
	state       *workerState
	outputQueue pipeline.Enqueuer[executor.Job]
	enabled     bool
}

func NewFetcherWorker(
	eventEmitter executor.EventEmitter,
	newTraceDBFn tracedb.FactoryFunc,
	dsRepo resourcemanager.Current[datastore.DataStore],
	updater executor.RunUpdater,
	subscriptionManager *subscription.Manager,
	tracer trace.Tracer,
	enabled bool,
) *traceFetcherWorker {
	state := &workerState{
		eventEmitter:        eventEmitter,
		newTraceDBFn:        newTraceDBFn,
		dsRepo:              dsRepo,
		updater:             updater,
		subscriptionManager: subscriptionManager,
		tracer:              tracer,
	}

	return &traceFetcherWorker{state: state, enabled: enabled}
}

func (w *traceFetcherWorker) SetInputQueue(queue pipeline.Enqueuer[executor.Job]) {
	w.state.inputQueue = queue
}

func (w *traceFetcherWorker) SetOutputQueue(queue pipeline.Enqueuer[executor.Job]) {
	w.outputQueue = queue
}

func (w *traceFetcherWorker) ProcessItem(ctx context.Context, job executor.Job) {
	if !w.enabled {
		w.outputQueue.Enqueue(ctx, job)
		return
	}

	ctx, span := w.state.tracer.Start(ctx, "Fetching trace")
	defer span.End()

	populateSpan(span, job, "", nil)

	traceDB, err := getTraceDB(ctx, w.state)
	if err != nil {
		log.Printf("[TracePoller] Test %s Run %d: GetDataStore error: %s", job.Test.ID, job.Run.ID, err.Error())
		handleError(ctx, job, err, w.state, span)
		return
	}

	traceID := job.Run.TraceID.String()
	trace, err := traceDB.GetTraceByID(ctx, traceID)
	if err != nil {
		log.Printf("[TracePoller] Test %s Run %d: GetTraceByID (traceID %s) error: %s", job.Test.ID, job.Run.ID, traceID, err.Error())

		if isTraceNotFoundError(err) {
			job.Headers.SetBool("traceNotFound", true)
		} else {
			job.Run.LastError = err
			handleDBError(w.state.updater.Update(ctx, job.Run))
		}

		w.outputQueue.Enqueue(ctx, job)
		return
	}

	spansBefore := 0
	if job.Run.Trace != nil {
		spansBefore = len(job.Run.Trace.Flat)
	}

	collectedSpans := len(trace.Flat) - spansBefore
	job.Headers.SetInt("collectedSpans", collectedSpans)
	job.Headers.SetBool("traceNotFound", false)

	trace.ID = job.Run.TraceID
	job.Run.Trace = &trace

	handleDBError(w.state.updater.Update(ctx, job.Run))
	w.outputQueue.Enqueue(ctx, job)
}
