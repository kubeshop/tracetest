package executor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/subscription"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type AssertionRequest struct {
	carrier propagation.MapCarrier
	Test    model.Test
	Run     model.Run
}

func (r AssertionRequest) Context() context.Context {
	ctx := context.Background()
	return otel.GetTextMapPropagator().Extract(ctx, r.carrier)
}

type AssertionRunner interface {
	RunAssertions(ctx context.Context, request AssertionRequest)
	WorkerPool
}

type defaultAssertionRunner struct {
	updater             RunUpdater
	assertionExecutor   AssertionExecutor
	outputsProcessor    OutputsProcessorFn
	inputChannel        chan AssertionRequest
	exitChannel         chan bool
	subscriptionManager *subscription.Manager
	eventEmitter        EventEmitter
}

var _ WorkerPool = &defaultAssertionRunner{}
var _ AssertionRunner = &defaultAssertionRunner{}

func NewAssertionRunner(
	updater RunUpdater,
	assertionExecutor AssertionExecutor,
	op OutputsProcessorFn,
	subscriptionManager *subscription.Manager,
	eventEmitter EventEmitter,
) AssertionRunner {
	return &defaultAssertionRunner{
		outputsProcessor:    op,
		updater:             updater,
		assertionExecutor:   assertionExecutor,
		inputChannel:        make(chan AssertionRequest, 1),
		subscriptionManager: subscriptionManager,
		eventEmitter:        eventEmitter,
	}
}

func (e *defaultAssertionRunner) Start(workers int) {
	e.exitChannel = make(chan bool, workers)

	for i := 0; i < workers; i++ {
		go e.startWorker()
	}
}

func (e *defaultAssertionRunner) Stop() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			e.exitChannel <- true
			return
		}
	}
}

func (e *defaultAssertionRunner) startWorker() {
	for {
		select {
		case <-e.exitChannel:
			fmt.Println("Exiting assertion executor worker")
			return
		case request := <-e.inputChannel:
			ctx := request.Context()
			run, err := e.runAssertionsAndUpdateResult(ctx, request)

			log.Printf("[AssertionRunner] Test %s Run %d: update channel start\n", request.Test.ID, request.Run.ID)
			e.subscriptionManager.PublishUpdate(subscription.Message{
				ResourceID: run.TransactionStepResourceID(),
				Type:       "run_update",
				Content:    RunResult{Run: run, Err: err},
			})
			log.Printf("[AssertionRunner] Test %s Run %d: update channel complete\n", request.Test.ID, request.Run.ID)

			if err != nil {
				log.Printf("[AssertionRunner] Test %s Run %d: error with runAssertionsAndUpdateResult: %s\n", request.Test.ID, request.Run.ID, err.Error())
			}
		}
	}
}

