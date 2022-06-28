package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type AssertionRequest struct {
	carrier propagation.MapCarrier
	Test    model.Test
	Run     model.Run
}

type AssertionRunner interface {
	RunAssertions(ctx context.Context, request AssertionRequest)
	WorkerPool
}

type defaultAssertionRunner struct {
	updater           RunUpdater
	assertionExecutor AssertionExecutor
	inputChannel      chan AssertionRequest
	exitChannel       chan bool
}

var _ WorkerPool = &defaultAssertionRunner{}
var _ AssertionRunner = &defaultAssertionRunner{}

func NewAssertionRunner(updater RunUpdater, assertionExecutor AssertionExecutor) AssertionRunner {
	return &defaultAssertionRunner{
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
			ctx := e.getCtxFromRequest(assertionRequest)
			err := e.runAssertionsAndUpdateResult(ctx, assertionRequest)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func (e *defaultAssertionRunner) getCtxFromRequest(req AssertionRequest) context.Context {
	ctx := context.Background()
	otel.GetTextMapPropagator().Extract(ctx, req.carrier)

	return ctx
}

func (e *defaultAssertionRunner) runAssertionsAndUpdateResult(ctx context.Context, request AssertionRequest) error {
	run, err := e.executeAssertions(ctx, request)
	if err != nil {
		return e.updater.Update(ctx, run.Failed(err))
	}

	err = e.updater.Update(ctx, run)
	if err != nil {
		return fmt.Errorf("could not save result on database: %w", err)
	}

	return nil
}

func (e *defaultAssertionRunner) executeAssertions(ctx context.Context, req AssertionRequest) (model.Run, error) {
	run := req.Run
	if run.Trace == nil {
		return model.Run{}, fmt.Errorf("trace not available")
	}

	run = run.SuccessfullyAsserted(
		e.assertionExecutor.Assert(ctx, req.Test.Definition, *run.Trace),
	)

	return run, nil
}

func (e *defaultAssertionRunner) RunAssertions(ctx context.Context, request AssertionRequest) {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	request.carrier = carrier

	e.inputChannel <- request
}
