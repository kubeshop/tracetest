package app

import (
	"context"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
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
	runner            executor.PersistentRunner
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

func (rf runnerFacade) RunTest(ctx context.Context, test test.Test, rm test.RunMetadata, env environment.Environment) test.Run {
	return rf.runner.Run(ctx, test, rm, env)
}

func (rf runnerFacade) RunTransaction(ctx context.Context, tr transaction.Transaction, rm test.RunMetadata, env environment.Environment) transaction.TransactionRun {
	return rf.transactionRunner.Run(ctx, tr, rm, env)
}

func (rf runnerFacade) RunAssertions(ctx context.Context, request executor.AssertionRequest) {
	rf.assertionRunner.RunAssertions(ctx, request)
}

func newRunnerFacades(
	ppRepo *pollingprofile.Repository,
	dsRepo *datastore.Repository,
	lintRepo *analyzer.Repository,
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
		assertionRunner,
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
		tracePoller,
		tracer,
		subscriptionManager,
		tracedb.Factory(runRepo),
		dsRepo,
		eventEmitter,
		ppRepo,
	)

	transactionRunner := executor.NewTransactionRunner(
		runner,
		testRepo,
		transactionRunRepository,
		subscriptionManager,
	)

	return &runnerFacade{
		sm:                subscriptionManager,
		runner:            runner,
		transactionRunner: transactionRunner,
		assertionRunner:   assertionRunner,
		tracePoller:       tracePoller,
		linterRunner:      linterRunner,
	}
}
