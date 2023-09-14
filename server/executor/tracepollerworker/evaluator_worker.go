package tracepollerworker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"

	"go.opentelemetry.io/otel/trace"
)

type PollingStopStrategy interface {
	Evaluate(ctx context.Context, job *executor.Job, traceDB tracedb.TraceDB) (bool, string)
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
	ctx, span := w.state.tracer.Start(ctx, "Evaluating trace")
	defer span.End()

	traceNotFound := job.Headers.GetBool("traceNotFound")

	if traceNotFound && !tracePollerTimedOut(ctx, job) {
		// Edge case: the trace still not available on Data Store during polling, we need to poll/fetch trace again
		populateSpan(span, job, "", nil)

		emitEvent(ctx, w.state, events.TracePollingIterationInfo(job.Test.ID, job.Run.ID, 0, job.EnqueueCount(), false, "trace not found on data store"))
		enqueueTraceFetchJob(ctx, job, w.state)
		return
	}

	// if an error happened on last iteration validate it
	if job.Run.LastError != nil {
		err := job.Run.LastError
		reason := ""

		if traceNotFound && tracePollerTimedOut(ctx, job) {
			reason = fmt.Sprintf("Timed out without finding trace, trace id \"%s\"", job.Run.TraceID.String())
			log.Println("[TracePoller] Timed-out")
		} else {
			reason = fmt.Sprintf("Unexpected error: %s", err.Error())
			log.Println("[TracePoller] Unknown error", err)
		}

		emitEvent(ctx, w.state, events.TracePollingIterationInfo(job.Test.ID, job.Run.ID, 0, job.EnqueueCount(), false, reason))
		emitEvent(ctx, w.state, events.TracePollingError(job.Test.ID, job.Run.ID, reason, err))
		emitEvent(ctx, w.state, events.TraceFetchingError(job.Test.ID, job.Run.ID, err))

		successful := false
		populateSpan(span, job, reason, &successful)

		run := job.Run.TraceFailed(err)
		analytics.SendEvent("test_run_finished", "error", "", &map[string]string{
			"finalState": string(run.State),
		})

		handleDBError(w.state.updater.Update(ctx, run))

		w.state.subscriptionManager.PublishUpdate(subscription.Message{
			ResourceID: run.TransactionStepResourceID(),
			Type:       "update_run",
			Content:    executor.RunResult{Run: run, Err: err},
		})

		handleError(ctx, job, err, w.state, span)
		return
	}

	// otherwise, validate if the polling process should stop
	traceDB, err := getTraceDB(ctx, w.state)
	if err != nil {
		log.Printf("[TracePoller] Test %s Run %d: GetDataStore error: %s", job.Test.ID, job.Run.ID, err.Error())
		handleError(ctx, job, err, w.state, span)
		return
	}

	done, reason := w.stopStrategy.Evaluate(ctx, &job, traceDB)

	populateSpan(span, job, reason, &done)

	if !done { // trace polling is not done, try to fetch trace again
		emitEvent(ctx, w.state, events.TracePollingIterationInfo(job.Test.ID, job.Run.ID, len(job.Run.Trace.Flat), job.EnqueueCount(), false, reason))

		log.Printf("[TracePoller] Test %s Run %d: Not done polling. (%s)", job.Test.ID, job.Run.ID, reason)

		enqueueTraceFetchJob(ctx, job, w.state)
		return
	}
	log.Printf("[TracePoller] Test %s Run %d: Done polling. (%s)", job.Test.ID, job.Run.ID, reason)

	log.Printf("[TracePoller] Test %s Run %d: Start Sorting", job.Test.ID, job.Run.ID)
	sorted := job.Run.Trace.Sort()
	job.Run.Trace = &sorted
	log.Printf("[TracePoller] Test %s Run %d: Sorting complete", job.Test.ID, job.Run.ID)

	if !job.Run.Trace.HasRootSpan() {
		newRoot := test.NewTracetestRootSpan(job.Run)
		job.Run.Trace = job.Run.Trace.InsertRootSpan(newRoot)
	} else {
		job.Run.Trace.RootSpan = traces.AugmentRootSpan(job.Run.Trace.RootSpan, job.Run.TriggerResult)
	}
	job.Run = job.Run.SuccessfullyPolledTraces(job.Run.Trace)

	log.Printf("[TracePoller] Completed polling process for Test Run %d after %d iterations, number of spans collected: %d ", job.Run.ID, job.EnqueueCount()+1, len(job.Run.Trace.Flat))

	log.Printf("[TracePoller] Test %s Run %d: Start updating", job.Test.ID, job.Run.ID)
	handleDBError(w.state.updater.Update(ctx, job.Run))

	emitEvent(ctx, w.state, events.TracePollingSuccess(job.Test.ID, job.Run.ID, reason))
	emitEvent(ctx, w.state, events.TraceFetchingSuccess(job.Test.ID, job.Run.ID))

	log.Printf("[TracePoller] Test %s Run %d: Done polling (reason: %s). Completed polling after %d iterations, number of spans collected %d\n", job.Test.ID, job.Run.ID, reason, job.EnqueueCount()+1, len(job.Run.Trace.Flat))

	w.outputQueue.Enqueue(ctx, job)
}

func tracePollerTimedOut(ctx context.Context, job executor.Job) bool {
	if job.PollingProfile.Periodic == nil {
		return false
	}

	pp := *job.PollingProfile.Periodic
	timedOut := time.Since(job.Run.ServiceTriggeredAt) >= pp.TimeoutDuration()

	return timedOut
}

func enqueueTraceFetchJob(ctx context.Context, job executor.Job, state *workerState) {
	go func() {
		log.Printf("[TracePoller] Requeuing Test Run %d. Current iteration: %d\n", job.Run.ID, job.EnqueueCount())
		time.Sleep(job.PollingProfile.Periodic.RetryDelayDuration())

		job.IncreaseEnqueueCount()
		job.Headers.SetBool("requeued", true)

		select {
		default:
		case <-ctx.Done():
			return // user requested to stop the process
		}

		// inputQueue is set as the trace fetch queue by our pipeline engine
		state.inputQueue.Enqueue(ctx, job)
	}()
}
