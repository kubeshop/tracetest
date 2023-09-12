package tracepollerworker

import (
	"context"
	"log"
	"fmt"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
)

type tracePollerEvaluatorWorker struct {
	state *workerState
	outputQueue  pipeline.Enqueuer[executor.Job]
}

func NewEvaluatorWorker(
	eventEmitter executor.EventEmitter,
	newTraceDBFn tracedb.FactoryFunc,
	dsRepo resourcemanager.Current[datastore.DataStore],
	updater             executor.RunUpdater,
	subscriptionManager *subscription.Manager,
) *tracePollerEvaluatorWorker {
	state := &workerState{
		eventEmitter: eventEmitter,
		newTraceDBFn: newTraceDBFn,
		dsRepo: dsRepo,
		updater: updater,
		subscriptionManager: subscriptionManager,
	}

	return &tracePollerEvaluatorWorker{state: state}
}

func (w *tracePollerEvaluatorWorker) SetInputQueue(queue pipeline.Enqueuer[executor.Job]) {
	w.state.inputQueue = queue
}

func (w *tracePollerEvaluatorWorker) SetOutputQueue(queue pipeline.Enqueuer[executor.Job]) {
	w.outputQueue = queue
}

// TODO: add instrumentation and selector based features

func (w *tracePollerEvaluatorWorker) ProcessItem(ctx context.Context, job executor.Job) {
	traceDB, err := getTraceDB(ctx, w.state)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: GetDataStore error: %s", job.Test.ID, job.Run.ID, err.Error())
		handleError(ctx, job, err, w.state)
		return
	}

	done, reason := w.donePollingTraces(&job, traceDB, job.Run.Trace)

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

func (w *tracePollerEvaluatorWorker) donePollingTraces(job *executor.Job, traceDB tracedb.TraceDB, trace *traces.Trace) (bool, string) {
	if !traceDB.ShouldRetry() {
		return true, "TraceDB is not retryable"
	}

	maxTracePollRetry := job.PollingProfile.Periodic.MaxTracePollRetry()
	// we're done if we have the same amount of spans after polling or `maxTracePollRetry` times
	log.Printf("[PollerExecutor] Test %s Run %d: Job count %d, max retries: %d", job.Test.ID, job.Run.ID, job.EnqueueCount(), maxTracePollRetry)
	if job.EnqueueCount() >= maxTracePollRetry {
		return true, fmt.Sprintf("Hit MaxRetry of %d", maxTracePollRetry)
	}

	if job.Run.Trace == nil {
		return false, "First iteration"
	}

	haveNotCollectedSpansSinceLastPoll := len(trace.Flat) == len(job.Run.Trace.Flat)
	haveCollectedSpansInTestRun := len(trace.Flat) > 0
	haveCollectedOnlyRootNode := len(trace.Flat) == 1 && trace.HasRootSpan()

	// Today we consider that we finished collecting traces
	// if we haven't collected any new spans since our last poll
	// and we have collected at least one span for this test run
	// and we have not collected only the root span

	if haveNotCollectedSpansSinceLastPoll && haveCollectedSpansInTestRun && !haveCollectedOnlyRootNode {
		return true, fmt.Sprintf("Trace has no new spans. Spans found: %d", len(trace.Flat))
	}

	return false, fmt.Sprintf("New spans found. Before: %d After: %d", len(job.Run.Trace.Flat), len(trace.Flat))
}
