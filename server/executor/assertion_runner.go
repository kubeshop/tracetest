package executor

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/variableset"
)

type defaultAssertionRunner struct {
	updater             RunUpdater
	assertionExecutor   AssertionExecutor
	outputsProcessor    OutputsProcessorFn
	subscriptionManager *subscription.Manager
	eventEmitter        EventEmitter
}

func NewAssertionRunner(
	updater RunUpdater,
	assertionExecutor AssertionExecutor,
	op OutputsProcessorFn,
	subscriptionManager *subscription.Manager,
	eventEmitter EventEmitter,
) *defaultAssertionRunner {
	return &defaultAssertionRunner{
		outputsProcessor:    op,
		updater:             updater,
		assertionExecutor:   assertionExecutor,
		subscriptionManager: subscriptionManager,
		eventEmitter:        eventEmitter,
	}
}

func (e *defaultAssertionRunner) SetOutputQueue(pipeline.Enqueuer[Job]) {
	// this is a no-op, as assertion runner does not need to enqueue anything
}

func (e *defaultAssertionRunner) ProcessItem(ctx context.Context, job Job) {
	run, err := e.runAssertionsAndUpdateResult(ctx, job)

	log.Printf("[AssertionRunner] Test %s Run %d: update channel start\n", job.Test.ID, job.Run.ID)
	e.subscriptionManager.PublishUpdate(subscription.Message{
		ResourceID: run.TransactionStepResourceID(),
		Type:       "run_update",
		Content:    RunResult{Run: run, Err: err},
	})
	log.Printf("[AssertionRunner] Test %s Run %d: update channel complete\n", job.Test.ID, job.Run.ID)

	if err != nil {
		log.Printf("[AssertionRunner] Test %s Run %d: error with runAssertionsAndUpdateResult: %s\n", job.Test.ID, job.Run.ID, err.Error())
	}
}

func (e *defaultAssertionRunner) runAssertionsAndUpdateResult(ctx context.Context, job Job) (test.Run, error) {
	log.Printf("[AssertionRunner] Test %s Run %d: Starting\n", job.Test.ID, job.Run.ID)

	err := e.eventEmitter.Emit(ctx, events.TestSpecsRunStart(job.Test.ID, job.Run.ID))
	if err != nil {
		log.Printf("[AssertionRunner] Test %s Run %d: fail to emit TestSpecsRunStart event: %s\n", job.Test.ID, job.Run.ID, err.Error())
	}

	run, err := e.executeAssertions(ctx, job)
	if err != nil {
		log.Printf("[AssertionRunner] Test %s Run %d: error executing assertions: %s\n", job.Test.ID, job.Run.ID, err.Error())

		anotherErr := e.eventEmitter.Emit(ctx, events.TestSpecsRunError(job.Test.ID, job.Run.ID, err))
		if anotherErr != nil {
			log.Printf("[AssertionRunner] Test %s Run %d: fail to emit TestSpecsRunError event: %s\n", job.Test.ID, job.Run.ID, anotherErr.Error())
		}

		run = run.AssertionFailed(err)
		analytics.SendEvent("test_run_finished", "error", "", &map[string]string{
			"finalState": string(run.State),
		})

		return test.Run{}, e.updater.Update(ctx, run)
	}
	log.Printf("[AssertionRunner] Test %s Run %d: Success. pass: %d, fail: %d\n", job.Test.ID, job.Run.ID, run.Pass, run.Fail)

	err = e.updater.Update(ctx, run)
	if err != nil {
		log.Printf("[AssertionRunner] Test %s Run %d: error updating run: %s\n", job.Test.ID, job.Run.ID, err.Error())

		anotherErr := e.eventEmitter.Emit(ctx, events.TestSpecsRunPersistenceError(job.Test.ID, job.Run.ID, err))
		if anotherErr != nil {
			log.Printf("[AssertionRunner] Test %s Run %d: fail to emit TestSpecsRunPersistenceError event: %s\n", job.Test.ID, job.Run.ID, anotherErr.Error())
		}

		return test.Run{}, fmt.Errorf("could not save result on database: %w", err)
	}

	err = e.eventEmitter.Emit(ctx, events.TestSpecsRunSuccess(job.Test.ID, job.Run.ID))
	if err != nil {
		log.Printf("[AssertionRunner] Test %s Run %d: fail to emit TestSpecsRunSuccess event: %s\n", job.Test.ID, job.Run.ID, err.Error())
	}

	return run, nil
}

