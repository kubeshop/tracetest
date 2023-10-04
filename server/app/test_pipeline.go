package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/executor/tracepollerworker"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/metric"
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
	appConfig *config.AppConfig,
	meter metric.Meter,
) *executor.TestPipeline {
	eventEmitter := executor.NewEventEmitter(treRepo, subscriptionManager)

	execTestUpdater := (executor.CompositeUpdater{}).
		Add(executor.NewDBUpdater(runRepo)).
		Add(executor.NewSubscriptionUpdater(subscriptionManager))

	workerMetricMiddlewareBuilder := executor.NewWorkerMetricMiddlewareBuilder(meter)

	assertionRunner := executor.NewAssertionRunner(
		execTestUpdater,
		executor.NewAssertionExecutor(tracer),
		executor.InstrumentedOutputProcessor(tracer),
		subscriptionManager,
		eventEmitter,
	)

	linterRunner := executor.NewLinterRunner(
		execTestUpdater,
		subscriptionManager,
		eventEmitter,
		lintRepo,
	)

	tracePollerStarterWorker := tracepollerworker.NewStarterWorker(
		eventEmitter,
		tracedbFactory,
		dsRepo,
		execTestUpdater,
		subscriptionManager,
		tracer,
	)

	traceFetcherWorker := tracepollerworker.NewFetcherWorker(
		eventEmitter,
		tracedbFactory,
		dsRepo,
		execTestUpdater,
		subscriptionManager,
		tracer,
		appConfig.TestPipelineTraceFetchingEnabled(),
	)

	tracePollerEvaluatorWorker := tracepollerworker.NewEvaluatorWorker(
		eventEmitter,
		tracedbFactory,
		dsRepo,
		execTestUpdater,
		subscriptionManager,
		tracepollerworker.NewSelectorBasedPollingStopStrategy(eventEmitter, tracepollerworker.NewSpanCountPollingStopStrategy()),
		tracer,
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
		appConfig.TestPipelineTriggerExecutionEnabled(),
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
		WithInstanceID(instanceID).
		WithMetricMeter(meter)

	pgQueue := pipeline.NewPostgresQueueDriver[executor.Job](pool, pgChannelName)

	pipeline := pipeline.New(queueBuilder,
		pipeline.Step[executor.Job]{Processor: workerMetricMiddlewareBuilder.New("trigger_resolver", triggerResolverWorker), Driver: pgQueue.Channel("trigger_resolve")},
		pipeline.Step[executor.Job]{Processor: workerMetricMiddlewareBuilder.New("trigger_executer", triggerExecuterWorker), Driver: pgQueue.Channel("trigger_execute")},
		pipeline.Step[executor.Job]{Processor: workerMetricMiddlewareBuilder.New("trigger_result_processor", triggerResultProcessorWorker), Driver: pgQueue.Channel("trigger_result")},
		pipeline.Step[executor.Job]{Processor: workerMetricMiddlewareBuilder.New("trace_poller_starter", tracePollerStarterWorker), Driver: pgQueue.Channel("tracePoller_start")},
		pipeline.Step[executor.Job]{Processor: workerMetricMiddlewareBuilder.New("trace_fetcher", traceFetcherWorker), Driver: pgQueue.Channel("tracePoller_fetch")},
		pipeline.Step[executor.Job]{Processor: workerMetricMiddlewareBuilder.New("trace_poller_evaluator", tracePollerEvaluatorWorker), Driver: pgQueue.Channel("tracePoller_evaluate"), InputQueueOffset: -1},
		pipeline.Step[executor.Job]{Processor: workerMetricMiddlewareBuilder.New("linter_runner", linterRunner), Driver: pgQueue.Channel("linterRunner")},
		pipeline.Step[executor.Job]{Processor: workerMetricMiddlewareBuilder.New("assertion_runner", assertionRunner), Driver: pgQueue.Channel("assertionRunner")},
	)

	const assertionRunnerStepIndex = 7

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
