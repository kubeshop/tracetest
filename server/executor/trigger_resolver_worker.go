package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/datastore"
	triggerer "github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

type currentDataStoreGetter interface {
	Current(context.Context) (datastore.DataStore, error)
}

func NewTriggerResolverWorker(
	triggers *triggerer.Registry,
	updater RunUpdater,
	tracer trace.Tracer,
	newTraceDBFn tracedb.FactoryFunc,
	dsRepo currentDataStoreGetter,
	eventEmitter EventEmitter,
) *triggerResolverWorker {
	return &triggerResolverWorker{
		triggers:     triggers,
		updater:      updater,
		tracer:       tracer,
		newTraceDBFn: newTraceDBFn,
		dsRepo:       dsRepo,
		eventEmitter: eventEmitter,
	}
}

type triggerResolverWorker struct {
	triggers     *triggerer.Registry
	updater      RunUpdater
	tracer       trace.Tracer
	newTraceDBFn tracedb.FactoryFunc
	dsRepo       currentDataStoreGetter
	eventEmitter EventEmitter
	outputQueue  Enqueuer
}

func (r *triggerResolverWorker) SetOutputQueue(queue Enqueuer) {
	r.outputQueue = queue
}

func (r triggerResolverWorker) handleDBError(run test.Run, err error) {
	if err != nil {
		fmt.Printf("test %s run #%d trigger DB error: %s\n", run.TestID, run.ID, err.Error())
	}
}

func (r triggerResolverWorker) handleError(run test.Run, err error) {
	if err != nil {
		fmt.Printf("test %s run #%d trigger DB error: %s\n", run.TestID, run.ID, err.Error())
	}
}

func (r triggerResolverWorker) traceDB(ctx context.Context) (tracedb.TraceDB, error) {
	ds, err := r.dsRepo.Current(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get default datastore: %w", err)
	}

	tdb, err := r.newTraceDBFn(ds)
	if err != nil {
		return nil, fmt.Errorf(`cannot get tracedb from DataStore config with ID "%s": %w`, ds.ID, err)
	}

	return tdb, nil
}

func (r triggerResolverWorker) ProcessItem(ctx context.Context, job Job) {
	ctx, pollingSpan := r.tracer.Start(ctx, "Resolve trigger")
	defer pollingSpan.End()

	run := job.Run.Start()
	r.handleDBError(run, r.updater.Update(ctx, run))

	err := r.eventEmitter.Emit(ctx, events.TriggerCreatedInfo(job.Run.TestID, job.Run.ID))
	if err != nil {
		r.handleError(job.Run, err)
	}

	triggererObj, err := r.triggers.Get(job.Test.Trigger.Type)
	if err != nil {
		r.handleError(job.Run, err)
	}

	tdb, err := r.traceDB(ctx)
	if err != nil {
		r.handleError(job.Run, err)
	}

	traceID := tdb.GetTraceID()
	run.TraceID = traceID
	r.handleDBError(run, r.updater.Update(ctx, run))

	ds := []expression.DataStore{expression.VariableDataStore{
		Values: run.VariableSet.Values,
	}}

	executor := expression.NewExecutor(ds...)

	triggerOptions := &triggerer.TriggerOptions{
		TraceID:  traceID,
		Executor: executor,
	}

	err = r.eventEmitter.Emit(ctx, events.TriggerResolveStart(job.Run.TestID, job.Run.ID))
	if err != nil {
		r.handleError(job.Run, err)
	}

	resolvedTest, err := triggererObj.Resolve(ctx, job.Test, triggerOptions)
	if err != nil {
		emitErr := r.eventEmitter.Emit(ctx, events.TriggerResolveError(job.Run.TestID, job.Run.ID, err))
		if emitErr != nil {
			r.handleError(job.Run, emitErr)
		}

		r.handleError(job.Run, err)
	}

	err = r.eventEmitter.Emit(ctx, events.TriggerResolveSuccess(job.Run.TestID, job.Run.ID))
	if err != nil {
		r.handleError(job.Run, err)
	}

	if job.Test.Trigger.Type == trigger.TriggerTypeTraceID {
		traceIDFromParam, err := trace.TraceIDFromHex(job.Test.Trigger.TraceID.ID)
		if err == nil {
			run.TraceID = traceIDFromParam
		}
	}

	run.ResolvedTrigger = resolvedTest.Trigger
	r.handleDBError(run, r.updater.Update(ctx, run))
	job.Run = run

	r.outputQueue.Enqueue(ctx, job)
}
