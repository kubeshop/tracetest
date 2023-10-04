package app

import (
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/testsuite"
	"go.opentelemetry.io/otel/metric"
)

func buildTestSuitePipeline(
	tranRepo *testsuite.Repository,
	runRepo *testsuite.RunRepository,
	testRunner *executor.TestPipeline,
	subscriptionManager *subscription.Manager,
	meter metric.Meter,
) *executor.TestSuitesPipeline {
	tranRunner := executor.NewTestSuiteRunner(testRunner, runRepo, subscriptionManager)
	queueBuilder := executor.NewQueueConfigurer().
		WithTestSuiteGetter(tranRepo).
		WithTestSuiteRunGetter(runRepo).
		WithMetricMeter(meter)

	pipeline := pipeline.New(queueBuilder,
		pipeline.Step[executor.Job]{Processor: tranRunner, Driver: pipeline.NewInMemoryQueueDriver[executor.Job]("testSuiteRunner")},
	)

	return executor.NewTestSuitePipeline(
		pipeline,
		runRepo,
	)
}
