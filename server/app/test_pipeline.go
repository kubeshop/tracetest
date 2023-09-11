package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

func buildTestPipeline(
	pool *pgxpool.Pool,
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

	triggerResolverWorker := executor.NewTriggerResolverWorker(
		triggerRegistry,
		execTestUpdater,
		tracer,
		tracedbFactory,
		dsRepo,
		eventEmitter,
	)

	triggerExecuterWorker := executor.NewTriggerExecuterWorker(
		triggerRegistry,
		execTestUpdater,
		tracer,
		eventEmitter,
	)

	triggerResultProcessorWorker := executor.NewTriggerResultProcessorWorker(
		tracer,
		subscriptionManager,
		eventEmitter,
		execTestUpdater,
	)

	cancelRunHandlerFn := executor.HandleRunCancelation(execTestUpdater, tracer, eventEmitter)

	queueBuilder := executor.NewQueueConfigurer().
		WithCancelRunHandlerFn(cancelRunHandlerFn).
		WithSubscriptor(subscriptionManager).
		WithDataStoreGetter(dsRepo).
		WithPollingProfileGetter(ppRepo).
		WithTestGetter(testRepo).
		WithRunGetter(runRepo).
		WithInstanceID(instanceID)

	pgQueue := pipeline.NewPostgresQueueDriver[executor.Job](pool, pgChannelName)

	pipeline := pipeline.New(queueBuilder,
		pipeline.Step[executor.Job]{Processor: triggerResolverWorker, Driver: pgQueue.Channel("trigger_resolve")},
		pipeline.Step[executor.Job]{Processor: triggerExecuterWorker, Driver: pgQueue.Channel("trigger_execute")},
		pipeline.Step[executor.Job]{Processor: triggerResultProcessorWorker, Driver: pgQueue.Channel("trigger_result")},
		pipeline.Step[executor.Job]{Processor: tracePoller, Driver: pgQueue.Channel("tracePoller")},
		pipeline.Step[executor.Job]{Processor: linterRunner, Driver: pgQueue.Channel("linterRunner")},
		pipeline.Step[executor.Job]{Processor: assertionRunner, Driver: pgQueue.Channel("assertionRunner")},
	)

	const assertionRunnerStepIndex = 5

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
