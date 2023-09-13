package tracepollerworker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/tracedb/connection"

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

func handleError(ctx context.Context, job executor.Job, err error, state *workerState) {
	log.Printf("[PollerExecutor] Test %s Run %d: Update error: %s", job.Test.ID, job.Run.ID, err.Error())

	log.Printf("[TracePoller] Test %s Run %d: ExecuteRequest Error: %s", job.Test.ID, job.Run.ID, err.Error())
	jobFailed, reason := handleTraceDBError(ctx, job, err, state)

	if jobFailed {
		emitEvent(ctx, state, events.TracePollingError(job.Test.ID, job.Run.ID, reason, err))
		emitEvent(ctx, state, events.TraceFetchingError(job.Test.ID, job.Run.ID, err))
	}
}

func handleTraceDBError(ctx context.Context, job executor.Job, err error, state *workerState) (bool, string) {
	run := job.Run

	if job.PollingProfile.Periodic == nil {
		log.Println("[TracePoller] cannot get polling profile.")
		return true, "Cannot get polling profile"
	}

	pp := *job.PollingProfile.Periodic

	// Edge case: the trace still not avaiable on Data Store during polling
	if errors.Is(err, connection.ErrTraceNotFound) && time.Since(run.ServiceTriggeredAt) < pp.TimeoutDuration() {
		log.Println("[TracePoller] Trace not found on Data Store yet. Requeuing...")
		requeue(ctx, job, state)
		return false, "Trace not found" // job without fail
	}

	reason := ""

	if errors.Is(err, connection.ErrTraceNotFound) {
		reason = fmt.Sprintf("Timed out without finding trace, trace id \"%s\"", run.TraceID.String())

		err = fmt.Errorf("timed out waiting for traces after %s", pp.Timeout)
		fmt.Println("[TracePoller] Timed-out", err)
	} else {
		reason = "Unexpected error"

		err = fmt.Errorf("cannot fetch trace: %w", err)
		fmt.Println("[TracePoller] Unknown error", err)
	}

	run = run.TraceFailed(err)
	analytics.SendEvent("test_run_finished", "error", "", &map[string]string{
		"finalState": string(run.State),
	})

	handleDBError(state.updater.Update(ctx, run))

	state.subscriptionManager.PublishUpdate(subscription.Message{
		ResourceID: run.TransactionStepResourceID(),
		Type:       "update_run",
		Content:    executor.RunResult{Run: run, Err: err},
	})

	return true, reason // job failed
}

func requeue(ctx context.Context, job executor.Job, state *workerState) {
	go func() {
		fmt.Printf("[TracePoller] Requeuing Test Run %d. Current iteration: %d\n", job.Run.ID, job.EnqueueCount())
		time.Sleep(job.PollingProfile.Periodic.RetryDelayDuration())

		job.IncreaseEnqueueCount()
		job.Headers.SetBool("requeued", true)

		select {
		default:
		case <-ctx.Done():
			return
		}

		state.inputQueue.Enqueue(ctx, job)
	}()
}

func handleDBError(err error) {
	if err != nil {
		fmt.Printf("DB error when polling traces: %s\n", err.Error())
	}
}

// func isFirstRequest(job *executor.Job) bool {
// 	return !job.Headers.GetBool("requeued")
// }
