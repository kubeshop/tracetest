package app

import (
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

func buildTestPipeline(
	ppRepo *pollingprofile.Repository,
	dsRepo *datastore.Repository,
	lintRepo *analyzer.Repository,
	trRepo *testrunner.Repository,
	treRepo model.TestRunEventRepository,
	testRepo test.Repository,
	runRepo test.RunRepository,
	tracer trace.Tracer,
	subscriptionManager *subscription.Manager,
	triggerRegistry *trigger.Registry,
	tracedbFactory tracedb.FactoryFunc,
) *executor.TestPipeline {
	eventEmitter := executor.NewEventEmitter(treRepo, subscriptionManager)

	execTestUpdater := (executor.CompositeUpdater{}).
		Add(executor.NewDBUpdater(runRepo)).
		Add(executor.NewSubscriptionUpdater(subscriptionManager))

	assertionRunner := executor.NewAssertionRunner(
		execTestUpdater,
		executor.NewAssertionExecutor(tracer),
		executor.InstrumentedOutputProcessor(tracer),
		subscriptionManager,
		eventEmitter,
	)

	linterRunner := executor.NewlinterRunner(
		execTestUpdater,
		subscriptionManager,
		eventEmitter,
		lintRepo,
	)

	pollerExecutor := executor.NewSelectorBasedPoller(
		executor.NewPollerExecutor(
			tracer,
			execTestUpdater,
			tracedbFactory,
			dsRepo,
			eventEmitter,
		),
		eventEmitter,
	)

	tracePoller := executor.NewTracePoller(
		pollerExecutor,
		execTestUpdater,
		subscriptionManager,
		eventEmitter,
	)

	runner := executor.NewPersistentRunner(
		triggerRegistry,
		execTestUpdater,
		tracer,
		subscriptionManager,
		tracedbFactory,
		dsRepo,
		eventEmitter,
	)

	cancelRunHandlerFn := executor.HandleRunCancelation(execTestUpdater, tracer, eventEmitter)

	queueBuilder := executor.NewQueueBuilder().
		WithCancelRunHandlerFn(cancelRunHandlerFn).
		WithSubscriptor(subscriptionManager).
		WithDataStoreGetter(dsRepo).
		WithPollingProfileGetter(ppRepo).
		WithTestGetter(testRepo).
		WithRunGetter(runRepo)

	pipeline := executor.NewPipeline(queueBuilder,
		executor.PipelineStep{Processor: runner, Driver: executor.NewInMemoryQueueDriver("runner")},
		executor.PipelineStep{Processor: tracePoller, Driver: executor.NewInMemoryQueueDriver("tracePoller")},
		executor.PipelineStep{Processor: linterRunner, Driver: executor.NewInMemoryQueueDriver("linterRunner")},
		executor.PipelineStep{Processor: assertionRunner, Driver: executor.NewInMemoryQueueDriver("assertionRunner")},
	)

	pipeline.Start()

	const assertionRunnerStepIndex = 3

	return executor.NewTestPipeline(
		pipeline,
		subscriptionManager,
		pipeline.GetQueueForStep(assertionRunnerStepIndex), // assertion runner step
		runRepo,
		trRepo,
		ppRepo,
		dsRepo,
	)
}
