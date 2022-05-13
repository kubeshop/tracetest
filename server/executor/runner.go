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
	Run(openapi.Test) openapi.TestRun
}

type PersistentRunner interface {
	Runner
	WorkerPool
}

type Executor interface {
	Execute(*openapi.Test, trace.TraceID, trace.SpanID) (openapi.HttpResponse, error)
}

type testsDB interface {
	UpdateTest(context.Context, *openapi.Test) error
	CreateRun(_ context.Context, Id string, _ *openapi.TestRun) error
	runUpdater
}

func NewPersistentRunner(
	e Executor,
	testDB testdb.TestRepository,
	resultDB testdb.ResultRepository,
	tp TracePoller,
) PersistentRunner {
	return persistentRunner{
		executor:     e,
		tp:           tp,
		testsDB:      testDB,
		resultDB:     resultDB,
		idGen:        id.NewRandGenerator(),
		executeQueue: make(chan execReq, 5),
		exit:         make(chan bool, 1),
	}
}

type persistentRunner struct {
	executor Executor
	tp       TracePoller
	idGen    id.Generator
	testsDB  testdb.TestRepository
	resultDB testdb.ResultRepository

	executeQueue chan execReq
	exit         chan bool
}

type execReq struct {
	ctx  context.Context
	test openapi.Test
	run  openapi.TestRun
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
						"persistentRunner job. Id %s, TraceID %s, SpanID %s\n",
						job.run.Id,
						job.run.TraceId,
						job.run.SpanId,
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

func (r persistentRunner) Run(t openapi.Test) openapi.TestRun {
	// Start a new background context for the async process
	ctx := context.Background()

	result := r.newTestResult(t.TestId)
	r.handleDBError(r.resultDB.CreateResult(ctx, result.TestId, &result))

	r.executeQueue <- execReq{
		ctx:  ctx,
		test: t,
		run:  result,
	}

	return result
}

func (r persistentRunner) processExecQueue(job execReq) {
	result := job.result
	result.State = TestRunStateExecuting
	r.handleDBError(r.resultDB.UpdateResult(job.ctx, &result))

	tid, _ := trace.TraceIDFromHex(run.TraceId)
	sid, _ := trace.SpanIDFromHex(run.SpanId)

	response, err := r.executor.Execute(&job.test, tid, sid)
	run = r.handleExecutionResult(run, response, err)

	if job.test.ReferenceTestRun.Id == "" {
		job.test.ReferenceTestRun = openapi.TestRun{
			TraceId: run.TraceId,
		}
		r.handleDBError(r.testsDB.UpdateTest(job.ctx, &job.test))
	}

	r.handleDBError(r.resultDB.UpdateResult(job.ctx, &result))
	if result.State == TestRunStateAwaitingTrace {
		// start a new context
		r.tp.Poll(job.ctx, run)
	}
}

func (r persistentRunner) handleExecutionResult(result openapi.TestRun, resp openapi.HttpResponse, err error) openapi.TestRun {
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

func (r persistentRunner) newTestRun(Id string) openapi.TestRun {
	return openapi.TestRun{
		Id:        Id,
		TraceId:   r.idGen.TraceID().String(),
		SpanId:    r.idGen.SpanID().String(),
		CreatedAt: time.Now(),
		State:     TestRunStateCreated,
	}
}
