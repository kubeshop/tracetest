package executor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/linter"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/subscription"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type LinterRequest struct {
	carrier propagation.MapCarrier
	Test    model.Test
	Run     model.Run
}

func (r LinterRequest) Context() context.Context {
	ctx := context.Background()
	return otel.GetTextMapPropagator().Extract(ctx, r.carrier)
}

type LinterRunner interface {
	RunLinter(ctx context.Context, request LinterRequest)
	WorkerPool
}

type AnalyzerGetter interface {
	GetDefault(ctx context.Context) analyzer.Linter
}

type defaultlinterRunner struct {
	updater             RunUpdater
	inputChannel        chan LinterRequest
	exitChannel         chan bool
	subscriptionManager *subscription.Manager
	eventEmitter        EventEmitter
	analyzerGetter      AnalyzerGetter
	assertionRunner     AssertionRunner
}

var _ WorkerPool = &defaultlinterRunner{}
var _ LinterRunner = &defaultlinterRunner{}

func NewlinterRunner(
	updater RunUpdater,
	subscriptionManager *subscription.Manager,
	eventEmitter EventEmitter,
	assertionRunner AssertionRunner,
	analyzerGetter AnalyzerGetter,
) LinterRunner {
	return &defaultlinterRunner{
		updater:             updater,
		inputChannel:        make(chan LinterRequest, 1),
		subscriptionManager: subscriptionManager,
		eventEmitter:        eventEmitter,
		assertionRunner:     assertionRunner,
		analyzerGetter:      analyzerGetter,
	}
}

func (e *defaultlinterRunner) Start(workers int) {
	e.exitChannel = make(chan bool, workers)

	for i := 0; i < workers; i++ {
		go e.startWorker()
	}
}

func (e *defaultlinterRunner) Stop() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			e.exitChannel <- true
			return
		}
	}
}

func (e *defaultlinterRunner) startWorker() {
	for {
		select {
		case <-e.exitChannel:
			fmt.Println("Exiting linter executor worker")
			return
		case request := <-e.inputChannel:
			e.onRequest(request)
		}
	}
}

func (e *defaultlinterRunner) onRequest(request LinterRequest) {
	ctx := request.Context()
	lintResource := e.analyzerGetter.GetDefault(ctx)
	linter := linter.NewLinter(lintResource, linter.AvailablePlugins...)

	shouldSkip, reason := linter.ShouldSkip()
	if shouldSkip {
		log.Printf("[linterRunner] Skipping Tracelinter. Reason %s\n", reason)
		err := e.eventEmitter.Emit(ctx, events.TraceLinterSkip(request.Test.ID, request.Run.ID, reason))
		if err != nil {
			log.Printf("[linterRunner] Test %s Run %d: fail to emit TracelinterSkip event: %s\n", request.Test.ID, request.Run.ID, err.Error())
		}

		e.onFinish(ctx, request, request.Run)
		return
	}

	err := linter.IsValid()
	if err != nil {
		e.onError(ctx, request, request.Run, err)
		return
	}

	run, err := e.onRun(ctx, request, linter, lintResource)
	log.Printf("[linterRunner] Test %s Run %d: update channel start\n", request.Test.ID, request.Run.ID)
	e.subscriptionManager.PublishUpdate(subscription.Message{
		ResourceID: run.TransactionStepResourceID(),
		Type:       "run_update",
		Content:    RunResult{Run: run, Err: err},
	})
	log.Printf("[linterRunner] Test %s Run %d: update channel complete\n", request.Test.ID, request.Run.ID)

	if err != nil {
		log.Printf("[linterRunner] Test %s Run %d: error with runlinterRunLinterAndUpdateResult: %s\n", request.Test.ID, request.Run.ID, err.Error())
		return
	}

	e.onFinish(ctx, request, run)
}

func (e *defaultlinterRunner) RunLinter(ctx context.Context, request LinterRequest) {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	request.carrier = carrier

	e.inputChannel <- request
}

func (e *defaultlinterRunner) onFinish(ctx context.Context, request LinterRequest, run model.Run) {
	assertionRequest := AssertionRequest{
		Test: request.Test,
		Run:  run,
	}
	e.assertionRunner.RunAssertions(ctx, assertionRequest)
}

func (e *defaultlinterRunner) onRun(ctx context.Context, request LinterRequest, linter linter.Linter, analyzer analyzer.Linter) (model.Run, error) {
	run := request.Run
	log.Printf("[linterRunner] Test %s Run %d: Starting\n", request.Test.ID, request.Run.ID)

	err := e.eventEmitter.Emit(ctx, events.TraceLinterStart(request.Test.ID, request.Run.ID))
	if err != nil {
		log.Printf("[linterRunner] Test %s Run %d: fail to emit TracelinterStart event: %s\n", request.Test.ID, request.Run.ID, err.Error())
	}

	result, err := linter.Run(ctx, *run.Trace)
	if err != nil {
		return e.onError(ctx, request, run, err)
	}

	result.MinimumScore = analyzer.MinimumScore
	run = run.SuccessfulLinterExecution(result)
	err = e.updater.Update(ctx, run)
	if err != nil {
		log.Printf("[linterRunner] Test %s Run %d: error updating run: %s\n", request.Test.ID, request.Run.ID, err.Error())
		return model.Run{}, fmt.Errorf("could not save result on database: %w", err)
	}

	err = e.eventEmitter.Emit(ctx, events.TraceLinterSuccess(request.Test.ID, request.Run.ID))
	if err != nil {
		log.Printf("[linterRunner] Test %s Run %d: fail to emit TracelinterSuccess event: %s\n", request.Test.ID, request.Run.ID, err.Error())
	}

	return run, nil
}

func (e *defaultlinterRunner) onError(ctx context.Context, request LinterRequest, run model.Run, err error) (model.Run, error) {
	log.Printf("[linterRunner] Test %s Run %d: Linter failed. Reason: %s\n", request.Test.ID, request.Run.ID, err.Error())
	anotherErr := e.eventEmitter.Emit(ctx, events.TraceLinterError(request.Test.ID, request.Run.ID, err))
	if anotherErr != nil {
		log.Printf("[linterRunner] Test %s Run %d: fail to emit TracelinterError event: %s\n", request.Test.ID, request.Run.ID, anotherErr.Error())
	}

	analytics.SendEvent("test_run_finished", "error", "", &map[string]string{
		"finalState": string(run.State),
	})

	run = run.LinterError(err)
	return run, e.updater.Update(ctx, run)
}
