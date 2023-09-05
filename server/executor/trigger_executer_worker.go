package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/analytics"
	triggerer "github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/test"
	"go.opentelemetry.io/otel/trace"
)

func NewTriggerExecuterWorker(
	triggers *triggerer.Registry,
	updater RunUpdater,
	tracer trace.Tracer,
	eventEmitter EventEmitter,
) *triggerExecuterWorker {
	return &triggerExecuterWorker{
		triggers:     triggers,
		updater:      updater,
		tracer:       tracer,
		eventEmitter: eventEmitter,
	}
}

type triggerExecuterWorker struct {
	triggers     *triggerer.Registry
	updater      RunUpdater
	tracer       trace.Tracer
	eventEmitter EventEmitter
	outputQueue  Enqueuer
}

func (r *triggerExecuterWorker) SetOutputQueue(queue Enqueuer) {
	r.outputQueue = queue
}

func (r triggerExecuterWorker) handleDBError(run test.Run, err error) {
	if err != nil {
		fmt.Printf("test %s run #%d trigger DB error: %s\n", run.TestID, run.ID, err.Error())
	}
}

func (r triggerExecuterWorker) handleError(run test.Run, err error) {
	if err != nil {
		fmt.Printf("test %s run #%d trigger DB error: %s\n", run.TestID, run.ID, err.Error())
	}
}

func (r triggerExecuterWorker) ProcessItem(ctx context.Context, job Job) {
	err := r.eventEmitter.Emit(ctx, events.TriggerExecutionStart(job.Run.TestID, job.Run.ID))
	if err != nil {
		r.handleError(job.Run, err)
	}

	triggererObj, err := r.triggers.Get(job.Test.Trigger.Type)
	if err != nil {
		r.handleError(job.Run, err)
	}

	job.Test.Trigger = job.Run.ResolvedTrigger
	run := job.Run

	response, err := triggererObj.Trigger(ctx, job.Test, &triggerer.TriggerOptions{
		TraceID: run.TraceID,
	})
	run = r.handleExecutionResult(run, response, err)
	run.SpanID = response.SpanID

	r.handleDBError(run, r.updater.Update(ctx, run))

	job.Run = run
	r.outputQueue.Enqueue(ctx, job)
}

func (r triggerExecuterWorker) handleExecutionResult(run test.Run, response triggerer.Response, err error) test.Run {
	run = run.TriggerCompleted(response.Result)
	if err != nil {
		run = run.TriggerFailed(err)

		analytics.SendEvent("test_run_finished", "error", "", &map[string]string{
			"finalState": string(run.State),
		})

		return run
	}

	return run.SuccessfullyTriggered()
}
