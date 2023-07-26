package executor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

const PollingRequestStartIteration = 1

type TracePoller interface {
	Poll(context.Context, test.Test, test.Run, pollingprofile.PollingProfile)
}

type PersistentTracePoller interface {
	TracePoller
	WorkerPool
}

func NewPollResult(finished bool, reason string, run test.Run) PollResult {
	return PollResult{
		finished: finished,
		reason:   reason,
		run:      run,
	}
}

type PollResult struct {
	finished bool
	reason   string
	run      test.Run
}

func (pr PollResult) Finished() bool {
	return pr.finished
}

func (pr PollResult) Reason() string {
	return pr.reason
}

func (pr PollResult) Run() test.Run {
	return pr.run
}

type pollerExecutor interface {
	ExecuteRequest(context.Context, *Job) (PollResult, error)
}

type TraceFetcher interface {
	GetTraceByID(ctx context.Context, traceID string) (*v1.TracesData, error)
}

func NewTracePoller(
	pe pollerExecutor,
	updater RunUpdater,
	subscriptionManager *subscription.Manager,
	eventEmitter EventEmitter,
) *tracePoller {
	return &tracePoller{
		updater:             updater,
		pollerExecutor:      pe,
		subscriptionManager: subscriptionManager,
		eventEmitter:        eventEmitter,
	}
}

type tracePoller struct {
	updater             RunUpdater
	pollerExecutor      pollerExecutor
	subscriptionManager *subscription.Manager
	eventEmitter        EventEmitter
	inputQueue          Enqueuer
	outputQueue         Enqueuer
}

func (tp tracePoller) handleDBError(err error) {
	if err != nil {
		fmt.Printf("DB error when polling traces: %s\n", err.Error())
	}
}

func (tp tracePoller) enqueueJob(ctx context.Context, job Job) {
	tp.inputQueue.Enqueue(ctx, job)
}

func (tp tracePoller) isFirstRequest(job Job) bool {
	return job.EnqueueCount() == 0
}

func (tp *tracePoller) SetOutputQueue(queue Enqueuer) {
	tp.outputQueue = queue
}

func (tp *tracePoller) SetInputQueue(queue Enqueuer) {
	tp.inputQueue = queue
}

func (tp *tracePoller) ProcessItem(ctx context.Context, job Job) {
	if tp.isFirstRequest(job) {
		err := tp.eventEmitter.Emit(ctx, events.TraceFetchingStart(job.Test.ID, job.Run.ID))
		if err != nil {
			log.Printf("[TracePoller] Test %s Run %d: fail to emit TracePollingStart event: %s", job.Test.ID, job.Run.ID, err.Error())
		}
	}

	log.Println("TracePoller] processJob", job.EnqueueCount())

	result, err := tp.pollerExecutor.ExecuteRequest(ctx, &job)
	run := result.run
	if err != nil {
		log.Printf("[TracePoller] Test %s Run %d: ExecuteRequest Error: %s", job.Test.ID, job.Run.ID, err.Error())
		jobFailed, reason := tp.handleTraceDBError(ctx, job, err)

		if jobFailed {
			anotherErr := tp.eventEmitter.Emit(ctx, events.TracePollingError(job.Test.ID, job.Run.ID, reason, err))
			if anotherErr != nil {
				log.Printf("[TracePoller] Test %s Run %d: fail to emit TracePollingError event: %s \n", job.Test.ID, job.Run.ID, err.Error())
			}

			anotherErr = tp.eventEmitter.Emit(ctx, events.TraceFetchingError(job.Test.ID, job.Run.ID, err))
			if anotherErr != nil {
				log.Printf("[TracePoller] Test %s Run %d: fail to emit TracePollingError event: %s \n", job.Test.ID, job.Run.ID, err.Error())
			}
		}

		return
	}

	if !result.finished {
		tp.requeue(ctx, job)
		return
	}

	err = tp.eventEmitter.Emit(ctx, events.TracePollingSuccess(job.Test.ID, job.Run.ID, result.reason))
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: failed to emit TracePollingSuccess event: error: %s\n", job.Test.ID, job.Run.ID, err.Error())
	}

	log.Printf("[TracePoller] Test %s Run %d: Done polling (reason: %s). Completed polling after %d iterations, number of spans collected %d\n", job.Test.ID, job.Run.ID, result.reason, job.EnqueueCount()+1, len(run.Trace.Flat))

	err = tp.eventEmitter.Emit(ctx, events.TraceFetchingSuccess(job.Test.ID, job.Run.ID))
	if err != nil {
		log.Printf("[TracePoller] Test %s Run %d: fail to emit TracePollingSuccess event: %s \n", job.Test.ID, job.Run.ID, err.Error())
	}

	tp.handleDBError(tp.updater.Update(ctx, run))
	tp.outputQueue.Enqueue(ctx, job)
}

func (tp tracePoller) handleTraceDBError(ctx context.Context, job Job, err error) (bool, string) {
	run := job.Run

	if job.PollingProfile.Periodic == nil {
		log.Println("[TracePoller] cannot get polling profile.")
		return true, "Cannot get polling profile"
	}

	pp := *job.PollingProfile.Periodic

	// Edge case: the trace still not avaiable on Data Store during polling
	if errors.Is(err, connection.ErrTraceNotFound) && time.Since(run.ServiceTriggeredAt) < pp.TimeoutDuration() {
		log.Println("[TracePoller] Trace not found on Data Store yet. Requeuing...")
		tp.requeue(ctx, job)
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

	tp.handleDBError(tp.updater.Update(ctx, run))

	tp.subscriptionManager.PublishUpdate(subscription.Message{
		ResourceID: run.TransactionStepResourceID(),
		Type:       "update_run",
		Content:    RunResult{Run: run, Err: err},
	})

	return true, reason // job failed
}

func (tp tracePoller) requeue(ctx context.Context, job Job) {
	go func() {
		fmt.Printf("[TracePoller] Requeuing Test Run %d. Current iteration: %d\n", job.Run.ID, job.EnqueueCount())
		time.Sleep(job.PollingProfile.Periodic.RetryDelayDuration())

		job.IncreaseEnqueueCount()
		job.Headers.SetBool("requeued", true)
		tp.enqueueJob(ctx, job)
	}()
}

func isFirstRequest(job *Job) bool {
	return !job.Headers.GetBool("requeued")
}
