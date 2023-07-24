package executor

import (
	"context"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/transaction"
)

type TransactionPipeline struct {
	*Pipeline
	runs transactionsRunRepo
}

type transactionsRunRepo interface {
	CreateRun(context.Context, transaction.TransactionRun) (transaction.TransactionRun, error)
}

func NewTransactionPipeline(
	pipeline *Pipeline,
	runs transactionsRunRepo,
) *TransactionPipeline {
	return &TransactionPipeline{
		Pipeline: pipeline,
		runs:     runs,
	}
}

func (p *TransactionPipeline) Run(ctx context.Context, tran transaction.Transaction, metadata test.RunMetadata, environment environment.Environment, requiredGates *[]testrunner.RequiredGate) transaction.TransactionRun {
	tranRun := tran.NewRun()
	tranRun.Metadata = metadata
	tranRun.Environment = environment
	tranRun.RequiredGates = requiredGates

	tranRun, _ = p.runs.CreateRun(ctx, tranRun)

	job := Job{
		Transaction:    tran,
		TransactionRun: tranRun,
	}

	p.Pipeline.Begin(ctx, job)

	return tranRun
}
