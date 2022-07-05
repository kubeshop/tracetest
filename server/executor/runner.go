package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Runner interface {
	Run(context.Context, model.Test) model.Run
}

type PersistentRunner interface {
	Runner
	WorkerPool
}

func NewPersistentRunner(
	triggers *trigger.Registry,
	runs model.RunRepository,
	updater RunUpdater,
	tp TracePoller,
	tracer trace.Tracer,
) PersistentRunner {
	return persistentRunner{
		triggers:     triggers,
		runs:         runs,
		updater:      updater,
		tp:           tp,
		tracer:       tracer,
		executeQueue: make(chan execReq, 5),
		exit:         make(chan bool, 1),
	}
}

type persistentRunner struct {
	triggers *trigger.Registry
	tp       TracePoller
	runs     model.RunRepository
	updater  RunUpdater
	tracer   trace.Tracer

	executeQueue chan execReq
	exit         chan bool
}

type execReq struct {
	ctx     context.Context
	test    model.Test
	run     model.Run
	Headers propagation.MapCarrier
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

func getNewCtx(ctx context.Context) context.Context {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	return otel.GetTextMapPropagator().Extract(context.Background(), carrier)
}

func (r persistentRunner) Run(ctx context.Context, test model.Test) model.Run {
	ctx = getNewCtx(ctx)

	run, err := r.runs.CreateRun(ctx, test, model.NewRun())
	r.handleDBError(err)

	r.executeQueue <- execReq{
		ctx:  ctx,
		test: test,
		run:  run,
	}

	return run
}

func (r persistentRunner) processExecQueue(job execReq) {
	run := job.run.Start()
	r.handleDBError(r.updater.Update(job.ctx, run))

	// TODO: hardcoded trigger type.
	trigger, err := r.triggers.Get("http")
	if err != nil {
		panic(err)
	}

	response, err := trigger.Trigger(job.ctx, job.test, job.run.TraceID, job.run.SpanID)
	// TODO: hardcoded response type.
	run = r.handleExecutionResult(run, response.Response.(model.HTTPResponse), err)

	r.handleDBError(r.updater.Update(job.ctx, run))
	if run.State == model.RunStateAwaitingTrace {
		ctx, span := r.tracer.Start(job.ctx, "Start Polling trace")
		defer span.End()
		r.tp.Poll(ctx, job.test, run)
	}
}

func (r persistentRunner) handleExecutionResult(run model.Run, resp model.HTTPResponse, err error) model.Run {
	run.Response = resp
	if err != nil {
		return run.Failed(err)
	}

	return run.SuccessfullyExecuted()
}
