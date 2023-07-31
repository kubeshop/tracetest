package executor

import (
	"context"

	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/transaction"
	"github.com/kubeshop/tracetest/server/variableset"
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

func (p *TransactionPipeline) Run(ctx context.Context, tran transaction.Transaction, metadata test.RunMetadata, variableSet variableset.VariableSet, requiredGates *[]testrunner.RequiredGate) transaction.TransactionRun {
	tranRun := tran.NewRun()
	tranRun.Metadata = metadata
	tranRun.VariableSet = variableSet
	tranRun.RequiredGates = requiredGates

	tranRun, _ = p.runs.CreateRun(ctx, tranRun)

	job := NewJob()
	job.Transaction = tran
	job.TransactionRun = tranRun

	p.Pipeline.Begin(ctx, job)

	return tranRun
}
