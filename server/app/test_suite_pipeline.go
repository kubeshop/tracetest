package app

import (
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/testsuite"
)

func buildTestSuitePipeline(
	tranRepo *testsuite.Repository,
	runRepo *testsuite.RunRepository,
	testRunner *executor.TestPipeline,
	subscriptionManager *subscription.Manager,
) *executor.TestSuitesPipeline {
	tranRunner := executor.NewTestSuiteRunner(testRunner, runRepo, subscriptionManager)
	queueBuilder := executor.NewQueueBuilder().
		WithTestSuiteGetter(tranRepo).
		WithTestSuiteRunGetter(runRepo)

	pipeline := executor.NewPipeline(queueBuilder,
		executor.PipelineStep{Processor: tranRunner, Driver: executor.NewInMemoryQueueDriver("testSuiteRunner")},
	)

	return executor.NewTestSuitePipeline(
		pipeline,
		runRepo,
	)
}
