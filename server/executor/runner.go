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
	triggerer "github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type RunResult struct {
	Run test.Run
	Err error
}

type Runner interface {
	Run(context.Context, test.Test, test.RunMetadata, environment.Environment) test.Run
}

type PersistentRunner interface {
	Runner
	WorkerPool
}

func NewPersistentRunner(
	triggers *triggerer.Registry,
	runs test.RunRepository,
	updater RunUpdater,
	tp TracePoller,
	tracer trace.Tracer,
	subscriptionManager *subscription.Manager,
	newTraceDBFn traceDBFactoryFn,
	dsRepo resourcemanager.Current[datastore.DataStore],
	eventEmitter EventEmitter,
	ppGetter PollingProfileGetter,
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
		eventEmitter:        eventEmitter,
		ppGetter:            ppGetter,
		executeQueue:        make(chan execReq, 5),
		exit:                make(chan bool, 1),
	}
}

type persistentRunner struct {
	triggers            *triggerer.Registry
	tp                  TracePoller
	runs                test.RunRepository
	updater             RunUpdater
	tracer              trace.Tracer
	subscriptionManager *subscription.Manager
	newTraceDBFn        traceDBFactoryFn
	dsRepo              resourcemanager.Current[datastore.DataStore]
	eventEmitter        EventEmitter
	ppGetter            PollingProfileGetter

	executeQueue chan execReq
	exit         chan bool
}

type execReq struct {
	ctx      context.Context
	test     test.Test
	run      test.Run
	Headers  propagation.MapCarrier
	executor expression.Executor
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

func (r persistentRunner) Run(ctx context.Context, testObj test.Test, metadata test.RunMetadata, environment environment.Environment) test.Run {
	ctx, cancelCtx := context.WithCancel(
		getNewCtx(ctx),
	)

	run := test.NewRun()
	run.Metadata = metadata
	run.Environment = environment
	run, err := r.runs.CreateRun(ctx, testObj, run)
	r.handleDBError(run, err)

	r.listenForStopRequests(ctx, cancelCtx, run)

	ds := []expression.DataStore{expression.EnvironmentDataStore{
		Values: environment.Values,
	}}

	executor := expression.NewExecutor(ds...)

	r.executeQueue <- execReq{
		ctx:      ctx,
		test:     testObj,
		run:      run,
		executor: executor,
	}

	return run
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

func (r persistentRunner) processExecQueue(job execReq) {
	run := job.run.Start()
	r.handleDBError(run, r.updater.Update(job.ctx, run))

	err := r.eventEmitter.Emit(job.ctx, events.TriggerCreatedInfo(job.run.TestID, job.run.ID))
	if err != nil {
		r.handleError(job.run, err)
	}

	triggererObj, err := r.triggers.Get(job.test.Trigger.Type)
	if err != nil {
		r.handleError(job.run, err)
	}

	tdb, err := r.traceDB(job.ctx)
	if err != nil {
		r.handleError(job.run, err)
	}

	traceID := tdb.GetTraceID()
	run.TraceID = traceID
	r.handleDBError(run, r.updater.Update(job.ctx, run))

	triggerOptions := &triggerer.TriggerOptions{
		TraceID:  traceID,
		Executor: job.executor,
	}

	err = r.eventEmitter.Emit(job.ctx, events.TriggerResolveStart(job.run.TestID, job.run.ID))
	if err != nil {
		r.handleError(job.run, err)
	}

	resolvedTest, err := triggererObj.Resolve(job.ctx, job.test, triggerOptions)
	if err != nil {
		emitErr := r.eventEmitter.Emit(job.ctx, events.TriggerResolveError(job.run.TestID, job.run.ID, err))
		if emitErr != nil {
			r.handleError(job.run, emitErr)
		}

		r.handleError(job.run, err)
	}

	err = r.eventEmitter.Emit(job.ctx, events.TriggerResolveSuccess(job.run.TestID, job.run.ID))
	if err != nil {
		r.handleError(job.run, err)
	}

	if job.test.Trigger.Type == trigger.TriggerTypeTraceID {
		traceIDFromParam, err := trace.TraceIDFromHex(job.test.Trigger.TraceID.ID)
		if err == nil {
			run.TraceID = traceIDFromParam
		}
	}

	err = r.eventEmitter.Emit(job.ctx, events.TriggerExecutionStart(job.run.TestID, job.run.ID))
	if err != nil {
		r.handleError(job.run, err)
	}

	response, err := triggererObj.Trigger(job.ctx, resolvedTest, triggerOptions)
	run = r.handleExecutionResult(run, response, err)
	if err != nil {
		if isConnectionError(err) {
			r.emitUnreachableEndpointEvent(job, err)

			if isTargetLocalhost(job) && isServerRunningInsideContainer() {
				r.emitMismatchEndpointEvent(job, err)
			}
		}

		emitErr := r.eventEmitter.Emit(job.ctx, events.TriggerExecutionError(job.run.TestID, job.run.ID, err))
		if emitErr != nil {
			r.handleError(job.run, emitErr)
		}

		fmt.Printf("test %s run #%d trigger error: %s\n", run.TestID, run.ID, err.Error())
		r.subscriptionManager.PublishUpdate(subscription.Message{
			ResourceID: run.TransactionStepResourceID(),
			Type:       "run_update",
			Content:    RunResult{Run: run, Err: err},
		})
	} else {
		err = r.eventEmitter.Emit(job.ctx, events.TriggerExecutionSuccess(job.run.TestID, job.run.ID))
		if err != nil {
			r.handleDBError(job.run, err)
		}
	}

	run.SpanID = response.SpanID

	r.handleDBError(run, r.updater.Update(job.ctx, run))
	if run.State == test.RunStateAwaitingTrace {
		ctx, pollingSpan := r.tracer.Start(job.ctx, "Start Polling trace")
		defer pollingSpan.End()
		r.tp.Poll(ctx, job.test, run, r.ppGetter.GetDefault(ctx))
	}
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

func (r persistentRunner) emitUnreachableEndpointEvent(job execReq, err error) {
	var event model.TestRunEvent
	switch job.test.Trigger.Type {
	case trigger.TriggerTypeHTTP:
		event = events.TriggerHTTPUnreachableHostError(job.run.TestID, job.run.ID, err)
	case trigger.TriggerTypeGRPC:
		event = events.TriggergRPCUnreachableHostError(job.run.TestID, job.run.ID, err)
	}

	emitErr := r.eventEmitter.Emit(job.ctx, event)
	if emitErr != nil {
		r.handleError(job.run, emitErr)
	}
}

func (r persistentRunner) emitMismatchEndpointEvent(job execReq, err error) {
	emitErr := r.eventEmitter.Emit(job.ctx, events.TriggerDockerComposeHostMismatchError(job.run.TestID, job.run.ID))
	if emitErr != nil {
		r.handleError(job.run, emitErr)
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

func isTargetLocalhost(job execReq) bool {
	var endpoint string
	switch job.test.Trigger.Type {
	case trigger.TriggerTypeHTTP:
		endpoint = job.test.Trigger.HTTP.URL
	case trigger.TriggerTypeGRPC:
		endpoint = job.test.Trigger.GRPC.Address
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
