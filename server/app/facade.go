package app

import (
	"context"

	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

type runnerFacade struct {
	runner            executor.PersistentRunner
	transactionRunner executor.PersistentTransactionRunner
	assertionRunner   executor.AssertionRunner
	tracePoller       executor.PersistentTracePoller
}

func (rf runnerFacade) RunTest(ctx context.Context, test model.Test, rm model.RunMetadata, env model.Environment) model.Run {
	return rf.runner.Run(ctx, test, rm, env)
}

func (rf runnerFacade) RunTransaction(ctx context.Context, tr model.Transaction, rm model.RunMetadata, env model.Environment) model.TransactionRun {
	return rf.transactionRunner.Run(ctx, tr, rm, env)
}

func (rf runnerFacade) RunAssertions(ctx context.Context, request executor.AssertionRequest) {
	rf.assertionRunner.RunAssertions(ctx, request)
}

func newRunnerFacades(
	ppRepo *pollingprofile.Repository,
	testDB model.Repository,
	appTracer trace.Tracer,
	tracer trace.Tracer,
	subscriptionManager *subscription.Manager,
	triggerRegistry *trigger.Registry,
) *runnerFacade {
	eventEmitter := executor.NewEventEmitter(testDB, subscriptionManager)

	execTestUpdater := (executor.CompositeUpdater{}).
		Add(executor.NewDBUpdater(testDB)).
		Add(executor.NewSubscriptionUpdater(subscriptionManager))

	assertionRunner := executor.NewAssertionRunner(
		execTestUpdater,
		executor.NewAssertionExecutor(tracer, eventEmitter),
		executor.InstrumentedOutputProcessor(tracer),
		subscriptionManager,
		eventEmitter,
	)

	pollerExecutor := executor.NewPollerExecutor(
		ppRepo,
		tracer,
		execTestUpdater,
		tracedb.Factory(testDB),
		testDB,
	)

	tracePoller := executor.NewTracePoller(
		pollerExecutor,
		ppRepo,
		execTestUpdater,
		assertionRunner,
		subscriptionManager,
	)

	runner := executor.NewPersistentRunner(
		triggerRegistry,
		testDB,
		execTestUpdater,
		tracePoller,
		tracer,
		subscriptionManager,
		tracedb.Factory(testDB),
		testDB,
		eventEmitter,
	)

	transactionRunner := executor.NewTransactionRunner(
		runner,
		testDB,
		subscriptionManager,
	)

	return &runnerFacade{
		runner:            runner,
		transactionRunner: transactionRunner,
		assertionRunner:   assertionRunner,
		tracePoller:       tracePoller,
	}
}
