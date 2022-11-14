package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/model"
)

type TransactionRunner interface {
	Run(context.Context, model.Transaction, model.RunMetadata, model.Environment) model.TransactionRun
}

func NewTransactionRunner(runner Runner, db model.Repository, config config.Config) persistentTransactionRunner {
	return persistentTransactionRunner{
		testRunner:          runner,
		db:                  db,
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
	checkTestStateDelay time.Duration
	executionChannel    chan transactionRunJob
	exit                chan bool
}

func (r persistentTransactionRunner) Run(ctx context.Context, transaction model.Transaction, metadata model.RunMetadata, environment model.Environment) model.TransactionRun {
	run := model.NewTransactionRun(transaction)
	run.Metadata = metadata
	run.Environment = environment

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
						panic(err)
					}
				}
			}
		}()
	}
}

func (r persistentTransactionRunner) runTransaction(ctx context.Context, run model.TransactionRun) error {
	run.State = model.TransactionRunStateExecuting
	environment := run.Environment

	var err error

	for i := range run.Steps {
		run, err = r.runTransactionStep(ctx, run, i, environment)
		if err != nil {
			return fmt.Errorf("could not execute step %d of transaction %s: %w", i, run.TransactionID, err)
		}

		if run.State == model.TransactionRunStateFailed {
			break
		}

		environment = r.patchEnvironment(environment, run)
	}

	if run.State != model.TransactionRunStateFailed {
		run.State = model.TransactionRunStateFinished
	}

	return r.db.UpdateTransactionRun(ctx, run)
}

func (r persistentTransactionRunner) runTransactionStep(ctx context.Context, transactionRun model.TransactionRun, stepIndex int, environment model.Environment) (model.TransactionRun, error) {
	step := transactionRun.Steps[stepIndex]
	test, err := r.db.GetLatestTestVersion(ctx, step.ID)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("could not load transaction step: %w", err)
	}

	_, completedTestChannel := r.testRunner.Run(ctx, test, transactionRun.Metadata, environment)
	runResult := <-completedTestChannel
	testRun, err := runResult.Run, runResult.Err
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("could not run step: %w", err)
	}

	stepRun := model.TransactionStepRun{ID: testRun.ID, TestID: testRun.TestID, State: testRun.State}

	transactionRun.StepRuns = append(transactionRun.StepRuns, stepRun)

	if testRun.State == model.RunStateFailed {
		transactionRun.State = model.TransactionRunStateFailed
	} else {
		transactionRun.CurrentTest += 1
	}

	err = r.db.UpdateTransactionRun(ctx, transactionRun)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("could not update transaction run: %w", err)
	}

	return transactionRun, nil
}

func (r persistentTransactionRunner) patchEnvironment(baseEnvironment model.Environment, run model.TransactionRun) model.Environment {
	if run.CurrentTest == 0 {
		return baseEnvironment
	}

	lastExecutedTest := run.StepRuns[run.CurrentTest-1]
	lastEnvironment := lastExecutedTest.Environment
	newEnvVariables := make([]model.EnvironmentValue, 0)
	lastExecutedTest.Outputs.ForEach(func(key, val string) error {
		newEnvVariables = append(newEnvVariables, model.EnvironmentValue{
			Key:   key,
			Value: val,
		})

		return nil
	})

	newEnvironment := model.Environment{Values: newEnvVariables}

	return lastEnvironment.Merge(newEnvironment)
}
