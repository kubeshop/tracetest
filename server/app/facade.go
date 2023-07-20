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
	"github.com/kubeshop/tracetest/server/transaction"
	"go.opentelemetry.io/otel/trace"
)

func buildTestPipeline(
	ppRepo *pollingprofile.Repository,
	dsRepo *datastore.Repository,
	lintRepo *analyzer.Repository,
	trRepo *testrunner.Repository,
	db model.Repository,
	testRepo test.Repository,
	runRepo test.RunRepository,
	transactionRunRepository *transaction.RunRepository,
	appTracer trace.Tracer,
	tracer trace.Tracer,
	subscriptionManager *subscription.Manager,
	triggerRegistry *trigger.Registry,
) *TestPipeline {
	eventEmitter := executor.NewEventEmitter(db, subscriptionManager)

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

	pollerExecutor := executor.NewPollerExecutor(
		tracer,
		execTestUpdater,
		tracedb.Factory(runRepo),
		dsRepo,
		eventEmitter,
	)

	pollerExecutor = executor.NewSelectorBasedPoller(pollerExecutor, eventEmitter)
	tracePoller := executor.NewTracePoller(
		pollerExecutor,
		execTestUpdater,
		subscriptionManager,
		eventEmitter,
	)

	runner := executor.NewPersistentRunner(
		triggerRegistry,
		runRepo,
		execTestUpdater,
		tracer,
		subscriptionManager,
		tracedb.Factory(runRepo),
		dsRepo,
		eventEmitter,
	)

	queueBuilder := executor.NewQueueBuilder().
		WithDataStoreGetter(dsRepo).
		WithPollingProfileGetter(ppRepo).
		WithTestGetter(testRepo).
		WithRunGetter(runRepo)

	pipeline := NewPipeline(queueBuilder,
		PipelineStep{processor: runner, driver: executor.NewInMemoryQueueDriver("runner")},
		PipelineStep{processor: tracePoller, driver: executor.NewInMemoryQueueDriver("tracePoller")},
		PipelineStep{processor: linterRunner, driver: executor.NewInMemoryQueueDriver("linterRunner")},
		PipelineStep{processor: assertionRunner, driver: executor.NewInMemoryQueueDriver("assertionRunner")},
	)

	pipeline.Start()

	return NewTestPipeline(
		pipeline,
		pipeline.GetQueueForStep(3), // assertion runner step
		runRepo,
		trRepo,
		ppRepo,
		dsRepo,
	)
}
