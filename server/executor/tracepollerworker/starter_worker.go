package tracepollerworker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/testconnection"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

type tracePollerStarterWorker struct {
	state          *workerState
	dsTestPipeline dataStorePipeline
	outputQueue    pipeline.Enqueuer[executor.Job]
}

type dataStorePipeline interface {
	Run(context.Context, testconnection.Job)
	NewJob(context.Context, datastore.DataStore) testconnection.Job
	Subscribe(string, testconnection.NotifierFn) error
	Unsubscribe(string)
}

func NewStarterWorker(
	eventEmitter executor.EventEmitter,
	newTraceDBFn tracedb.FactoryFunc,
	dsRepo resourcemanager.Current[datastore.DataStore],
	updater executor.RunUpdater,
	subscriptionManager subscription.Manager,
	tracer trace.Tracer,
	dsTestPipeline dataStorePipeline,
) *tracePollerStarterWorker {
	state := &workerState{
		eventEmitter:        eventEmitter,
		newTraceDBFn:        newTraceDBFn,
		dsRepo:              dsRepo,
		updater:             updater,
		subscriptionManager: subscriptionManager,
		tracer:              tracer,
	}

	return &tracePollerStarterWorker{
		state:          state,
		dsTestPipeline: dsTestPipeline, // this is necessary just for this worker
	}
}

func (w *tracePollerStarterWorker) SetInputQueue(queue pipeline.Enqueuer[executor.Job]) {
	w.state.inputQueue = queue
}

func (w *tracePollerStarterWorker) SetOutputQueue(queue pipeline.Enqueuer[executor.Job]) {
	w.outputQueue = queue
}

func (w *tracePollerStarterWorker) ProcessItem(ctx context.Context, job executor.Job) {
	ctx, span := w.state.tracer.Start(ctx, "Start polling trace")
	defer span.End()

	if job.Run.SkipTraceCollection {
		emitEvent(ctx, w.state, events.TracePollingSkipped(job.Test.ID, job.Run.ID))
		w.outputQueue.Enqueue(ctx, job)
		return
	}

	populateSpan(span, job, "", nil)

	select {
	default:
	case <-ctx.Done():
		err := context.Cause(ctx)
		if errors.Is(err, executor.ErrSkipTraceCollection) {
			ctx = context.Background()
			emitEvent(ctx, w.state, events.TracePollingSkipped(job.Test.ID, job.Run.ID))
			w.outputQueue.Enqueue(ctx, job)
		}

		return
	}

	log.Println("[TracePoller] Starting to poll traces", job.EnqueueCount())

	traceDB, err := getTraceDB(ctx, w.state)
	if err != nil {
		log.Printf("[TracePoller] GetDataStore error: %s", err.Error())
		handleError(ctx, job, err, w.state, span)
		return
	}

	emitEvent(ctx, w.state, events.TraceFetchingStart(job.Test.ID, job.Run.ID))

	endpoints := traceDB.GetEndpoints()
	ds, err := w.state.dsRepo.Current(ctx)
	if err != nil {
		wrappedError := fmt.Errorf("could not get current datastore: %w", err)
		handleError(ctx, job, wrappedError, w.state, span)
		return
	}

	connectionResult, err := w.testConnection(ctx, traceDB, ds)
	if err != nil {
		log.Printf("[TracePoller] TestConnection error: %s", err.Error())
		handleError(ctx, job, err, w.state, span)
		return
	}

	if connectionResult != nil {
		emitEvent(ctx, w.state, events.TraceDataStoreConnectionInfo(job.Test.ID, job.Run.ID, *connectionResult))
	}

	emitEvent(ctx, w.state, events.TracePollingStart(job.Test.ID, job.Run.ID, string(ds.Type), endpoints))

	w.outputQueue.Enqueue(ctx, job)
}

func (w *tracePollerStarterWorker) testConnection(ctx context.Context, traceDB tracedb.TraceDB, ds datastore.DataStore) (*model.ConnectionResult, error) {
	_, ok := traceDB.(tracedb.TestableTraceDB)
	if !ok {
		return nil, nil
	}

	job := w.dsTestPipeline.NewJob(ctx, ds)

	wg := sync.WaitGroup{}
	err := w.dsTestPipeline.Subscribe(job.ID, func(result testconnection.Job) {
		job = result
		wg.Done()
	})

	if err != nil {
		return nil, err
	}

	w.dsTestPipeline.Run(ctx, job)
	wg.Add(1)
	wg.Wait()
	w.dsTestPipeline.Unsubscribe(job.ID)

	return &job.TestResult, nil
}
