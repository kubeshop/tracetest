package executor

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	triggerer "github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"go.opentelemetry.io/otel/trace"
)

func NewTriggerExecuterWorker(
	triggers *triggerer.Registry,
	updater RunUpdater,
	tracer trace.Tracer,
	eventEmitter EventEmitter,
	enabled bool,
) *triggerExecuterWorker {
	return &triggerExecuterWorker{
		triggers:     triggers,
		updater:      updater,
		tracer:       tracer,
		eventEmitter: eventEmitter,
		enabled:      enabled,
	}
}

type triggerExecuterWorker struct {
	triggers     *triggerer.Registry
	updater      RunUpdater
	tracer       trace.Tracer
	eventEmitter EventEmitter
	outputQueue  pipeline.Enqueuer[Job]
	enabled      bool
}

func (r *triggerExecuterWorker) SetOutputQueue(queue pipeline.Enqueuer[Job]) {
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
	if !r.enabled {
		r.outputQueue.Enqueue(ctx, job)
		return
	}

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
	if err != nil {
		response.Result.Error = &trigger.TriggerError{
			ConnectionError:    isConnectionError(err),
			RunningOnContainer: isServerRunningInsideContainer(),
			TargetsLocalhost:   isTargetLocalhost(job.Test.Trigger),
			ErrorMessage:       err.Error(),
		}
	}

	run.SpanID = response.SpanID
	run.TriggerResult = response.Result
	run = run.TriggerCompleted(run.TriggerResult)
	r.handleDBError(run, r.updater.Update(ctx, run))

	job.Run = run
	r.outputQueue.Enqueue(ctx, job)
}

func isConnectionError(err error) bool {
	for err != nil {
		// a dial error means we couldn't open a TCP connection (either host is not available or DNS doesn't exist)
		if strings.HasPrefix(err.Error(), "dial ") {
			return true
		}

		// it means a trigger timeout
		if errors.Is(err, context.DeadlineExceeded) {
			return true
		}

		err = errors.Unwrap(err)
	}

	return false
}

func isTargetLocalhost(t trigger.Trigger) bool {
	var endpoint string
	switch t.Type {
	case trigger.TriggerTypeHTTP:
		endpoint = t.HTTP.URL
	case trigger.TriggerTypeGRPC:
		endpoint = t.GRPC.Address
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

func isServerRunningInsideContainer() bool {
	// Check if running on Docker
	// Reference: https://paulbradley.org/indocker/
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// Check if running on k8s
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return true
	}

	return false
}
