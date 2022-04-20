package openapi

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Runner interface {
	Run(Test) TestRunResult
}

type PersistentRunner interface {
	Runner
	WorkerPool
}

type Executor interface {
	Execute(*Test, trace.TraceID, trace.SpanID) (HttpResponse, error)
}

type TestsDB interface {
	CreateResult(_ context.Context, testID string, _ *TestRunResult) error
	UpdateTest(context.Context, *Test) error
	ResultUpdater
}

func NewPersistentRunner(e Executor, resultDB TestsDB, tp TracePoller) PersistentRunner {
	return persistentRunner{
		executor:     e,
		tp:           tp,
		testsDB:      resultDB,
		idGen:        NewRandGenerator(),
		executeQueue: make(chan execReq, 5),
		exit:         make(chan bool, 1),
	}
}

type persistentRunner struct {
	executor Executor
	tp       TracePoller
	idGen    IDGenerator
	testsDB  TestsDB

	executeQueue chan execReq
	exit         chan bool
}

type execReq struct {
	ctx    context.Context
	test   Test
	result TestRunResult
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

func (r persistentRunner) Run(t Test) TestRunResult {
	// Start a new background context for the async process
	ctx := context.Background()

	result := r.newTestResult(t.TestId)
	r.handleDBError(r.testsDB.CreateResult(ctx, result.TestId, &result))

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
	r.handleDBError(r.testsDB.UpdateResult(job.ctx, &result))

	tid, _ := trace.TraceIDFromHex(result.TraceId)
	sid, _ := trace.SpanIDFromHex(result.SpanId)

	response, err := r.executor.Execute(&job.test, tid, sid)
	result = r.handleExecutionResult(result, response, err)

	if job.test.ReferenceTestRunResult.ResultId == "" {
		job.test.ReferenceTestRunResult = TestRunResult{
			TraceId: result.TraceId,
		}
		r.handleDBError(r.testsDB.UpdateTest(job.ctx, &job.test))
	}

	r.handleDBError(r.testsDB.UpdateResult(job.ctx, &result))
	if result.State == TestRunStateAwaitingTrace {
		// start a new context
		r.tp.Poll(job.ctx, result)
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
