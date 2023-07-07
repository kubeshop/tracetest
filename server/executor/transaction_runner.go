package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/transaction"
)

type TransactionRunner interface {
	Run(context.Context, transaction.Transaction, test.RunMetadata, environment.Environment, *testrunner.RequiredGates) transaction.TransactionRun
}

type PersistentTransactionRunner interface {
	TransactionRunner
	WorkerPool
}

type transactionRunRepository interface {
	transactionUpdater
	CreateRun(context.Context, transaction.TransactionRun) (transaction.TransactionRun, error)
}

func NewTransactionRunner(
	runner Runner,
	db test.Repository,
	transactionRuns transactionRunRepository,
	subscriptionManager *subscription.Manager,
) persistentTransactionRunner {
	updater := (CompositeTransactionUpdater{}).
		Add(NewDBTranasctionUpdater(transactionRuns)).
		Add(NewSubscriptionTransactionUpdater(subscriptionManager))

	return persistentTransactionRunner{
		testRunner:          runner,
		db:                  db,
		transactionRuns:     transactionRuns,
		updater:             updater,
		subscriptionManager: subscriptionManager,
		executionChannel:    make(chan transactionRunJob, 1),
		exit:                make(chan bool),
	}
}

type transactionRunJob struct {
	ctx         context.Context
	transaction transaction.Transaction
	run         transaction.TransactionRun
}

type persistentTransactionRunner struct {
	testRunner          Runner
	db                  test.Repository
	transactionRuns     transactionRunRepository
	updater             TransactionRunUpdater
	subscriptionManager *subscription.Manager
	executionChannel    chan transactionRunJob
	exit                chan bool
}

func (r persistentTransactionRunner) Run(ctx context.Context, transaction transaction.Transaction, metadata test.RunMetadata, environment environment.Environment, requiredGates *testrunner.RequiredGates) transaction.TransactionRun {
	run := transaction.NewRun()
	run.Metadata = metadata
	run.Environment = environment
	run.RequiredGates = requiredGates

	ctx = getNewCtx(ctx)

	run, _ = r.transactionRuns.CreateRun(ctx, run)

	r.executionChannel <- transactionRunJob{
		ctx:         ctx,
		transaction: transaction,
		run:         run,
	}

	return run
}

func (r persistentTransactionRunner) Stop() {
	r.exit <- true
}

func (r persistentTransactionRunner) Start(workers int) {
	for i := 0; i < workers; i++ {
		go func() {
			fmt.Println("PersistentTransactionRunner start goroutine")
			for {
				select {
				case <-r.exit:
					fmt.Println("PersistentTransactionRunner exit goroutine")
					return
				case job := <-r.executionChannel:
					err := r.runTransaction(job.ctx, job.transaction, job.run)
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}()
	}
}

func (r persistentTransactionRunner) runTransaction(ctx context.Context, tran transaction.Transaction, run transaction.TransactionRun) error {
	run.State = transaction.TransactionRunStateExecuting

	var err error

	for step, test := range tran.Steps {
		run, err = r.runTransactionStep(ctx, run, step, test)
		if err != nil {
			return fmt.Errorf("could not execute step %d of transaction %s: %w", step, run.TransactionID, err)
		}

		if run.State == transaction.TransactionRunStateFailed {
			break
		}

		run.Environment = mergeOutputsIntoEnv(run.Environment, run.Steps[step].Outputs)
		err = r.transactionRuns.UpdateRun(ctx, run)
		if err != nil {
			return fmt.Errorf("coult not update transaction step: %w", err)
		}
	}

	if run.State != transaction.TransactionRunStateFailed {
		run.State = transaction.TransactionRunStateFinished
	}

	return r.updater.Update(ctx, run)
}

func (r persistentTransactionRunner) runTransactionStep(ctx context.Context, tr transaction.TransactionRun, step int, testObj test.Test) (transaction.TransactionRun, error) {
	testRun := r.testRunner.Run(ctx, testObj, tr.Metadata, tr.Environment, tr.RequiredGates)
	tr, err := r.updateStepRun(ctx, tr, step, testRun)
	if err != nil {
		return transaction.TransactionRun{}, fmt.Errorf("could not update transaction run: %w", err)
	}

	done := make(chan bool)
	// listen for updates and propagate them as if they were transaction updates
	r.subscriptionManager.Subscribe(testRun.ResourceID(), subscription.NewSubscriberFunction(
		func(m subscription.Message) error {
			testRun := m.Content.(test.Run)
			if testRun.LastError != nil {
				tr.State = transaction.TransactionRunStateFailed
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

func (r persistentTransactionRunner) updateStepRun(ctx context.Context, tr transaction.TransactionRun, step int, run test.Run) (transaction.TransactionRun, error) {
	if len(tr.Steps) <= step {
		tr.Steps = append(tr.Steps, test.Run{})
	}

	tr.Steps[step] = run
	err := r.updater.Update(ctx, tr)
	if err != nil {
		return transaction.TransactionRun{}, fmt.Errorf("could not update transaction run: %w", err)
	}

	return tr, nil
}

func mergeOutputsIntoEnv(env environment.Environment, outputs maps.Ordered[string, test.RunOutput]) environment.Environment {
	newEnv := make([]environment.EnvironmentValue, 0, outputs.Len())
	outputs.ForEach(func(key string, val test.RunOutput) error {
		newEnv = append(newEnv, environment.EnvironmentValue{
			Key:   key,
			Value: val.Value,
		})

		return nil
	})

	return env.Merge(environment.Environment{
		Values: newEnv,
	})
}
