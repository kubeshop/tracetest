package testrunner

import (
	"context"
	"fmt"
	"time"

	openapi "github.com/kubeshop/tracetest/server/go"
	"go.opentelemetry.io/otel/trace"
)

type PersistentRunner interface {
	Run(context.Context, openapi.Test)
	Start(workers int)
	Stop()
}

type Executor interface {
	Execute(*openapi.Test, trace.TraceID, trace.SpanID) (openapi.HttpResponse, error)
}

type ResultDB interface {
	CreateResult(ctx context.Context, res *openapi.TestRunResult) error
	UpdateResult(ctx context.Context, res *openapi.TestRunResult) error
}

func NewPersistentRunner(e Executor, resultDB ResultDB) PersistentRunner {
	return persistentRunner{
		executor:     e,
		resultDB:     resultDB,
		idGen:        NewRandGenerator(),
		executeQueue: make(chan execReq, 5),
		exit:         make(chan bool, 1),
	}
}

type persistentRunner struct {
	executor Executor
	idGen    IDGenerator
	resultDB ResultDB
	// traceDB tracedb.TraceDB

	executeQueue chan execReq

	exit chan bool
}

type execReq struct {
	ctx    context.Context
	test   openapi.Test
	result openapi.TestRunResult
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

func (r persistentRunner) Run(ctx context.Context, t openapi.Test) {
	result := r.newTestResult(t.TestId)
	// TODO: handle error
	_ = r.resultDB.CreateResult(ctx, &result)

	r.executeQueue <- execReq{
		ctx:    ctx,
		test:   t,
		result: result,
	}
}

func (r persistentRunner) processExecQueue(job execReq) {
	result := job.result
	result.State = TestRunStateExecuting
	// TODO: handle error
	_ = r.resultDB.UpdateResult(job.ctx, &result)

	tid, _ := trace.TraceIDFromHex(result.TraceId)
	sid, _ := trace.SpanIDFromHex(result.SpanId)

	fmt.Println("Starting test", job.test.TestId)
	response, err := r.executor.Execute(&job.test, tid, sid)
	fmt.Println("Completed test", job.test.TestId)
	result = r.handleExecutionResult(result, response, err)

	// TODO: handle error
	_ = r.resultDB.UpdateResult(job.ctx, &result)
}

func (r persistentRunner) handleExecutionResult(result openapi.TestRunResult, resp openapi.HttpResponse, err error) openapi.TestRunResult {
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

func (r persistentRunner) newTestResult(testID string) openapi.TestRunResult {
	return openapi.TestRunResult{
		TestId:    testID,
		TraceId:   r.idGen.TraceID().String(),
		SpanId:    r.idGen.SpanID().String(),
		CreatedAt: time.Now(),
		State:     TestRunStateCreated,
	}
}
