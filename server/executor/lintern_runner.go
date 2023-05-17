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
	lintern              lintern.Lintern
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
	lintern lintern.Lintern,
	subscriptionManager *subscription.Manager,
	eventEmitter EventEmitter,
	assertionRunner AssertionRunner,
) LinternRunner {
	return &defaultLinternRunner{
		updater:             updater,
		lintern:             lintern,
		inputChannel:        make(chan LinternRequest, 1),
		subscriptionManager: subscriptionManager,
		eventEmitter:        eventEmitter,
		assertionRunner:     assertionRunner,
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
			run, err := e.runLinternAndUpdateResult(ctx, request)

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

			assertionRequest := AssertionRequest{
				Test: request.Test,
				Run:  run,
			}
			e.assertionRunner.RunAssertions(ctx, assertionRequest)
		}
	}
}

func (e *defaultLinternRunner) RunLintern(ctx context.Context, request LinternRequest) {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	request.carrier = carrier

	e.inputChannel <- request
}

func (e *defaultLinternRunner) runLinternAndUpdateResult(ctx context.Context, request LinternRequest) (model.Run, error) {
	run := request.Run
	log.Printf("[LinternRunner] Test %s Run %d: Starting\n", request.Test.ID, request.Run.ID)

	err := e.eventEmitter.Emit(ctx, events.TraceLinternStart(request.Test.ID, request.Run.ID))
	if err != nil {
		log.Printf("[LinternRunner] Test %s Run %d: fail to emit TraceLinternStart event: %s\n", request.Test.ID, request.Run.ID, err.Error())
	}

	result, err := e.lintern.Run(ctx, *run.Trace)
	if err != nil {
		log.Printf("[LinternRunner] Test %s Run %d: error executing lintern: %s\n", request.Test.ID, request.Run.ID, err.Error())

		anotherErr := e.eventEmitter.Emit(ctx, events.TraceLinternFailed(request.Test.ID, request.Run.ID, err))
		if anotherErr != nil {
			log.Printf("[LinternRunner] Test %s Run %d: fail to emit TraceLinternFailed event: %s\n", request.Test.ID, request.Run.ID, anotherErr.Error())
		}

		run = run.LinternFailed(err)
		analytics.SendEvent("test_run_finished", "error", "", &map[string]string{
			"finalState": string(run.State),
		})

		return model.Run{}, e.updater.Update(ctx, run)
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
