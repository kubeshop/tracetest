package executor

import (
	"context"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/linter"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
)

type AnalyzerGetter interface {
	GetDefault(ctx context.Context) analyzer.Linter
}

type defaultLinterRunner struct {
	updater             RunUpdater
	subscriptionManager *subscription.Manager
	eventEmitter        EventEmitter
	analyzerGetter      AnalyzerGetter
	outputQueue         pipeline.Enqueuer[Job]
}

func NewlinterRunner(
	updater RunUpdater,
	subscriptionManager *subscription.Manager,
	eventEmitter EventEmitter,
	analyzerGetter AnalyzerGetter,
) *defaultLinterRunner {
	return &defaultLinterRunner{
		updater:             updater,
		subscriptionManager: subscriptionManager,
		eventEmitter:        eventEmitter,
		analyzerGetter:      analyzerGetter,
	}
}

func (e *defaultLinterRunner) SetOutputQueue(queue pipeline.Enqueuer[Job]) {
	e.outputQueue = queue
}

func (e *defaultLinterRunner) ProcessItem(ctx context.Context, job Job) {
	lintResource := e.analyzerGetter.GetDefault(ctx)

	shouldSkip := lintResource.ShouldSkip()
	if shouldSkip {
		e.doSkip(ctx, job)
		return
	}

	lintResource, err := lintResource.WithMetadata()
	if err != nil {
		log.Printf("[linterRunner] Test %s Run %d: error with WithMetadata: %s\n", job.Test.ID, job.Run.ID, err.Error())
		e.outputQueue.Enqueue(ctx, job)
		return
	}

	// in the future, the registry should be dynamic based on user plugins
	linter := linter.NewLinter(linter.DefaultPluginRegistry)

	run, err := e.lint(ctx, job, linter, lintResource)
	log.Printf("[linterRunner] Test %s Run %d: update channel start\n", job.Test.ID, job.Run.ID)
	e.subscriptionManager.PublishUpdate(subscription.Message{
		ResourceID: run.TransactionStepResourceID(),
		Type:       "run_update",
		Content:    RunResult{Run: run, Err: err},
	})
	log.Printf("[linterRunner] Test %s Run %d: update channel complete\n", job.Test.ID, job.Run.ID)

	if err != nil {
		log.Printf("[linterRunner] Test %s Run %d: error with runlinterRunLinterAndUpdateResult: %s\n", job.Test.ID, job.Run.ID, err.Error())
		return
	}

	job.Run = run

	e.outputQueue.Enqueue(ctx, job)
}

func (e *defaultLinterRunner) doSkip(ctx context.Context, job Job) {
	log.Printf("[linterRunner] Skipping Trace Analyzer")
	err := e.eventEmitter.Emit(ctx, events.TraceLinterSkip(job.Test.ID, job.Run.ID))
	if err != nil {
		log.Printf("[linterRunner] Test %s Run %d: fail to emit TracelinterSkip event: %s\n", job.Test.ID, job.Run.ID, err.Error())
	}

	e.outputQueue.Enqueue(ctx, job)
}

func (e *defaultLinterRunner) lint(ctx context.Context, job Job, linterObj linter.Linter, analyzerObj analyzer.Linter) (test.Run, error) {
	run := job.Run
	log.Printf("[linterRunner] Test %s Run %d: Starting\n", job.Test.ID, job.Run.ID)

	err := e.eventEmitter.Emit(ctx, events.TraceLinterStart(job.Test.ID, job.Run.ID))
	if err != nil {
		log.Printf("[linterRunner] Test %s Run %d: fail to emit TracelinterStart event: %s\n", job.Test.ID, job.Run.ID, err.Error())
	}

	result, err := linterObj.Run(ctx, *run.Trace, analyzerObj)
	if err != nil {
		return e.onError(ctx, job, run, err)
	}

	result.MinimumScore = analyzerObj.MinimumScore
	run = run.SuccessfulLinterExecution(result)
	err = e.updater.Update(ctx, run)
	if err != nil {
		log.Printf("[linterRunner] Test %s Run %d: error updating run: %s\n", job.Test.ID, job.Run.ID, err.Error())
		return test.Run{}, fmt.Errorf("could not save result on database: %w", err)
	}

	err = e.eventEmitter.Emit(ctx, events.TraceLinterSuccess(job.Test.ID, job.Run.ID))
	if err != nil {
		log.Printf("[linterRunner] Test %s Run %d: fail to emit TracelinterSuccess event: %s\n", job.Test.ID, job.Run.ID, err.Error())
	}

	return run, nil
}

func (e *defaultLinterRunner) onError(ctx context.Context, job Job, run test.Run, err error) (test.Run, error) {
	log.Printf("[linterRunner] Test %s Run %d: Linter failed. Reason: %s\n", job.Test.ID, job.Run.ID, err.Error())
	anotherErr := e.eventEmitter.Emit(ctx, events.TraceLinterError(job.Test.ID, job.Run.ID, err))
	if anotherErr != nil {
		log.Printf("[linterRunner] Test %s Run %d: fail to emit TracelinterError event: %s\n", job.Test.ID, job.Run.ID, anotherErr.Error())
	}

	analytics.SendEvent("test_run_finished", "error", "", &map[string]string{
		"finalState": string(run.State),
	})

	run = run.LinterError(err)
	return run, e.updater.Update(ctx, run)
}
