package executor

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type TracePoller interface {
	Poll(context.Context, model.Test, model.Run)
}

type PersistentTracePoller interface {
	TracePoller
	WorkerPool
}

type TraceFetcher interface {
	GetTraceByID(ctx context.Context, traceID string) (*v1.TracesData, error)
}

func NewTracePoller(
	tf TraceFetcher,
	tests model.Repository,
	maxWaitTimeForTrace time.Duration,
	subscriptionManager *subscription.Manager,
	assertionRunner AssertionRunner,
) PersistentTracePoller {
	retryDelay := 1 * time.Second
	maxTracePollRetry := int(math.Ceil(float64(maxWaitTimeForTrace) / float64(retryDelay)))
	return tracePoller{
		tests:               tests,
		traceDB:             tf,
		maxWaitTimeForTrace: maxWaitTimeForTrace,
		maxTracePollRetry:   maxTracePollRetry,
		retryDelay:          retryDelay,
		executeQueue:        make(chan tracePollReq, 5),
		exit:                make(chan bool, 1),
		subscriptions:       subscriptionManager,
		assertionRunner:     assertionRunner,
	}
}

type tracePoller struct {
	tests               model.Repository
	traceDB             TraceFetcher
	maxWaitTimeForTrace time.Duration
	retryDelay          time.Duration
	maxTracePollRetry   int

	executeQueue chan tracePollReq
	exit         chan bool

	subscriptions   *subscription.Manager
	assertionRunner AssertionRunner
}

type tracePollReq struct {
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
					fmt.Printf(
						"tracePoller job. TestID %s, TraceID %s, SpanID %s\n",
						job.test.ID,
						job.run.TraceID,
						job.run.SpanID,
					)
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
	tp.enqueueJob(tracePollReq{
		ctx:  ctx,
		test: test,
		run:  run,
	})
}

func (tp tracePoller) enqueueJob(job tracePollReq) {
	tp.executeQueue <- job
}

func (tp tracePoller) processJob(job tracePollReq) {
	run := job.run

	otelTrace, err := tp.traceDB.GetTraceByID(job.ctx, run.TraceID.String())
	if err != nil {
		tp.handleTraceDBError(job, err)
		return
	}

	run.State = model.RunStateAwaitingTestResults

	trace := traces.FromOtel(otelTrace)
	trace.ID = run.TraceID
	if !tp.donePollingTraces(job, trace) {
		fmt.Println("Not done polling traces. Requeue")
		run.Trace = &trace
		job.run = run
		job.count = job.count + 1
		tp.requeue(job)
		return
	}

	run.Trace = augmentData(&trace, run.Response)
	run.State = model.RunStateAwaitingTestResults

	fmt.Printf("completed polling result %s after %d times, number of spans: %d \n", job.run.ID, job.count, len(run.Trace.Flat))

	tp.handleDBError(tp.tests.UpdateRun(job.ctx, run))

	resource := fmt.Sprintf("test/%s/run/%s", job.test.ID, run.ID)
	tp.subscriptions.PublishUpdate(resource, subscription.Message{
		Type:    "result_update",
		Content: run,
	})

	err = tp.runAssertions(job.ctx, job.test, run)
	if err != nil {
		fmt.Printf("could not run assertions: %s\n", err.Error())
	}
}

func (tp tracePoller) runAssertions(ctx context.Context, test model.Test, run model.Run) error {
	assertionRequest := AssertionRequest{
		Test: test,
		Run:  run,
	}

	tp.assertionRunner.RunAssertions(assertionRequest)

	return nil
}

func augmentData(trace *traces.Trace, resp model.HTTPResponse) *traces.Trace {
	if trace == nil {
		return trace
	}

	trace.RootSpan.Attributes["tracetest.response.status"] = fmt.Sprintf("%d", resp.StatusCode)
	trace.RootSpan.Attributes["tracetest.response.body"] = resp.Body
	trace.RootSpan.Attributes["tracetest.response.headers"] = fmt.Sprintf("%d", resp.StatusCode)

	return trace
}

func (tp tracePoller) donePollingTraces(job tracePollReq, trace traces.Trace) bool {
	// we're done if we have the same amount of spans after polling or `maxTracePollRetry` times
	if job.count == tp.maxTracePollRetry {
		return true
	}

	if job.run.Trace == nil {
		return false
	}

	if len(trace.Flat) > 0 && len(trace.Flat) == len(job.run.Trace.Flat) {
		return true
	}

	return false

}

func (tp tracePoller) handleTraceDBError(job tracePollReq, err error) {
	run := job.run
	if errors.Is(err, tracedb.ErrTraceNotFound) {
		if time.Since(run.CompletedAt) < tp.maxWaitTimeForTrace {
			tp.requeue(job)
			return
		}
		err = fmt.Errorf("timed out waiting for traces after %s", tp.maxWaitTimeForTrace.String())
		fmt.Println("timedout", err)
	} else {
		err = fmt.Errorf("cannot fetch trace: %w", err)
		fmt.Println("other", err)
	}

	run.State = model.RunStateFailed
	run.LastError = err

	tp.handleDBError(tp.tests.UpdateRun(job.ctx, run))

}

func (tp tracePoller) requeue(job tracePollReq) {
	go func() {
		fmt.Printf("requeuing result %s for %d time\n", job.run.ID, job.count)
		time.Sleep(tp.retryDelay)
		tp.enqueueJob(job)
	}()
}
