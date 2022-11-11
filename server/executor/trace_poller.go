package executor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type TracePoller interface {
	Poll(context.Context, model.Test, model.Run, chan RunResult)
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
	}
}

type tracePoller struct {
	updater             RunUpdater
	pollerExecutor      PollerExecutor
	maxWaitTimeForTrace time.Duration
	assertionRunner     AssertionRunner

	retryDelay        time.Duration
	maxTracePollRetry int

	executeQueue chan PollingRequest
	exit         chan bool
}

type PollingRequest struct {
	ctx     context.Context
	test    model.Test
	run     model.Run
	channel chan RunResult
	count   int
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

func (tp tracePoller) Poll(ctx context.Context, test model.Test, run model.Run, resultChannel chan RunResult) {
	tp.enqueueJob(PollingRequest{
		ctx:     ctx,
		test:    test,
		run:     run,
		channel: resultChannel,
	})
}

func (tp tracePoller) enqueueJob(job PollingRequest) {
	tp.executeQueue <- job
}

func (tp tracePoller) processJob(job PollingRequest) {
	finished, run, err := tp.pollerExecutor.ExecuteRequest(&job)
	if err != nil {
		tp.handleTraceDBError(job, err)
		return
	}

	if !finished {
		job.count += 1
		tp.requeue(job)
		return
	}

	fmt.Printf("completed polling result %d after %d times, number of spans: %d \n", job.run.ID, job.count, len(run.Trace.Flat))

	tp.handleDBError(tp.updater.Update(job.ctx, run))
	err = tp.runAssertions(job)
	if err != nil {
		fmt.Printf("could not run assertions: %s\n", err.Error())
	}

	job.channel <- RunResult{
		Run: run,
		Err: err,
	}
}

func (tp tracePoller) runAssertions(job PollingRequest) error {
	assertionRequest := AssertionRequest{
		Test:    job.test,
		Run:     job.run,
		channel: job.channel,
	}

	tp.assertionRunner.RunAssertions(job.ctx, assertionRequest)

	return nil
}

func augmentData(trace *traces.Trace, result model.TriggerResult) *traces.Trace {
	if trace == nil {
		return trace
	}

	switch result.Type {
	case model.TriggerTypeHTTP:
		resp := result.HTTP
		jsonheaders, _ := json.Marshal(resp.Headers)
		trace.RootSpan.Attributes["tracetest.response.status"] = fmt.Sprintf("%d", resp.StatusCode)
		trace.RootSpan.Attributes["tracetest.response.body"] = resp.Body
		trace.RootSpan.Attributes["tracetest.response.headers"] = string(jsonheaders)
	case model.TriggerTypeGRPC:
		resp := result.GRPC
		jsonheaders, _ := json.Marshal(resp.Metadata)
		trace.RootSpan.Attributes["tracetest.response.status"] = fmt.Sprintf("%d", resp.StatusCode)
		trace.RootSpan.Attributes["tracetest.response.body"] = resp.Body
		trace.RootSpan.Attributes["tracetest.response.headers"] = string(jsonheaders)
	}

	return trace
}

func (tp tracePoller) handleTraceDBError(job PollingRequest, err error) {
	run := job.run
	if errors.Is(err, tracedb.ErrTraceNotFound) {
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

}

func (tp tracePoller) requeue(job PollingRequest) {
	go func() {
		fmt.Printf("requeuing result %d for %d time\n", job.run.ID, job.count)
		time.Sleep(tp.retryDelay)
		tp.enqueueJob(job)
	}()
}
