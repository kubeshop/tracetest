package executor

import (
	"context"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/testsuite"
	"github.com/kubeshop/tracetest/server/variableset"
)

type testSuiteRunRepository interface {
	transactionUpdater
	CreateRun(context.Context, testsuite.TestSuiteRun) (testsuite.TestSuiteRun, error)
}

type testRunner interface {
	Run(context.Context, test.Test, test.RunMetadata, variableset.VariableSet, *[]testrunner.RequiredGate) test.Run
}

func NewTestSuiteRunner(
	testRunner testRunner,
	transactionRuns testSuiteRunRepository,
	subscriptionManager *subscription.Manager,
) *persistentTransactionRunner {
	updater := (CompositeTransactionUpdater{}).
		Add(NewDBTranasctionUpdater(transactionRuns)).
		Add(NewSubscriptionTransactionUpdater(subscriptionManager))

	return &persistentTransactionRunner{
		testRunner:          testRunner,
		transactionRuns:     transactionRuns,
		updater:             updater,
		subscriptionManager: subscriptionManager,
	}
}

type persistentTransactionRunner struct {
	testRunner          testRunner
	transactionRuns     testSuiteRunRepository
	updater             TestSuiteRunUpdater
	subscriptionManager *subscription.Manager
}

func (r *persistentTransactionRunner) SetOutputQueue(_ pipeline.Enqueuer[Job]) {
	// this is a no-op, as transaction runner does not need to enqueue anything
}

func (r persistentTransactionRunner) ProcessItem(ctx context.Context, job Job) {
	tran := job.TestSuite
	run := job.TestSuiteRun

	run.State = testsuite.TestSuiteStateExecuting
	err := r.updater.Update(ctx, run)
	if err != nil {
		log.Printf("[TransactionRunner] could not update transaction run: %s", err.Error())
		return
	}

	log.Printf("[TransactionRunner] running transaction %s with %d steps", run.TestSuiteID, len(tran.Steps))
	for step, test := range tran.Steps {
		run, err = r.runTransactionStep(ctx, run, step, test)
		if err != nil {
			log.Printf("[TransactionRunner] could not execute step %d of transaction %s: %s", step, run.TestSuiteID, err.Error())
			return
		}

		if run.State == testsuite.TestSuiteStateFailed {
			break
		}

		run.VariableSet = mergeOutputsIntoEnv(run.VariableSet, run.Steps[step].Outputs)
		err = r.transactionRuns.UpdateRun(ctx, run)
		if err != nil {
			log.Printf("[TransactionRunner] could not update transaction step: %s", err.Error())
			return
		}
	}

	if run.State != testsuite.TestSuiteStateFailed {
		run.State = testsuite.TestSuiteStateFinished
	}

	err = r.updater.Update(ctx, run)
	if err != nil {
		log.Printf("[TransactionRunner] could not update transaction run: %s", err.Error())
		return
	}
}

func (r persistentTransactionRunner) runTransactionStep(ctx context.Context, tr testsuite.TestSuiteRun, step int, testObj test.Test) (testsuite.TestSuiteRun, error) {
	testRun := r.testRunner.Run(ctx, testObj, tr.RunMetadata(step), tr.VariableSet, tr.RequiredGates)
	tr, err := r.updateStepRun(ctx, tr, step, testRun)
	if err != nil {
		return testsuite.TestSuiteRun{}, fmt.Errorf("could not update transaction run: %w", err)
	}

	done := make(chan bool)
	// listen for updates and propagate them as if they were transaction updates
	r.subscriptionManager.Subscribe(testRun.ResourceID(), subscription.NewSubscriberFunction(
		func(m subscription.Message) error {
			testRun := m.Content.(test.Run)
			if testRun.LastError != nil {
				tr.State = testsuite.TestSuiteStateFailed
				tr.LastError = testRun.LastError
			}

			tr, err = r.updateStepRun(ctx, tr, step, testRun)
			if err != nil {
				done <- true
				return err
			}

			r.subscriptionManager.PublishUpdate(subscription.Message{
				ResourceID: tr.ResourceID(),
				Type:       "result_update",
				Content:    tr,
			})

			if testRun.State.IsFinal() {
				done <- true
			}

			return nil
		}),
	)
	// TODO: this will block indefinitely. we need to set a timeout or something
	<-done

	return tr, err
}

func (r persistentTransactionRunner) updateStepRun(ctx context.Context, tr testsuite.TestSuiteRun, step int, run test.Run) (testsuite.TestSuiteRun, error) {
	if len(tr.Steps) <= step {
		tr.Steps = append(tr.Steps, test.Run{})
	}

	tr.Steps[step] = run
	err := r.updater.Update(ctx, tr)
	if err != nil {
		return testsuite.TestSuiteRun{}, fmt.Errorf("could not update transaction run: %w", err)
	}

	return tr, nil
}

func mergeOutputsIntoEnv(variableSet variableset.VariableSet, outputs maps.Ordered[string, test.RunOutput]) variableset.VariableSet {
	newEnv := make([]variableset.VariableSetValue, 0, outputs.Len())
	outputs.ForEach(func(key string, val test.RunOutput) error {
		newEnv = append(newEnv, variableset.VariableSetValue{
			Key:   key,
			Value: val.Value,
		})

		return nil
	})

	return variableSet.Merge(variableset.VariableSet{
		Values: newEnv,
	})
}
