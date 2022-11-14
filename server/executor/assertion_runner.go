package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type AssertionRequest struct {
	carrier propagation.MapCarrier
	Test    model.Test
	Run     model.Run
	channel chan RunResult
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
	updater           RunUpdater
	assertionExecutor AssertionExecutor
	outputsProcessor  OutputsProcessorFn
	inputChannel      chan AssertionRequest
	exitChannel       chan bool
}

var _ WorkerPool = &defaultAssertionRunner{}
var _ AssertionRunner = &defaultAssertionRunner{}

func NewAssertionRunner(updater RunUpdater, assertionExecutor AssertionExecutor, op OutputsProcessorFn) AssertionRunner {
	return &defaultAssertionRunner{
		outputsProcessor:  op,
		updater:           updater,
		assertionExecutor: assertionExecutor,
		inputChannel:      make(chan AssertionRequest, 1),
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
		case assertionRequest := <-e.inputChannel:
			ctx := assertionRequest.Context()
			run, err := e.runAssertionsAndUpdateResult(ctx, assertionRequest)

			runResult := RunResult{Run: run, Err: err}
			assertionRequest.channel <- runResult

			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func (e *defaultAssertionRunner) runAssertionsAndUpdateResult(ctx context.Context, request AssertionRequest) (model.Run, error) {
	run, err := e.executeAssertions(ctx, request)
	if err != nil {
		return model.Run{}, e.updater.Update(ctx, run.Failed(err))
	}

	err = e.updater.Update(ctx, run)
	if err != nil {
		return model.Run{}, fmt.Errorf("could not save result on database: %w", err)
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

	assertionResult, allPassed := e.assertionExecutor.Assert(ctx, req.Test.Specs, *run.Trace, ds)

	run = run.SuccessfullyAsserted(
		outputs,
		assertionResult,
		allPassed,
	)

	return run, nil
}

func (e *defaultAssertionRunner) RunAssertions(ctx context.Context, request AssertionRequest) {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	request.carrier = carrier

	e.inputChannel <- request
}
