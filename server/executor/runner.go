package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/id"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/testdb"
	"go.opentelemetry.io/otel/trace"
)

type Runner interface {
	Run(openapi.Test) openapi.TestRunResult
}

type PersistentRunner interface {
	Runner
	WorkerPool
}

type Executor interface {
	Execute(*openapi.Test, trace.TraceID, trace.SpanID) (openapi.HttpResponse, error)
}

func NewPersistentRunner(
	e Executor,
	testDB testdb.TestRepository,
	resultDB testdb.ResultRepository,
	tp TracePoller,
	ar AssertionRunner,
) PersistentRunner {
	return persistentRunner{
		executor:     e,
		tp:           tp,
		testsDB:      testDB,
		resultDB:     resultDB,
		ar:           ar,
		idGen:        id.NewRandGenerator(),
		executeQueue: make(chan execReq, 5),
		exit:         make(chan bool, 1),
	}
}

type persistentRunner struct {
	executor Executor
	tp       TracePoller
	idGen    id.Generator
	ar       AssertionRunner
	testsDB  testdb.TestRepository
	resultDB testdb.ResultRepository

	executeQueue chan execReq
	exit         chan bool
}

type execReq struct {
	ctx    context.Context
	test   openapi.Test
	result openapi.TestRunResult
}

func (r persistentRunner) handleDBError(err error) {
	if err != nil {
		fmt.Printf("DB error when running test: %s\n", err.Error())
	}
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

func (r persistentRunner) Run(t openapi.Test) openapi.TestRunResult {
	// Start a new background context for the async process
	ctx := context.Background()

	result := r.newTestResult(t.TestId)
	r.handleDBError(r.resultDB.CreateResult(ctx, result.TestId, &result))

	r.executeQueue <- execReq{
		ctx:    ctx,
		test:   t,
		result: result,
	}

	return result
}

func (r persistentRunner) processExecQueue(job execReq) {
	result := job.result
	result.State = TestRunStateExecuting
	r.handleDBError(r.resultDB.UpdateResult(job.ctx, &result))

	tid, _ := trace.TraceIDFromHex(result.TraceId)
	sid, _ := trace.SpanIDFromHex(result.SpanId)

	response, err := r.executor.Execute(&job.test, tid, sid)
	result = r.handleExecutionResult(result, response, err)

	if job.test.ReferenceTestRunResult.ResultId == "" {
		job.test.ReferenceTestRunResult = openapi.TestRunResult{
			TraceId: result.TraceId,
		}
		r.handleDBError(r.testsDB.UpdateTest(job.ctx, &job.test))
	}

	r.handleDBError(r.resultDB.UpdateResult(job.ctx, &result))
	if result.State == TestRunStateAwaitingTrace {
		// start a new context
		r.tp.Poll(job.ctx, result)
	}
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
		ResultId:  r.idGen.UUID(),
		TraceId:   r.idGen.TraceID().String(),
		SpanId:    r.idGen.SpanID().String(),
		CreatedAt: time.Now(),
		State:     TestRunStateCreated,
	}
}
