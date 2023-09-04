package executor

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/events"
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
) *persistentRunner {
	return &persistentRunner{
		tracer:              tracer,
		subscriptionManager: subscriptionManager,
		eventEmitter:        eventEmitter,
	}
}

type persistentRunner struct {
	tracer              trace.Tracer
	subscriptionManager *subscription.Manager
	eventEmitter        EventEmitter
	outputQueue         Enqueuer
}

func (r *persistentRunner) SetOutputQueue(queue Enqueuer) {
	r.outputQueue = queue
}

func (r persistentRunner) handleDBError(run test.Run, err error) {
	if err != nil {
		fmt.Printf("test %s run #%d trigger DB error: %s\n", run.TestID, run.ID, err.Error())
	}
}

func (r persistentRunner) handleError(run test.Run, err error) {
	if err != nil {
		fmt.Printf("test %s run #%d trigger DB error: %s\n", run.TestID, run.ID, err.Error())
	}
}

func (r persistentRunner) ProcessItem(ctx context.Context, job Job) {
	ctx, pollingSpan := r.tracer.Start(ctx, "Start processing trigger response")
	defer pollingSpan.End()

	triggerResult := job.Run.TriggerResult
	if triggerResult.Error != nil {
		err := triggerResult.Error.Error()
		if triggerResult.Error.ConnectionError {
			r.emitUnreachableEndpointEvent(ctx, job, err)

			if isTargetLocalhost(job) && triggerResult.Error.RunningOnContainer {
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

	r.outputQueue.Enqueue(ctx, job)
}

func (r persistentRunner) emitUnreachableEndpointEvent(ctx context.Context, job Job, err error) {
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

func (r persistentRunner) emitMismatchEndpointEvent(ctx context.Context, job Job, err error) {
	emitErr := r.eventEmitter.Emit(ctx, events.TriggerDockerComposeHostMismatchError(job.Run.TestID, job.Run.ID))
	if emitErr != nil {
		r.handleError(job.Run, emitErr)
	}
}

func isTargetLocalhost(job Job) bool {
	var endpoint string
	switch job.Test.Trigger.Type {
	case trigger.TriggerTypeHTTP:
		endpoint = job.Test.Trigger.HTTP.URL
	case trigger.TriggerTypeGRPC:
		endpoint = job.Test.Trigger.GRPC.Address
	}

	url, err := url.Parse(endpoint)
	if err != nil {
		return false
	}

	// removes port
	host := url.Host
	colonPosition := strings.Index(url.Host, ":")
	if colonPosition >= 0 {
		host = host[0:colonPosition]
	}

	return host == "localhost" || host == "127.0.0.1"
}
