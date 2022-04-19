package openapi

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/kubeshop/tracetest/server/go/tracedb"
	"go.opentelemetry.io/otel/trace"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type TracePoller interface {
	Poll(context.Context, TestRunResult)
}

type PersistentTracePoller interface {
	TracePoller
	WorkerPool
}

type ResultUpdater interface {
	UpdateResult(ctx context.Context, res *TestRunResult) error
}

type TraceFetcher interface {
	GetTraceByID(ctx context.Context, traceID string) (*v1.TracesData, error)
}

func NewTracePoller(tf TraceFetcher, ru ResultUpdater, maxWaitTimeForTrace time.Duration) PersistentTracePoller {
	retryDelay := 500 * time.Millisecond
	maxTracePollRetry := int(math.Ceil(float64(maxWaitTimeForTrace) / float64(retryDelay)))
	return tracePoller{
		traceDB:             tf,
		resultDB:            ru,
		maxWaitTimeForTrace: maxWaitTimeForTrace,
		maxTracePollRetry:   maxTracePollRetry,
		retryDelay:          retryDelay,
		executeQueue:        make(chan tracePollReq, 5),
		exit:                make(chan bool, 1),
	}
}

type tracePoller struct {
	resultDB            ResultUpdater
	traceDB             TraceFetcher
	maxWaitTimeForTrace time.Duration
	retryDelay          time.Duration
	maxTracePollRetry   int

	executeQueue chan tracePollReq
	exit         chan bool
}

type tracePollReq struct {
	ctx    context.Context
	result TestRunResult
	count  int
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
						job.result.TestId,
						job.result.TraceId,
						job.result.SpanId,
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

func (tp tracePoller) Poll(ctx context.Context, result TestRunResult) {
	tp.enqueueJob(tracePollReq{
		ctx:    ctx,
		result: result,
	})
}

func (tp tracePoller) enqueueJob(job tracePollReq) {
	tp.executeQueue <- job
}

func (tp tracePoller) processJob(job tracePollReq) {
	res := job.result
	tr, err := tp.traceDB.GetTraceByID(job.ctx, res.TraceId)

	if err != nil {
		tp.handleTraceDBError(job, err)
		return
	}

	sid, err := trace.SpanIDFromHex(res.SpanId)
	if err != nil {
		res.State = TestRunStateFailed
		res.LastErrorState = err.Error()
		fmt.Printf("DB error when polling traces: %s\n", err.Error())
		tp.handleDBError(tp.resultDB.UpdateResult(job.ctx, &res))
		return
	}

	tid, err := trace.TraceIDFromHex(res.TraceId)
	if err != nil {
		res.State = TestRunStateFailed
		res.LastErrorState = err.Error()
		fmt.Printf("DB error when polling traces: %s\n", err.Error())
		tp.handleDBError(tp.resultDB.UpdateResult(job.ctx, &res))
		return
	}

	res.State = TestRunStateAwaitingTestResults
	res.Trace = mapTrace(
		FixParent(tr, string(tid[:]), string(sid[:]), res.Response),
	)

	if !tp.donePollingTraces(job, res) {
		job.result = res
		job.count = job.count + 1
		tp.requeue(job)
		return
	}

	fmt.Printf("completed polling result %s after %d times\n", job.result.ResultId, job.count)

	tp.handleDBError(tp.resultDB.UpdateResult(job.ctx, &res))
}

func (tp tracePoller) donePollingTraces(job tracePollReq, currentResults TestRunResult) bool {
	// we're done if we have the same amount of spans after polling `maxTracePollRetry` times
	return len(currentResults.Trace.ResourceSpans) == len(job.result.Trace.ResourceSpans) && job.count == tp.maxTracePollRetry
}

func (tp tracePoller) handleTraceDBError(job tracePollReq, err error) {
	res := job.result
	if errors.Is(err, tracedb.ErrTraceNotFound) {
		if time.Since(res.CompletedAt) < tp.maxWaitTimeForTrace {
			tp.requeue(job)
			return
		}
		err = fmt.Errorf("timed out waiting for traces after %s", tp.maxWaitTimeForTrace.String())
		fmt.Println("timedout", err)
	} else {
		err = fmt.Errorf("cannot fetch trace: %w", err)
		fmt.Println("other", err)
	}

	res.State = TestRunStateFailed
	res.LastErrorState = err.Error()

	tp.handleDBError(tp.resultDB.UpdateResult(job.ctx, &res))

}

func (tp tracePoller) requeue(job tracePollReq) {
	go func() {
		fmt.Printf("requeuing result %s for %d time\n", job.result.ResultId, job.count)
		time.Sleep(tp.retryDelay)
		tp.enqueueJob(job)
	}()
}
