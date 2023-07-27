package app

import (
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/transaction"
)

func buildTransactionPipeline(
	tranRepo *transaction.Repository,
	runRepo *transaction.RunRepository,
	testRunner *executor.TestPipeline,
	subscriptionManager *subscription.Manager,
) *executor.TransactionPipeline {
	tranRunner := executor.NewTransactionRunner(testRunner, runRepo, subscriptionManager)
	queueBuilder := executor.NewQueueBuilder().
		WithTransactionGetter(tranRepo).
		WithTransactionRunGetter(runRepo)

	pipeline := executor.NewPipeline(queueBuilder,
		executor.PipelineStep{Processor: tranRunner, Driver: executor.NewInMemoryQueueDriver("transactionRunner")},
	)

	return executor.NewTransactionPipeline(
		pipeline,
		runRepo,
	)
}
