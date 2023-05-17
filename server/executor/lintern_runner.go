package executor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/lintern"
	lintern_resource "github.com/kubeshop/tracetest/server/lintern/resource"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/subscription"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type LinternRequest struct {
	carrier propagation.MapCarrier
	Test    model.Test
	Run     model.Run
}

func (r LinternRequest) Context() context.Context {
	ctx := context.Background()
	return otel.GetTextMapPropagator().Extract(ctx, r.carrier)
}

type LinternRunner interface {
	RunLintern(ctx context.Context, request LinternRequest)
	WorkerPool
}

type LinternResourceGetter interface {
	GetDefault(ctx context.Context) lintern_resource.Lintern
}

type defaultLinternRunner struct {
	updater              RunUpdater
	inputChannel         chan LinternRequest
	exitChannel          chan bool
	subscriptionManager  *subscription.Manager
	eventEmitter         EventEmitter
	linterResourceGetter LinternResourceGetter
	assertionRunner      AssertionRunner
}

var _ WorkerPool = &defaultLinternRunner{}
var _ LinternRunner = &defaultLinternRunner{}

func NewLinternRunner(
	updater RunUpdater,
	subscriptionManager *subscription.Manager,
	eventEmitter EventEmitter,
	assertionRunner AssertionRunner,
	linterResourceGetter LinternResourceGetter,
) LinternRunner {
	return &defaultLinternRunner{
		updater:              updater,
		inputChannel:         make(chan LinternRequest, 1),
		subscriptionManager:  subscriptionManager,
		eventEmitter:         eventEmitter,
		assertionRunner:      assertionRunner,
		linterResourceGetter: linterResourceGetter,
	}
}

func (e *defaultLinternRunner) Start(workers int) {
	e.exitChannel = make(chan bool, workers)

	for i := 0; i < workers; i++ {
		go e.startWorker()
	}
}

func (e *defaultLinternRunner) Stop() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			e.exitChannel <- true
			return
		}
	}
}

func (e *defaultLinternRunner) startWorker() {
	for {
		select {
		case <-e.exitChannel:
			fmt.Println("Exiting lintern executor worker")
			return
		case request := <-e.inputChannel:
			ctx := request.Context()
			lintResource := e.linterResourceGetter.GetDefault(ctx)
			lintern := lintern.NewLintern(lintResource, lintern.AvailablePlugins...)

			shouldSkip, reason := lintern.ShouldSkip()
			if shouldSkip {
				log.Printf("[LinternRunner] Skipping TraceLintern. Reason %s\n", reason)
				err := e.eventEmitter.Emit(ctx, events.TraceLinternSkip(request.Test.ID, request.Run.ID, reason))
				if err != nil {
					log.Printf("[LinternRunner] Test %s Run %d: fail to emit TraceLinternSkip event: %s\n", request.Test.ID, request.Run.ID, err.Error())
				}

				e.onFinish(ctx, request, request.Run)
				return
			}

			err := lintern.IsValid()
			if err != nil {
				e.onError(ctx, request, request.Run, err)
				return
			}

			run, err := e.onRun(ctx, request, lintern, lintResource)
			log.Printf("[LinternRunner] Test %s Run %d: update channel start\n", request.Test.ID, request.Run.ID)
			e.subscriptionManager.PublishUpdate(subscription.Message{
				ResourceID: run.TransactionStepResourceID(),
				Type:       "run_update",
				Content:    RunResult{Run: run, Err: err},
			})
			log.Printf("[LinternRunner] Test %s Run %d: update channel complete\n", request.Test.ID, request.Run.ID)

			if err != nil {
				log.Printf("[LinternRunner] Test %s Run %d: error with runLinternAndUpdateResult: %s\n", request.Test.ID, request.Run.ID, err.Error())
				return
			}

			err = lintResource.ValidateResult(run.Lintern)
			if err != nil {
				e.onError(ctx, request, run, err)
				return
			}

			e.onFinish(ctx, request, run)
		}
	}
}

func (e *defaultLinternRunner) RunLintern(ctx context.Context, request LinternRequest) {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	request.carrier = carrier

	e.inputChannel <- request
}

func (e *defaultLinternRunner) onFinish(ctx context.Context, request LinternRequest, run model.Run) {
	assertionRequest := AssertionRequest{
		Test: request.Test,
		Run:  run,
	}
	e.assertionRunner.RunAssertions(ctx, assertionRequest)
}

func (e *defaultLinternRunner) onRun(ctx context.Context, request LinternRequest, lintern lintern.Lintern, linternResource lintern_resource.Lintern) (model.Run, error) {
	run := request.Run
	log.Printf("[LinternRunner] Test %s Run %d: Starting\n", request.Test.ID, request.Run.ID)

	err := e.eventEmitter.Emit(ctx, events.TraceLinternStart(request.Test.ID, request.Run.ID))
	if err != nil {
		log.Printf("[LinternRunner] Test %s Run %d: fail to emit TraceLinternStart event: %s\n", request.Test.ID, request.Run.ID, err.Error())
	}

	result, err := lintern.Run(ctx, *run.Trace)
	if err != nil {
		return e.onError(ctx, request, run, err)
	}

	run = run.SuccessfulLinternExecution(result)
	err = e.updater.Update(ctx, run)
	if err != nil {
		log.Printf("[LinternRunner] Test %s Run %d: error updating run: %s\n", request.Test.ID, request.Run.ID, err.Error())
		return model.Run{}, fmt.Errorf("could not save result on database: %w", err)
	}

	err = e.eventEmitter.Emit(ctx, events.TraceLinternSuccess(request.Test.ID, request.Run.ID))
	if err != nil {
		log.Printf("[LinternRunner] Test %s Run %d: fail to emit TraceLinternSuccess event: %s\n", request.Test.ID, request.Run.ID, err.Error())
	}

	return run, nil
}

func (e *defaultLinternRunner) onError(ctx context.Context, request LinternRequest, run model.Run, err error) (model.Run, error) {
	log.Printf("[LinternRunner] Test %s Run %d: Linter failed. Reason: %s\n", request.Test.ID, request.Run.ID, err.Error())
	anotherErr := e.eventEmitter.Emit(ctx, events.TraceLinternError(request.Test.ID, request.Run.ID, err))
	if anotherErr != nil {
		log.Printf("[LinternRunner] Test %s Run %d: fail to emit TraceLinternError event: %s\n", request.Test.ID, request.Run.ID, anotherErr.Error())
	}

	analytics.SendEvent("test_run_finished", "error", "", &map[string]string{
		"finalState": string(run.State),
	})

	run = run.LinternError(err)
	return run, e.updater.Update(ctx, run)
}
