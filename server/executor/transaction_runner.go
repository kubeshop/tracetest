package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/subscription"
)

type TransactionRunner interface {
	Run(context.Context, model.Transaction, model.RunMetadata, model.Environment) model.TransactionRun
}

func NewTransactionRunner(
	runner Runner,
	db model.Repository,
	subscriptionManager *subscription.Manager,
	config config.Config,
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
		checkTestStateDelay: config.PoolingRetryDelay() / 2,
	}
}

type transactionRunJob struct {
	ctx context.Context
	run model.TransactionRun
}

type persistentTransactionRunner struct {
	testRunner          Runner
	db                  model.Repository
	updater             TransactionRunUpdater
	subscriptionManager *subscription.Manager
	checkTestStateDelay time.Duration
	executionChannel    chan transactionRunJob
	exit                chan bool
}

func (r persistentTransactionRunner) Run(ctx context.Context, transaction model.Transaction, metadata model.RunMetadata, environment model.Environment) model.TransactionRun {
	run := model.NewTransactionRun(transaction)
	run.Metadata = metadata
	run.Environment = environment

	ctx = getNewCtx(ctx)

	run, _ = r.db.CreateTransactionRun(ctx, run)

	r.executionChannel <- transactionRunJob{
		ctx: ctx,
		run: run,
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
					err := r.runTransaction(job.ctx, job.run)
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}()
	}
}

func (r persistentTransactionRunner) runTransaction(ctx context.Context, run model.TransactionRun) error {
	run.State = model.TransactionRunStateExecuting

	var err error

	for i := range run.Steps {
		run, err = r.runTransactionStep(ctx, run, i)
		if err != nil {
			return fmt.Errorf("could not execute step %d of transaction %s: %w", i, run.TransactionID, err)
		}

		if run.State == model.TransactionRunStateFailed {
			break
		}

		run.Environment = run.InjectOutputsIntoEnvironment(run.Environment)
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

func (r persistentTransactionRunner) runTransactionStep(ctx context.Context, transactionRun model.TransactionRun, stepIndex int) (model.TransactionRun, error) {
	step := transactionRun.Steps[stepIndex]
	test, err := r.db.GetLatestTestVersion(ctx, step.ID)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("could not load transaction step: %w", err)
	}

	testRun, completedTestChannel := r.testRunner.Run(ctx, test, transactionRun.Metadata, transactionRun.Environment)
	stepRun := createStepRun(testRun)
	transactionRun.StepRuns = append(transactionRun.StepRuns, stepRun)

	// listen for updates and propagate them as if they were transaction updates
	r.subscriptionManager.Subscribe(testRun.ResourceID(), subscription.NewSubscriberFunction(
		func(m subscription.Message) error {
			updatedRun := m.Content.(model.Run)

			for i, run := range transactionRun.StepRuns {
				if run.ID == updatedRun.ID && run.TestID == updatedRun.TestID {
					transactionRun.StepRuns[i] = createStepRun(updatedRun)
					break
				}
			}

			r.subscriptionManager.PublishUpdate(subscription.Message{
				ResourceID: transactionRun.ResourceID(),
				Type:       "result_update",
				Content:    transactionRun,
			})

			return nil
		}),
	)

	err = r.updater.Update(ctx, transactionRun)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("could not update transaction run: %w", err)
	}

	runResult := <-completedTestChannel
	testRun, err = runResult.Run, runResult.Err
	if err != nil {
		transactionRun.State = model.TransactionRunStateFailed
		transactionRun.StepRuns[stepIndex] = createStepRun(testRun)
		r.updater.Update(ctx, transactionRun)

		return model.TransactionRun{}, fmt.Errorf("could not run step: %w", err)
	}

	stepRun = createStepRun(testRun)

	transactionRun.StepRuns[stepIndex] = stepRun

	if testRun.State == model.RunStateFailed {
		transactionRun.State = model.TransactionRunStateFailed
	} else {
		transactionRun.CurrentTest += 1
	}

	err = r.updater.Update(ctx, transactionRun)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("could not update transaction run: %w", err)
	}

	return transactionRun, nil
}

func createStepRun(testRun model.Run) model.TransactionStepRun {
	return model.TransactionStepRun{
		ID:          testRun.ID,
		TestID:      testRun.TestID,
		State:       testRun.State,
		Environment: testRun.Environment,
		Outputs:     testRun.Outputs,
	}
}
