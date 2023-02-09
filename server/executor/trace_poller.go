package executor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

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

func NewTracePoller(
	pe PollerExecutor,
	updater RunUpdater,
	assertionRunner AssertionRunner,
	retryDelay time.Duration,
	maxWaitTimeForTrace time.Duration,
	subscriptionManager *subscription.Manager,
) PersistentTracePoller {
	maxTracePollRetry := int(math.Ceil(float64(maxWaitTimeForTrace) / float64(retryDelay)))
	return tracePoller{
		updater:             updater,
		pollerExecutor:      pe,
		maxWaitTimeForTrace: maxWaitTimeForTrace,
		maxTracePollRetry:   maxTracePollRetry,
		retryDelay:          retryDelay,
		executeQueue:        make(chan PollingRequest, 5),
		exit:                make(chan bool, 1),
		assertionRunner:     assertionRunner,
		subscriptionManager: subscriptionManager,
	}
}

type tracePoller struct {
	updater             RunUpdater
	pollerExecutor      PollerExecutor
	maxWaitTimeForTrace time.Duration
	assertionRunner     AssertionRunner

	retryDelay        time.Duration
	maxTracePollRetry int

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
					log.Printf("[TracePoller] Test %s Run %d: recieved job\n", job.test.ID, job.run.ID)
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

	job := PollingRequest{
		ctx:  ctx,
		test: test,
		run:  run,
	}

	tp.enqueueJob(job)
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

	log.Printf("[TracePoller] Test %s Run %d: Done polling. completed polling after %d times, number of spans %d\n", job.test.ID, job.run.ID, job.count, len(run.Trace.Flat))

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
	if errors.Is(err, connection.ErrTraceNotFound) {
		if time.Since(run.ServiceTriggeredAt) < tp.maxWaitTimeForTrace {
			tp.requeue(job)
			return
		}
		err = fmt.Errorf("timed out waiting for traces after %s", tp.maxWaitTimeForTrace.String())
		fmt.Println("timedout", err)
	} else {
		err = fmt.Errorf("cannot fetch trace: %w", err)
		fmt.Println("other", err)
	}

	tp.handleDBError(tp.updater.Update(job.ctx, run.Failed(err)))

	tp.subscriptionManager.PublishUpdate(subscription.Message{
		ResourceID: run.TransactionStepResourceID(),
		Type:       "update_run",
		Content:    RunResult{Run: run, Err: err},
	})

}

func (tp tracePoller) requeue(job PollingRequest) {
	go func() {
		fmt.Printf("requeuing result %d for %d time\n", job.run.ID, job.count)
		time.Sleep(tp.retryDelay)
		tp.enqueueJob(job)
	}()
}
