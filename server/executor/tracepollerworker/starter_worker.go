package tracepollerworker

import (
	"context"
	"fmt"
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

type tracePollerStarterWorker struct {
	state       *workerState
	outputQueue pipeline.Enqueuer[executor.Job]
}

func NewStarterWorker(
	eventEmitter executor.EventEmitter,
	newTraceDBFn tracedb.FactoryFunc,
	dsRepo resourcemanager.Current[datastore.DataStore],
	updater executor.RunUpdater,
	subscriptionManager *subscription.Manager,
	tracer trace.Tracer,
) *tracePollerStarterWorker {
	state := &workerState{
		eventEmitter:        eventEmitter,
		newTraceDBFn:        newTraceDBFn,
		dsRepo:              dsRepo,
		updater:             updater,
		subscriptionManager: subscriptionManager,
		tracer:              tracer,
	}

	return &tracePollerStarterWorker{state: state}
}

func (w *tracePollerStarterWorker) SetInputQueue(queue pipeline.Enqueuer[executor.Job]) {
	w.state.inputQueue = queue
}

func (w *tracePollerStarterWorker) SetOutputQueue(queue pipeline.Enqueuer[executor.Job]) {
	w.outputQueue = queue
}

func (w *tracePollerStarterWorker) ProcessItem(ctx context.Context, job executor.Job) {
	ctx, span := w.state.tracer.Start(ctx, "Trace Polling")
	defer span.End()

	select {
	default:
	case <-ctx.Done():
		return
	}

	log.Println("[TracePoller] Starting to poll traces", job.EnqueueCount())

	traceDB, err := getTraceDB(ctx, w.state)
	if err != nil {
		log.Printf("[TracePoller] GetDataStore error: %s", err.Error())
		handleError(ctx, job, err, w.state)
		return
	}

	emitEvent(ctx, w.state, events.TraceFetchingStart(job.Test.ID, job.Run.ID))

	err = w.testConnection(ctx, traceDB, &job)
	if err != nil {
		handleError(ctx, job, err, w.state)
		return
	}

	w.outputQueue.Enqueue(ctx, job)
}

func (w *tracePollerStarterWorker) testConnection(ctx context.Context, traceDB tracedb.TraceDB, job *executor.Job) error {
	if testableTraceDB, ok := traceDB.(tracedb.TestableTraceDB); ok {
		connectionResult := testableTraceDB.TestConnection(ctx)

		emitEvent(ctx, w.state, events.TraceDataStoreConnectionInfo(job.Test.ID, job.Run.ID, connectionResult))
	}

	endpoints := traceDB.GetEndpoints()
	ds, err := w.state.dsRepo.Current(ctx)
	if err != nil {
		return fmt.Errorf("could not get current datastore: %w", err)
	}

	emitEvent(ctx, w.state, events.TracePollingStart(job.Test.ID, job.Run.ID, string(ds.Type), endpoints))

	return nil
}
