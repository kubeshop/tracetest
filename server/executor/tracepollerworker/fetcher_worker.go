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
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

type traceFetcherWorker struct {
	state       *workerState
	outputQueue pipeline.Enqueuer[executor.Job]
}

func NewFetcherWorker(
	eventEmitter executor.EventEmitter,
	newTraceDBFn tracedb.FactoryFunc,
	dsRepo resourcemanager.Current[datastore.DataStore],
	updater executor.RunUpdater,
	subscriptionManager *subscription.Manager,
	tracer trace.Tracer,
) *traceFetcherWorker {
	state := &workerState{
		eventEmitter:        eventEmitter,
		newTraceDBFn:        newTraceDBFn,
		dsRepo:              dsRepo,
		updater:             updater,
		subscriptionManager: subscriptionManager,
		tracer:              tracer,
	}

	return &traceFetcherWorker{state: state}
}

func (w *traceFetcherWorker) SetInputQueue(queue pipeline.Enqueuer[executor.Job]) {
	w.state.inputQueue = queue
}

func (w *traceFetcherWorker) SetOutputQueue(queue pipeline.Enqueuer[executor.Job]) {
	w.outputQueue = queue
}

func (w *traceFetcherWorker) ProcessItem(ctx context.Context, job executor.Job) {
	ctx, span := w.state.tracer.Start(ctx, "Trace Fetching")
	defer span.End()

	traceDB, err := getTraceDB(ctx, w.state)
	if err != nil {
		log.Printf("[TracePoller] Test %s Run %d: GetDataStore error: %s", job.Test.ID, job.Run.ID, err.Error())
		handleError(ctx, job, err, w.state)
		return
	}

	traceID := job.Run.TraceID.String()
	trace, err := traceDB.GetTraceByID(ctx, traceID)
	if err != nil {
		log.Printf("[TracePoller] Test %s Run %d: GetTraceByID (traceID %s) error: %s", job.Test.ID, job.Run.ID, traceID, err.Error())

		emitEvent(ctx, w.state, events.TracePollingIterationInfo(job.Test.ID, job.Run.ID, 0, job.EnqueueCount(), false, err.Error()))

		handleError(ctx, job, err, w.state)
		return
	}

	trace.ID = job.Run.TraceID
	job.Run.Trace = &trace

	err = w.state.updater.Update(ctx, job.Run)
	if err != nil {
		handleError(ctx, job, err, w.state)
		return
	}

	w.outputQueue.Enqueue(ctx, job)
}