func (e *defaultAssertionRunner) executeAssertions(ctx context.Context, req Job) (test.Run, error) {
	run := req.Run
	if run.Trace == nil {
		return test.Run{}, fmt.Errorf("trace not available")
	}

	ds := []expression.DataStore{expression.VariableDataStore{
		Values: req.Run.VariableSet.Values,
	}}

	outputs, err := e.outputsProcessor(ctx, req.Test.Outputs, *run.Trace, ds)
	if err != nil {
		return test.Run{}, fmt.Errorf("cannot process outputs: %w", err)
	}
	e.validateOutputResolution(ctx, req, outputs)

	newVariableSet := createVariableSet(req.Run.VariableSet, outputs)

	ds = []expression.DataStore{expression.VariableDataStore{Values: newVariableSet.Values}}

	assertionResult, allPassed := e.assertionExecutor.Assert(ctx, req.Test.Specs, *run.Trace, ds)

	e.emitFailedAssertions(ctx, req, assertionResult)

	run = run.SuccessfullyAsserted(
		outputs,
		newVariableSet,
		assertionResult,
		allPassed,
	)

	analytics.SendEvent("test_run_finished", "successful", "", &map[string]string{
		"finalState": string(run.State),
	})

	return run, nil
}

func (e *defaultAssertionRunner) emitFailedAssertions(ctx context.Context, req Job, result maps.Ordered[test.SpanQuery, []test.AssertionResult]) {
	for _, assertionResults := range result.Unordered() {
		for _, assertionResult := range assertionResults {
			for _, spanAssertionResult := range assertionResult.Results {

				if errors.Is(spanAssertionResult.CompareErr, expression.ErrExpressionResolution) {
					unwrappedError := errors.Unwrap(spanAssertionResult.CompareErr)
					e.eventEmitter.Emit(ctx, events.TestSpecsAssertionWarning(
						req.Run.TestID,
						req.Run.ID,
						unwrappedError,
						spanAssertionResult.SpanIDString(),
						string(assertionResult.Assertion),
					))
				}

				if errors.Is(spanAssertionResult.CompareErr, expression.ErrInvalidSyntax) {
					e.eventEmitter.Emit(ctx, events.TestSpecsAssertionWarning(
						req.Run.TestID,
						req.Run.ID,
						spanAssertionResult.CompareErr,
						spanAssertionResult.SpanIDString(),
						string(assertionResult.Assertion),
					))
				}

			}
		}
	}
}

func createVariableSet(env variableset.VariableSet, outputs maps.Ordered[string, test.RunOutput]) variableset.VariableSet {
	outputVariables := make([]variableset.VariableSetValue, 0)
	outputs.ForEach(func(key string, val test.RunOutput) error {
		outputVariables = append(outputVariables, variableset.VariableSetValue{
			Key:   val.Name,
			Value: val.Value,
		})

		return nil
	})

	outputEnv := variableset.VariableSet{Values: outputVariables}

	return env.Merge(outputEnv)
}
func (e *defaultAssertionRunner) validateOutputResolution(ctx context.Context, job Job, outputs maps.Ordered[string, test.RunOutput]) {
	err := outputs.ForEach(func(outputName string, outputModel test.RunOutput) error {
		if outputModel.Resolved {
			return nil
		}

		anotherErr := e.eventEmitter.Emit(ctx, events.TestOutputGenerationWarning(job.Test.ID, job.Run.ID, outputModel.Error, outputName))
		if anotherErr != nil {
			log.Printf("[AssertionRunner] Test %s Run %d: fail to emit TestOutputGenerationWarning event: %s\n", job.Test.ID, job.Run.ID, anotherErr.Error())
		}

		return nil
	})

	if err != nil {
		log.Printf("[AssertionRunner] Test %s Run %d: fail to validate outputs: %s\n", job.Test.ID, job.Run.ID, err.Error())
	}
}