func (e *defaultAssertionRunner) runAssertionsAndUpdateResult(ctx context.Context, request AssertionRequest) (model.Run, error) {
	log.Printf("[AssertionRunner] Test %s Run %d: Starting\n", request.Test.ID, request.Run.ID)

	err := e.eventEmitter.Emit(ctx, events.TestSpecsRunStart(request.Test.ID, request.Run.ID))
	if err != nil {
		log.Printf("[AssertionRunner] Test %s Run %d: fail to emit TestSpecsRunStart event: %s\n", request.Test.ID, request.Run.ID, err.Error())
	}

	run, err := e.executeAssertions(ctx, request)
	if err != nil {
		log.Printf("[AssertionRunner] Test %s Run %d: error executing assertions: %s\n", request.Test.ID, request.Run.ID, err.Error())

		anotherErr := e.eventEmitter.Emit(ctx, events.TestSpecsRunError(request.Test.ID, request.Run.ID, err))
		if anotherErr != nil {
			log.Printf("[AssertionRunner] Test %s Run %d: fail to emit TestSpecsRunError event: %s\n", request.Test.ID, request.Run.ID, anotherErr.Error())
		}

		return model.Run{}, e.updater.Update(ctx, run.Failed(err))
	}
	log.Printf("[AssertionRunner] Test %s Run %d: Success. pass: %d, fail: %d\n", request.Test.ID, request.Run.ID, run.Pass, run.Fail)

	err = e.updater.Update(ctx, run)
	if err != nil {
		log.Printf("[AssertionRunner] Test %s Run %d: error updating run: %s\n", request.Test.ID, request.Run.ID, err.Error())

		anotherErr := e.eventEmitter.Emit(ctx, events.TestSpecsRunError(request.Test.ID, request.Run.ID, err))
		if anotherErr != nil {
			log.Printf("[AssertionRunner] Test %s Run %d: fail to emit TestSpecsRunError event: %s\n", request.Test.ID, request.Run.ID, anotherErr.Error())
		}

		return model.Run{}, fmt.Errorf("could not save result on database: %w", err)
	}

	err = e.eventEmitter.Emit(ctx, events.TestSpecsRunSuccess(request.Test.ID, request.Run.ID))
	if err != nil {
		log.Printf("[AssertionRunner] Test %s Run %d: fail to emit TestSpecsRunSuccess event: %s\n", request.Test.ID, request.Run.ID, err.Error())
	}

	return run, nil
}

func (e *defaultAssertionRunner) executeAssertions(ctx context.Context, req AssertionRequest) (model.Run, error) {
	run := req.Run
	if run.Trace == nil {
		return model.Run{}, fmt.Errorf("trace not available")
	}

	ds := []expression.DataStore{expression.EnvironmentDataStore{
		Values: req.Run.Environment.Values,
	}}

	outputs, err := e.outputsProcessor(ctx, req.Test.Outputs, *run.Trace, ds)
	if err != nil {
		return model.Run{}, fmt.Errorf("cannot process outputs: %w", err)
	}
	e.validateOutputResolution(ctx, req, outputs)

	newEnvironment := createEnvironment(req.Run.Environment, outputs)

	ds = []expression.DataStore{expression.EnvironmentDataStore{Values: newEnvironment.Values}}

	assertionResult, allPassed := e.assertionExecutor.Assert(ctx, req.Test.Specs, *run.Trace, ds)

	run = run.SuccessfullyAsserted(
		outputs,
		newEnvironment,
		assertionResult,
		allPassed,
	)

	return run, nil
}

func createEnvironment(environment model.Environment, outputs model.OrderedMap[string, model.RunOutput]) model.Environment {
	outputVariables := make([]model.EnvironmentValue, 0)
	outputs.ForEach(func(key string, val model.RunOutput) error {
		outputVariables = append(outputVariables, model.EnvironmentValue{
			Key:   val.Name,
			Value: val.Value,
		})

		return nil
	})

	outputEnv := model.Environment{Values: outputVariables}

	return environment.Merge(outputEnv)
}

func (e *defaultAssertionRunner) RunAssertions(ctx context.Context, request AssertionRequest) {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	request.carrier = carrier

	e.inputChannel <- request
}

func (e *defaultAssertionRunner) validateOutputResolution(ctx context.Context, request AssertionRequest, outputs model.OrderedMap[string, model.RunOutput]) {
	err := outputs.ForEach(func(outputName string, outputModel model.RunOutput) error {
		if outputModel.Resolved {
			return nil
		}

		anotherErr := e.eventEmitter.Emit(ctx, events.TestOutputGenerationWarning(request.Test.ID, request.Run.ID, outputName))
		if anotherErr != nil {
			log.Printf("[AssertionRunner] Test %s Run %d: fail to emit TestOutputGenerationWarning event: %s\n", request.Test.ID, request.Run.ID, anotherErr.Error())
		}

		return nil
	})

	if err != nil {
		log.Printf("[AssertionRunner] Test %s Run %d: fail to validate outputs: %s\n", request.Test.ID, request.Run.ID, err.Error())
	}
}
