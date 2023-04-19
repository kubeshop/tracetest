package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/subscription"
)

type TransactionRunner interface {
	Run(context.Context, model.Transaction, model.RunMetadata, model.Environment) model.TransactionRun
}

type PersistentTransactionRunner interface {
	TransactionRunner
	WorkerPool
}

func NewTransactionRunner(
	runner Runner,
	db model.Repository,
	subscriptionManager *subscription.Manager,
) persistentTransactionRunner {
	updater := (CompositeTransactionUpdater{}).
		Add(NewDBTranasctionUpdater(db)).
		Add(NewSubscriptionTransactionUpdater(subscriptionManager))

	return persistentTransactionRunner{
		testRunner:          runner,
		db:                  db,
		updater:             updater,
		subscriptionManager: subscriptionManager,
		executionChannel:    make(chan transactionRunJob, 1),
		exit:                make(chan bool),
	}
}

type transactionRunJob struct {
	ctx         context.Context
	transaction model.Transaction
	run         model.TransactionRun
}

type persistentTransactionRunner struct {
	testRunner          Runner
	db                  model.Repository
	updater             TransactionRunUpdater
	subscriptionManager *subscription.Manager
	executionChannel    chan transactionRunJob
	exit                chan bool
}

func (r persistentTransactionRunner) Run(ctx context.Context, transaction model.Transaction, metadata model.RunMetadata, environment model.Environment) model.TransactionRun {
	run := transaction.NewRun()
	run.Metadata = metadata
	run.Environment = environment

	ctx = getNewCtx(ctx)

	run, _ = r.db.CreateTransactionRun(ctx, run)

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

func (r persistentTransactionRunner) runTransaction(ctx context.Context, transaction model.Transaction, run model.TransactionRun) error {
	run.State = model.TransactionRunStateExecuting

	var err error

	for step, test := range transaction.Steps {
		run, err = r.runTransactionStep(ctx, run, step, test)
		if err != nil {
			return fmt.Errorf("could not execute step %d of transaction %s: %w", step, run.TransactionID, err)
		}

		if run.State == model.TransactionRunStateFailed {
			break
		}

		run.Environment = mergeOutputsIntoEnv(run.Environment, run.Steps[step].Outputs)
		err = r.db.UpdateTransactionRun(ctx, run)
		if err != nil {
			return fmt.Errorf("coult not update transaction step: %w", err)
		}
	}

	if run.State != model.TransactionRunStateFailed {
		run.State = model.TransactionRunStateFinished
	}

	return r.updater.Update(ctx, run)
}

func (r persistentTransactionRunner) runTransactionStep(ctx context.Context, tr model.TransactionRun, step int, test model.Test) (model.TransactionRun, error) {
	testRun := r.testRunner.Run(ctx, test, tr.Metadata, tr.Environment)
	tr, err := r.updateStepRun(ctx, tr, step, testRun)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("could not update transaction run: %w", err)
	}

	done := make(chan bool)
	// listen for updates and propagate them as if they were transaction updates
	r.subscriptionManager.Subscribe(testRun.ResourceID(), subscription.NewSubscriberFunction(
		func(m subscription.Message) error {
			testRun := m.Content.(model.Run)
			if testRun.LastError != nil {
				tr.State = model.TransactionRunStateFailed
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

func (r persistentTransactionRunner) updateStepRun(ctx context.Context, tr model.TransactionRun, step int, run model.Run) (model.TransactionRun, error) {
	if len(tr.Steps) <= step {
		tr.Steps = append(tr.Steps, model.Run{})
	}

	tr.Steps[step] = run
	err := r.updater.Update(ctx, tr)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("could not update transaction run: %w", err)
	}

	return tr, nil
}

func mergeOutputsIntoEnv(env model.Environment, outputs maps.Ordered[string, model.RunOutput]) model.Environment {
	newEnv := make([]model.EnvironmentValue, 0, outputs.Len())
	outputs.ForEach(func(key string, val model.RunOutput) error {
		newEnv = append(newEnv, model.EnvironmentValue{
			Key:   key,
			Value: val.Value,
		})

		return nil
	})

	return env.Merge(model.Environment{
		Values: newEnv,
	})
}
