package executor

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	triggerer "github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

type RunResult struct {
	Run test.Run
	Err error
}

type Runner interface {
	Run(context.Context, test.Test, test.RunMetadata, environment.Environment, *[]testrunner.RequiredGate) test.Run
}

type PersistentRunner interface {
	Runner
	QueueItemProcessor
}

type TestRunnerGetter interface {
	GetDefault(ctx context.Context) testrunner.TestRunner
}

func NewPersistentRunner(
	triggers *triggerer.Registry,
	runs test.RunRepository,
	updater RunUpdater,
	tracer trace.Tracer,
	subscriptionManager *subscription.Manager,
	newTraceDBFn traceDBFactoryFn,
	dsRepo resourcemanager.Current[datastore.DataStore],
	eventEmitter EventEmitter,
) *persistentRunner {
	return &persistentRunner{
		triggers:            triggers,
		runs:                runs,
		updater:             updater,
		tracer:              tracer,
		newTraceDBFn:        newTraceDBFn,
		dsRepo:              dsRepo,
		subscriptionManager: subscriptionManager,
		eventEmitter:        eventEmitter,
	}
}

type persistentRunner struct {
	triggers            *triggerer.Registry
	runs                test.RunRepository
	updater             RunUpdater
	tracer              trace.Tracer
	subscriptionManager *subscription.Manager
	newTraceDBFn        traceDBFactoryFn
	dsRepo              resourcemanager.Current[datastore.DataStore]
	eventEmitter        EventEmitter
	inputQueue          Enqueuer
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

func (r persistentRunner) traceDB(ctx context.Context) (tracedb.TraceDB, error) {
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

func (r persistentRunner) ProcessItem(ctx context.Context, job Job) {
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

	ds := []expression.DataStore{expression.EnvironmentDataStore{
		Values: run.Environment.Values,
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

	err = r.eventEmitter.Emit(ctx, events.TriggerExecutionStart(job.Run.TestID, job.Run.ID))
	if err != nil {
		r.handleError(job.Run, err)
	}

	response, err := triggererObj.Trigger(ctx, resolvedTest, triggerOptions)
	run = r.handleExecutionResult(run, response, err)
	if err != nil {
		if isConnectionError(err) {
			r.emitUnreachableEndpointEvent(ctx, job, err)

			if isTargetLocalhost(job) && isServerRunningInsideContainer() {
				r.emitMismatchEndpointEvent(ctx, job, err)
			}
		}

		emitErr := r.eventEmitter.Emit(ctx, events.TriggerExecutionError(job.Run.TestID, job.Run.ID, err))
		if emitErr != nil {
			r.handleError(job.Run, emitErr)
		}

		fmt.Printf("test %s run #%d trigger error: %s\n", run.TestID, run.ID, err.Error())
		r.subscriptionManager.PublishUpdate(subscription.Message{
			ResourceID: run.TransactionStepResourceID(),
			Type:       "run_update",
			Content:    RunResult{Run: run, Err: err},
		})
	} else {
		err = r.eventEmitter.Emit(ctx, events.TriggerExecutionSuccess(job.Run.TestID, job.Run.ID))
		if err != nil {
			r.handleDBError(job.Run, err)
		}
	}

	run.SpanID = response.SpanID

	r.handleDBError(run, r.updater.Update(ctx, run))
	if run.State != test.RunStateAwaitingTrace {
		return
	}

	ctx, pollingSpan := r.tracer.Start(ctx, "Start Polling trace")
	defer pollingSpan.End()
	r.outputQueue.Enqueue(ctx, job)
}

func (r persistentRunner) handleExecutionResult(run test.Run, response triggerer.Response, err error) test.Run {
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
