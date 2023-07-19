package app

import (
	"context"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/transaction"
	"go.opentelemetry.io/otel/trace"
)

type runnerFacade struct {
	sm                *subscription.Manager
	testsPipeline     *TestPipeline
	transactionRunner executor.PersistentTransactionRunner
	assertionRunner   executor.AssertionRunner
	tracePoller       executor.PersistentTracePoller
	linterRunner      executor.LinterRunner
}

func (rf runnerFacade) StopTest(testID id.ID, runID int) {
	sr := executor.StopRequest{
		TestID: testID,
		RunID:  runID,
	}

	rf.sm.PublishUpdate(subscription.Message{
		ResourceID: sr.ResourceID(),
		Content:    sr,
	})
}

func (rf runnerFacade) RunTest(ctx context.Context, testObj test.Test, rm test.RunMetadata, env environment.Environment, gates *[]testrunner.RequiredGate) test.Run {
	return rf.testsPipeline.Run(ctx, testObj, rm, env, gates)
}

func (rf runnerFacade) RunTransaction(ctx context.Context, tr transaction.Transaction, rm test.RunMetadata, env environment.Environment, gates *[]testrunner.RequiredGate) transaction.TransactionRun {
	return rf.transactionRunner.Run(ctx, tr, rm, env, gates)
}

func (rf runnerFacade) RunAssertions(ctx context.Context, request executor.AssertionRequest) {
	rf.assertionRunner.RunAssertions(ctx, request)
}

func newRunnerFacades(
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
) *runnerFacade {
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
		ppRepo,
		tracer,
		execTestUpdater,
		tracedb.Factory(runRepo),
		dsRepo,
		eventEmitter,
	)

	pollerExecutor = executor.NewSelectorBasedPoller(pollerExecutor, eventEmitter)
	tracePoller := executor.NewTracePoller(
		pollerExecutor,
		ppRepo,
		execTestUpdater,
		linterRunner,
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
		ppRepo,
		trRepo,
	)

	queueBuilder := executor.NewQueueBuilder(runRepo, testRepo)
	pipeline := NewPipeline(queueBuilder,
		PipelineStep{processor: runner, driver: executor.NewInMemoryQueueDriver()},
		PipelineStep{processor: pollerExecutor, driver: executor.NewInMemoryQueueDriver()},
		PipelineStep{processor: linterRunner, driver: executor.NewInMemoryQueueDriver()},
		PipelineStep{processor: assertionRunner, driver: executor.NewInMemoryQueueDriver()},
	)

	pipeline.Start()

	testPipeline := NewTestPipeline(pipeline, runRepo, trRepo)

	transactionRunner := executor.NewTransactionRunner(
		runner,
		testRepo,
		transactionRunRepository,
		subscriptionManager,
	)

	return &runnerFacade{
		sm:                subscriptionManager,
		testsPipeline:     testPipeline,
		transactionRunner: transactionRunner,
		assertionRunner:   assertionRunner,
		tracePoller:       tracePoller,
		linterRunner:      linterRunner,
	}
}
