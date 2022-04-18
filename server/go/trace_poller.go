package openapi

import (
	"context"
	"errors"
	"fmt"
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
	return tracePoller{
		traceDB:             tf,
		resultDB:            ru,
		maxWaitTimeForTrace: maxWaitTimeForTrace,
		executeQueue:        make(chan tracePollReq, 5),
		exit:                make(chan bool, 1),
	}
}

type tracePoller struct {
	resultDB            ResultUpdater
	traceDB             TraceFetcher
	maxWaitTimeForTrace time.Duration

	executeQueue chan tracePollReq

	exit chan bool
}

type tracePollReq struct {
	ctx    context.Context
	result TestRunResult
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
	tp.executeQueue <- tracePollReq{
		ctx:    ctx,
		result: result,
	}
}

func (tp tracePoller) processJob(job tracePollReq) {
	res := job.result
	tr, err := tp.traceDB.GetTraceByID(job.ctx, res.TraceId)

	if err != nil {
		tp.handleError(job, err)
		return
	}

	// TODO: handle errors
	sid, _ := trace.SpanIDFromHex(res.SpanId)
	tid, _ := trace.TraceIDFromHex(res.TraceId)

	res.State = TestRunStateAwaitingTestResults
	res.Trace = mapTrace(
		FixParent(tr, string(tid[:]), string(sid[:]), res.Response),
	)

	// TODO: handle error
	_ = tp.resultDB.UpdateResult(job.ctx, &res)
}

func (tp tracePoller) handleError(job tracePollReq, err error) {
	res := job.result
	if errors.Is(err, tracedb.ErrTraceNotFound) {
		if time.Since(res.CompletedAt) < tp.maxWaitTimeForTrace {
			fmt.Println("requeue")
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

	// TODO: handle error
	_ = tp.resultDB.UpdateResult(job.ctx, &res)

}

func (tp tracePoller) requeue(job tracePollReq) {
	go func() {
		time.Sleep(500 * time.Millisecond)
		tp.Poll(job.ctx, job.result)
	}()
}
