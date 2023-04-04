package executor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type TracePoller interface {
	Poll(context.Context, model.Test, model.Run)
}

type PersistentTracePoller interface {
	TracePoller
	WorkerPool
}

type PollerExecutor interface {
	ExecuteRequest(*PollingRequest) (bool, model.Run, error)
}

type TraceFetcher interface {
	GetTraceByID(ctx context.Context, traceID string) (*v1.TracesData, error)
}

type PollingProfileGetter interface {
	GetDefault(ctx context.Context) pollingprofile.PollingProfile
}

func NewTracePoller(
	pe PollerExecutor,
	ppGetter PollingProfileGetter,
	updater RunUpdater,
	assertionRunner AssertionRunner,
	subscriptionManager *subscription.Manager,
) PersistentTracePoller {
	return tracePoller{
		updater:             updater,
		ppGetter:            ppGetter,
		pollerExecutor:      pe,
		executeQueue:        make(chan PollingRequest, 5),
		exit:                make(chan bool, 1),
		assertionRunner:     assertionRunner,
		subscriptionManager: subscriptionManager,
	}
}

type tracePoller struct {
	updater         RunUpdater
	ppGetter        PollingProfileGetter
	pollerExecutor  PollerExecutor
	assertionRunner AssertionRunner

	subscriptionManager *subscription.Manager

	executeQueue chan PollingRequest
	exit         chan bool
}

type PollingRequest struct {
	ctx   context.Context
	test  model.Test
	run   model.Run
	count int
}

func NewPollingRequest(ctx context.Context, test model.Test, run model.Run, count int) *PollingRequest {
	return &PollingRequest{
		ctx:   ctx,
		test:  test,
		run:   run,
		count: count,
	}
}

func (tp tracePoller) handleDBError(err error) {
	if err != nil {
		fmt.Printf("DB error when polling traces: %s\n", err.Error())
	}
}

func (tp tracePoller) Start(workers int) {
	for i := 0; i < workers; i++ {
		go func() {
			fmt.Println("tracePoller start goroutine")
			for {
				select {
				case <-tp.exit:
					fmt.Println("tracePoller exit goroutine")
					return
				case job := <-tp.executeQueue:
					log.Printf("[TracePoller] Test %s Run %d: Received job\n", job.test.ID, job.run.ID)
					tp.processJob(job)
				}
			}
		}()
	}
}

func (tp tracePoller) Stop() {
	fmt.Println("tracePoller stopping")
	tp.exit <- true
}

func (tp tracePoller) Poll(ctx context.Context, test model.Test, run model.Run) {
	log.Printf("[TracePoller] Test %s Run %d: Poll\n", test.ID, run.ID)

	job := NewPollingRequest(ctx, test, run, 0)

	tp.enqueueJob(*job)
}

func (tp tracePoller) enqueueJob(job PollingRequest) {
	tp.executeQueue <- job
}

func (tp tracePoller) processJob(job PollingRequest) {
	finished, run, err := tp.pollerExecutor.ExecuteRequest(&job)
	if err != nil {
		log.Printf("[TracePoller] Test %s Run %d: ExecuteRequest Error: %s\n", job.test.ID, job.run.ID, err.Error())
		tp.handleTraceDBError(job, err)
		return
	}

	if !finished {
		job.count += 1
		tp.requeue(job)
		return
	}

	log.Printf("[TracePoller] Test %s Run %d: Done polling. Completed polling after %d iterations, number of spans collected %d\n", job.test.ID, job.run.ID, job.count+1, len(run.Trace.Flat))

	tp.handleDBError(tp.updater.Update(job.ctx, run))

	job.run = run
	tp.runAssertions(job)
}

func (tp tracePoller) runAssertions(job PollingRequest) {
	assertionRequest := AssertionRequest{
		Test: job.test,
		Run:  job.run,
	}

	tp.assertionRunner.RunAssertions(job.ctx, assertionRequest)
}

func (tp tracePoller) handleTraceDBError(job PollingRequest, err error) {
	run := job.run

	pp := *tp.ppGetter.GetDefault(job.ctx).Periodic

	// Edge case: the trace still not avaiable on Data Store during polling
	if errors.Is(err, connection.ErrTraceNotFound) && time.Since(run.ServiceTriggeredAt) < pp.TimeoutDuration() {
		log.Println("[TracePoller] Trace not found on Data Store yet. Requeuing...")
		tp.requeue(job)
		return
	}

	if errors.Is(err, connection.ErrTraceNotFound) {
		err = fmt.Errorf("timed out waiting for traces after %s", pp.Timeout)
		fmt.Println("[TracePoller] Timed-out", err)
	} else {
		err = fmt.Errorf("cannot fetch trace: %w", err)
		fmt.Println("[TracePoller] Unknown error", err)
	}

	tp.handleDBError(tp.updater.Update(job.ctx, run.TraceFailed(err)))

	tp.subscriptionManager.PublishUpdate(subscription.Message{
		ResourceID: run.TransactionStepResourceID(),
		Type:       "update_run",
		Content:    RunResult{Run: run, Err: err},
	})
}

func (tp tracePoller) requeue(job PollingRequest) {
	go func() {
		pp := *tp.ppGetter.GetDefault(job.ctx).Periodic
		fmt.Printf("[TracePoller] Requeuing Test Run %d. Current iteration: %d\n", job.run.ID, job.count)
		time.Sleep(pp.RetryDelayDuration())
		tp.enqueueJob(job)
	}()
}
