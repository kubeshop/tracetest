package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/model"
)

type TransactionRunner interface {
	Run(context.Context, model.Transaction, model.RunMetadata, model.Environment) model.TransactionRun
}

func NewTransactionRunner(runner Runner, db model.Repository) PersistentTransactionRunner {
	return PersistentTransactionRunner{
		testRunner:       runner,
		db:               db,
		executionChannel: make(chan transactionRunJob, 1),
		exit:             make(chan bool),
	}
}

type transactionRunJob struct {
	ctx context.Context
	run model.TransactionRun
}

type PersistentTransactionRunner struct {
	testRunner       Runner
	db               model.Repository
	executionChannel chan transactionRunJob
	exit             chan bool
}

func (r PersistentTransactionRunner) Run(ctx context.Context, transaction model.Transaction, metadata model.RunMetadata, environment model.Environment) model.TransactionRun {
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

func (r PersistentTransactionRunner) Start(workers int) {
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

func (r PersistentTransactionRunner) runTransaction(ctx context.Context, run model.TransactionRun) error {
	run.State = model.TransactionRunStateExecuting
	environment := run.Environment

	var err error

	for i := range run.Steps {
		run, err = r.runTransactionStep(ctx, run, i, environment)
		if err != nil {
			return fmt.Errorf("could not execute step %d of transaction %s: %w", i, run.TransactionID, err)
		}

		if run.State == model.TransactionRunStateExecuting {
			environment = r.patchEnvironment(environment, run)
		}
	}

	if !run.State.IsFinal() {
		run.State = model.TransactionRunStateFinished
	}

	return r.db.UpdateTransactionRun(ctx, run)
}

func (r PersistentTransactionRunner) runTransactionStep(ctx context.Context, transactionRun model.TransactionRun, stepIndex int, environment model.Environment) (model.TransactionRun, error) {
	step := transactionRun.Steps[stepIndex]
	testRun := r.testRunner.Run(ctx, step, transactionRun.Metadata, environment)

	testRun, err := r.waitTestRunIsFinished(ctx, testRun)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("could not wait for step execution: %w", err)
	}

	transactionRun.StepRuns = append(transactionRun.StepRuns, testRun)

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

func (r PersistentTransactionRunner) waitTestRunIsFinished(ctx context.Context, testRun model.Run) (model.Run, error) {
	var err error
	for !testRun.State.IsFinal() {
		time.Sleep(2 * time.Second)

		testRun, err = r.db.GetRun(ctx, testRun.TestID, testRun.ID)
		if err != nil {
			return model.Run{}, err
		}
	}

	return testRun, nil
}

func (r PersistentTransactionRunner) patchEnvironment(baseEnvironment model.Environment, run model.TransactionRun) model.Environment {
	newEnvVariables := make([]model.EnvironmentValue, 0)

	for i := 0; i < run.CurrentTest; i++ {
		testRun := run.StepRuns[i]
		testRun.Outputs.ForEach(func(key, val string) error {
			newEnvVariables = append(newEnvVariables, model.EnvironmentValue{
				Key:   key,
				Value: val,
			})

			return nil
		})
	}

	newEnvironment := model.Environment{Values: newEnvVariables}

	return baseEnvironment.Merge(newEnvironment)
}
