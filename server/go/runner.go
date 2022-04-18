package openapi

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Runner interface {
	Run(context.Context, Test) (runID string)
}

type PersistentRunner interface {
	Runner
	WorkerPool
}

type Executor interface {
	Execute(*Test, trace.TraceID, trace.SpanID) (HttpResponse, error)
}

type ResultDB interface {
	CreateResult(_ context.Context, testID string, _ *TestRunResult) error
	ResultUpdater
}

func NewPersistentRunner(e Executor, resultDB ResultDB, tp TracePoller) PersistentRunner {
	return persistentRunner{
		executor:     e,
		tp:           tp,
		resultDB:     resultDB,
		idGen:        NewRandGenerator(),
		executeQueue: make(chan execReq, 5),
		exit:         make(chan bool, 1),
	}
}

type persistentRunner struct {
	executor Executor
	tp       TracePoller
	idGen    IDGenerator
	resultDB ResultDB

	executeQueue chan execReq
	exit         chan bool
}

type execReq struct {
	ctx    context.Context
	test   Test
	result TestRunResult
}

func (r persistentRunner) Start(workers int) {
	for i := 0; i < workers; i++ {
		go func() {
			fmt.Println("persistentRunner start goroutine")
			for {
				select {
				case <-r.exit:
					fmt.Println("persistentRunner exit goroutine")
					return
				case job := <-r.executeQueue:
					fmt.Printf(
						"persistentRunner job. TestID %s, TraceID %s, SpanID %s\n",
						job.result.TestId,
						job.result.TraceId,
						job.result.SpanId,
					)
					r.processExecQueue(job)
				}
			}
		}()
	}
}

func (r persistentRunner) Stop() {
	fmt.Println("persistentRunner stopping")
	r.exit <- true
}

func (r persistentRunner) Run(ctx context.Context, t Test) string {
	result := r.newTestResult(t.TestId)
	// TODO: handle error
	_ = r.resultDB.CreateResult(ctx, result.TestId, &result)

	r.executeQueue <- execReq{
		ctx:    ctx,
		test:   t,
		result: result,
	}

	return result.ResultId
}

func (r persistentRunner) processExecQueue(job execReq) {
	result := job.result
	result.State = TestRunStateExecuting
	// TODO: handle error
	_ = r.resultDB.UpdateResult(job.ctx, &result)

	tid, _ := trace.TraceIDFromHex(result.TraceId)
	sid, _ := trace.SpanIDFromHex(result.SpanId)

	response, err := r.executor.Execute(&job.test, tid, sid)
	result = r.handleExecutionResult(result, response, err)

	// TODO: handle error
	_ = r.resultDB.UpdateResult(job.ctx, &result)
	if result.State == TestRunStateAwaitingTrace {
		// start a new context
		r.tp.Poll(context.Background(), result)
	}
}

func (r persistentRunner) handleExecutionResult(result TestRunResult, resp HttpResponse, err error) TestRunResult {
	result.CompletedAt = time.Now()
	result.Response = resp
	if err != nil {
		result.State = TestRunStateFailed
		result.LastErrorState = err.Error()
	} else {
		result.State = TestRunStateAwaitingTrace
	}
	return result
}

func (r persistentRunner) newTestResult(testID string) TestRunResult {
	return TestRunResult{
		TestId:    testID,
		ResultId:  r.idGen.UUID(),
		TraceId:   r.idGen.TraceID().String(),
		SpanId:    r.idGen.SpanID().String(),
		CreatedAt: time.Now(),
		State:     TestRunStateCreated,
	}
}
