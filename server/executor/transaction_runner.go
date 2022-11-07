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

func NewTransactionRunner(runner Runner) PersistentTransactionRunner {
	return PersistentTransactionRunner{
		testRunner:       runner,
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
	executionChannel chan transactionRunJob
	exit             chan bool
}

func (r PersistentTransactionRunner) Run(ctx context.Context, transaction model.Transaction, metadata model.RunMetadata, environment model.Environment) model.TransactionRun {
	run := model.TransactionRun{
		TransactionID:      transaction.ID,
		TransactionVersion: transaction.Version,
		CreatedAt:          time.Now(),
		State:              model.TransactionRunStateCreated,
		Metadata:           metadata,
		Environment:        environment,
		Steps:              transaction.Steps,
		StepRuns:           make([]model.Run, len(transaction.Steps)),
		CurrentTest:        0,
	}

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
					r.runTransaction(job.ctx, job.run)
				}
			}
		}()
	}
}

func (r PersistentTransactionRunner) runTransaction(ctx context.Context, run model.TransactionRun) {
	for i := range run.Steps {
		r.runTransactionStep(ctx, &run, i)
	}
}

func (r PersistentTransactionRunner) runTransactionStep(ctx context.Context, transactionRun *model.TransactionRun, stepIndex int) {
	step := transactionRun.Steps[stepIndex]
	testRun := r.testRunner.Run(ctx, step, transactionRun.Metadata, transactionRun.Environment)
	transactionRun.StepRuns[stepIndex] = testRun
}
