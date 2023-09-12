package tracepollerworker

import (
	"context"
	"log"
	"fmt"

	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/executor"
)

type tracePollerStarterWorker struct {
	state *workerState
	outputQueue  pipeline.Enqueuer[executor.Job]
}

func NewStarterWorker(
	eventEmitter executor.EventEmitter,
	newTraceDBFn tracedb.FactoryFunc,
	dsRepo resourcemanager.Current[datastore.DataStore],
	updater             executor.RunUpdater,
	subscriptionManager *subscription.Manager,
) *tracePollerStarterWorker {
	state := &workerState{
		eventEmitter: eventEmitter,
		newTraceDBFn: newTraceDBFn,
		dsRepo: dsRepo,
		updater: updater,
		subscriptionManager: subscriptionManager,
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
	select {
		default:
		case <-ctx.Done():
			return
	}

	log.Println("[TracePoller] processJob", job.EnqueueCount())

	log.Printf("[PollerExecutor] Test %s Run %d: ExecuteRequest", job.Test.ID, job.Run.ID)

	traceDB, err := getTraceDB(ctx, w.state)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: GetDataStore error: %s", job.Test.ID, job.Run.ID, err.Error())
		handleError(ctx, job, err, w.state)
		return
	}

	if isFirstRequest(&job) {
		err := w.state.eventEmitter.Emit(ctx, events.TraceFetchingStart(job.Test.ID, job.Run.ID))
		if err != nil {
			log.Printf("[TracePoller] Test %s Run %d: fail to emit TracePollingStart event: %s", job.Test.ID, job.Run.ID, err.Error())
		}

		err = w.testConnection(ctx, traceDB, &job)
		if err != nil {
			handleError(ctx, job, err, w.state)
			return
		}
	}

	w.outputQueue.Enqueue(ctx, job)
}

func (w *tracePollerStarterWorker) testConnection(ctx context.Context, traceDB tracedb.TraceDB, job *executor.Job) error {
	if testableTraceDB, ok := traceDB.(tracedb.TestableTraceDB); ok {
		connectionResult := testableTraceDB.TestConnection(ctx)

		err := w.state.eventEmitter.Emit(ctx, events.TraceDataStoreConnectionInfo(job.Test.ID, job.Run.ID, connectionResult))
		if err != nil {
			log.Printf("[PollerExecutor] Test %s Run %d: failed to emit TraceDataStoreConnectionInfo event: error: %s", job.Test.ID, job.Run.ID, err.Error())
		}
	}

	endpoints := traceDB.GetEndpoints()
	ds, err := w.state.dsRepo.Current(ctx)
	if err != nil {
		return fmt.Errorf("could not get current datastore: %w", err)
	}

	err = w.state.eventEmitter.Emit(ctx, events.TracePollingStart(job.Test.ID, job.Run.ID, string(ds.Type), endpoints))
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: failed to emit TracePollingStart event: error: %s", job.Test.ID, job.Run.ID, err.Error())
	}

	return nil
}
