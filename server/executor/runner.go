package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/trace"
)

type Runner interface {
	Run(model.Test) model.Run
}

type PersistentRunner interface {
	Runner
	WorkerPool
}

type Executor interface {
	Execute(model.Test, trace.TraceID, trace.SpanID) (model.HTTPResponse, error)
}

func NewPersistentRunner(
	e Executor,
	tests model.Repository,
	tp TracePoller,
) PersistentRunner {
	return persistentRunner{
		executor:     e,
		tests:        tests,
		tp:           tp,
		idGen:        id.NewRandGenerator(),
		executeQueue: make(chan execReq, 5),
		exit:         make(chan bool, 1),
	}
}

type persistentRunner struct {
	executor Executor
	tp       TracePoller
	idGen    id.Generator
	tests    model.Repository

	executeQueue chan execReq
	exit         chan bool
}

type execReq struct {
	ctx  context.Context
	test model.Test
	run  model.Run
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
					job.run.StartAt = time.Now()
					fmt.Printf(
						"persistentRunner job. ID %s, testID %s, TraceID %s, SpanID %s\n",
						job.run.ID,
						job.test.ID,
						job.run.TraceID,
						job.run.SpanID,
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

func (r persistentRunner) Run(test model.Test) model.Run {
	// Start a new background context for the async process
	ctx := context.Background()

	run, err := r.tests.CreateRun(ctx, test, r.newTestRun())
	r.handleDBError(err)

	r.executeQueue <- execReq{
		ctx:  ctx,
		test: test,
		run:  run,
	}

	return run
}

func (r persistentRunner) processExecQueue(job execReq) {
	run := job.run
	run.State = model.RunStateExecuting
	run.ServiceTriggeredAt = time.Now()
	r.handleDBError(r.tests.UpdateRun(job.ctx, run))

	response, err := r.executor.Execute(job.test, job.run.TraceID, job.run.SpanID)
	run = r.handleExecutionResult(run, response, err)

	if job.test.ReferenceRun == nil {
		job.test.ReferenceRun = &run
		r.handleDBError(r.tests.UpdateTestVersion(job.ctx, job.test))
	}

	r.handleDBError(r.tests.UpdateRun(job.ctx, run))
	if run.State == model.RunStateAwaitingTrace {
		// start a new context
		r.tp.Poll(job.ctx, job.test, run)
	}
}

func (r persistentRunner) handleExecutionResult(run model.Run, resp model.HTTPResponse, err error) model.Run {
	run.Response = resp
	if err != nil {
		run.State = model.RunStateFailed
		run.LastError = err
	} else {
		run.State = model.RunStateAwaitingTrace
	}
	return run
}

func (r persistentRunner) newTestRun() model.Run {
	return model.Run{
		ID:                 r.idGen.UUID(),
		TraceID:            r.idGen.TraceID(),
		SpanID:             r.idGen.SpanID(),
		State:              model.RunStateCreated,
		CreatedAt:          time.Now(),
		StartAt:            time.Time{},
		ServiceTriggeredAt: time.Time{},
		ObtainedTraceAt:    time.Time{},
		CompletedAt:        time.Time{}, // zero value
	}
}
