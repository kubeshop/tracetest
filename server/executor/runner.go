package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type RunResult struct {
	Run model.Run
	Err error
}

type Runner interface {
	Run(context.Context, model.Test, model.RunMetadata, model.Environment) model.Run
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
	subscriptionManager *subscription.Manager,
	newTraceDBFn traceDBFactoryFn,
	dsRepo model.DataStoreRepository,
	eventEmitter EventEmitter,
) PersistentRunner {
	return persistentRunner{
		triggers:            triggers,
		runs:                runs,
		updater:             updater,
		tp:                  tp,
		tracer:              tracer,
		newTraceDBFn:        newTraceDBFn,
		dsRepo:              dsRepo,
		subscriptionManager: subscriptionManager,
		executeQueue:        make(chan execReq, 5),
		exit:                make(chan bool, 1),
	}
}

type persistentRunner struct {
	triggers            *trigger.Registry
	tp                  TracePoller
	runs                model.RunRepository
	updater             RunUpdater
	tracer              trace.Tracer
	subscriptionManager *subscription.Manager
	newTraceDBFn        traceDBFactoryFn
	dsRepo              model.DataStoreRepository

	executeQueue chan execReq
	exit         chan bool
}

type execReq struct {
	ctx                 context.Context
	test                model.Test
	run                 model.Run
	subscriptionManager *subscription.Manager
	Headers             propagation.MapCarrier
	executor            expression.Executor
}

func (r persistentRunner) handleDBError(run model.Run, err error) {
	if err != nil {
		fmt.Printf("test %s run #%d trigger DB error: %s\n", run.TestID, run.ID, err.Error())
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
						"persistentRunner job. ID %d, testID %s, TraceID %s, SpanID %s\n",
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

func (r persistentRunner) Run(ctx context.Context, test model.Test, metadata model.RunMetadata, environment model.Environment) model.Run {
	ctx = getNewCtx(ctx)

	run := model.NewRun()
	run.Metadata = metadata
	run.Environment = environment
	run, err := r.runs.CreateRun(ctx, test, run)
	r.handleDBError(run, err)

	ds := []expression.DataStore{expression.EnvironmentDataStore{
		Values: environment.Values,
	}}

	executor := expression.NewExecutor(ds...)

	r.executeQueue <- execReq{
		ctx:      ctx,
		test:     test,
		run:      run,
		executor: executor,
	}

	return run
}

func (r persistentRunner) traceDB(ctx context.Context) (tracedb.TraceDB, error) {
	ds, err := r.dsRepo.DefaultDataStore(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get default datastore: %w", err)
	}

	tdb, err := r.newTraceDBFn(ds)
	if err != nil {
		return nil, fmt.Errorf(`cannot get tracedb from DataStore config with ID "%s": %w`, ds.ID, err)
	}

	return tdb, nil
}

func (r persistentRunner) processExecQueue(job execReq) {
	run := job.run.Start()
	r.handleDBError(run, r.updater.Update(job.ctx, run))

	triggerer, err := r.triggers.Get(job.test.ServiceUnderTest.Type)
	if err != nil {
		// TODO: actually handle the error
		panic(err)
	}

	tdb, err := r.traceDB(job.ctx)
	if err != nil {
		panic(err)
	}

	traceID := tdb.GetTraceID()
	run.TraceID = traceID
	r.handleDBError(run, r.updater.Update(job.ctx, run))

	triggerOptions := &trigger.TriggerOptions{
		TraceID:  traceID,
		Executor: job.executor,
	}

	resolvedTest, err := triggerer.Resolve(job.ctx, job.test, triggerOptions)
	if err != nil {
		panic(err)
	}

	if job.test.ServiceUnderTest.Type == model.TriggerTypeTRACEID {
		traceIDFromParam, err := trace.TraceIDFromHex(job.test.ServiceUnderTest.TRACEID.ID)
		if err == nil {
			run.TraceID = traceIDFromParam
		}
	}

	response, err := triggerer.Trigger(job.ctx, resolvedTest, triggerOptions)
	run = r.handleExecutionResult(run, response, err)
	if err != nil {
		fmt.Printf("test %s run #%d trigger error: %s\n", run.TestID, run.ID, err.Error())
		r.subscriptionManager.PublishUpdate(subscription.Message{
			ResourceID: run.TransactionStepResourceID(),
			Type:       "run_update",
			Content:    RunResult{Run: run, Err: err},
		})
	}

	run.SpanID = response.SpanID

	r.handleDBError(run, r.updater.Update(job.ctx, run))
	if run.State == model.RunStateAwaitingTrace {
		ctx, pollingSpan := r.tracer.Start(job.ctx, "Start Polling trace")
		defer pollingSpan.End()
		r.tp.Poll(ctx, job.test, run)
	}
}

func (r persistentRunner) handleExecutionResult(run model.Run, response trigger.Response, err error) model.Run {
	run = run.TriggerCompleted(response.Result)
	if err != nil {
		return run.Failed(err)
	}

	return run.SuccessfullyTriggered()
}
