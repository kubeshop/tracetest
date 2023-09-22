package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"go.opentelemetry.io/otel/trace"
)

type RunResult struct {
	Run test.Run
	Err error
}

func NewTriggerResultProcessorWorker(
	tracer trace.Tracer,
	subscriptionManager *subscription.Manager,
	eventEmitter EventEmitter,
	updater RunUpdater,
) *triggerResultProcessorWorker {
	return &triggerResultProcessorWorker{
		tracer:              tracer,
		subscriptionManager: subscriptionManager,
		eventEmitter:        eventEmitter,
		updater:             updater,
	}
}

type triggerResultProcessorWorker struct {
	tracer              trace.Tracer
	subscriptionManager *subscription.Manager
	eventEmitter        EventEmitter
	outputQueue         pipeline.Enqueuer[Job]
	updater             RunUpdater
}

func (r *triggerResultProcessorWorker) SetOutputQueue(queue pipeline.Enqueuer[Job]) {
	r.outputQueue = queue
}

func (r triggerResultProcessorWorker) handleDBError(run test.Run, err error) {
	if err != nil {
		fmt.Printf("test %s run #%d trigger DB error: %s\n", run.TestID, run.ID, err.Error())
	}
}

func (r triggerResultProcessorWorker) handleError(run test.Run, err error) {
	if err != nil {
		fmt.Printf("test %s run #%d trigger DB error: %s\n", run.TestID, run.ID, err.Error())
	}
}

func (r triggerResultProcessorWorker) ProcessItem(ctx context.Context, job Job) {
	ctx, pollingSpan := r.tracer.Start(ctx, "Start processing trigger response")
	defer pollingSpan.End()

	job.Run = r.handleExecutionResult(job.Run, ctx)
	triggerResult := job.Run.TriggerResult
	if triggerResult.Error != nil {
		err := triggerResult.Error.Error()
		if triggerResult.Error.ConnectionError {
			r.emitUnreachableEndpointEvent(ctx, job, err)

			if triggerResult.Error.TargetsLocalhost && triggerResult.Error.RunningOnContainer {
				r.emitMismatchEndpointEvent(ctx, job, err)
			}
		}

		emitErr := r.eventEmitter.Emit(ctx, events.TriggerExecutionError(job.Run.TestID, job.Run.ID, err))
		if emitErr != nil {
			r.handleError(job.Run, emitErr)
		}

		fmt.Printf("test %s run #%d trigger error: %s\n", job.Run.TestID, job.Run.ID, err.Error())
		r.subscriptionManager.PublishUpdate(subscription.Message{
			ResourceID: job.Run.TransactionStepResourceID(),
			Type:       "run_update",
			Content:    RunResult{Run: job.Run, Err: err},
		})
	} else {
		err := r.eventEmitter.Emit(ctx, events.TriggerExecutionSuccess(job.Run.TestID, job.Run.ID))
		if err != nil {
			r.handleDBError(job.Run, err)
		}
	}

	job.Run.State = test.RunStateAwaitingTrace

	r.handleDBError(job.Run, r.updater.Update(ctx, job.Run))

	r.outputQueue.Enqueue(ctx, job)
}

func (r triggerResultProcessorWorker) emitUnreachableEndpointEvent(ctx context.Context, job Job, err error) {
	var event model.TestRunEvent
	switch job.Test.Trigger.Type {
	case trigger.TriggerTypeHTTP:
		event = events.TriggerHTTPUnreachableHostError(job.Run.TestID, job.Run.ID, err)
	case trigger.TriggerTypeGRPC:
		event = events.TriggergRPCUnreachableHostError(job.Run.TestID, job.Run.ID, err)
	}

	emitErr := r.eventEmitter.Emit(ctx, event)
	if emitErr != nil {
		r.handleError(job.Run, emitErr)
	}
}

func (r triggerResultProcessorWorker) emitMismatchEndpointEvent(ctx context.Context, job Job, err error) {
	emitErr := r.eventEmitter.Emit(ctx, events.TriggerDockerComposeHostMismatchError(job.Run.TestID, job.Run.ID))
	if emitErr != nil {
		r.handleError(job.Run, emitErr)
	}
}

func (r triggerResultProcessorWorker) handleExecutionResult(run test.Run, ctx context.Context) test.Run {
	run = run.TriggerCompleted(run.TriggerResult)
	if run.TriggerResult.Error != nil {
		run = run.TriggerFailed(fmt.Errorf(run.TriggerResult.Error.ErrorMessage))

		analytics.SendEvent("test_run_finished", "error", "", &map[string]string{
			"finalState": string(run.State),
			"tenant_id":  middleware.TenantIDFromContext(ctx),
		})

		return run
	}

	return run.SuccessfullyTriggered()
}
