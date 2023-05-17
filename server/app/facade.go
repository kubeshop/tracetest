package app

import (
	"context"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	lintern_resource "github.com/kubeshop/tracetest/server/lintern/resource"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/tracedb/datastoreresource"
	"go.opentelemetry.io/otel/trace"
)

type runnerFacade struct {
	sm                *subscription.Manager
	runner            executor.PersistentRunner
	transactionRunner executor.PersistentTransactionRunner
	assertionRunner   executor.AssertionRunner
	tracePoller       executor.PersistentTracePoller
	linternRunner     executor.LinternRunner
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

func (rf runnerFacade) RunTest(ctx context.Context, test model.Test, rm model.RunMetadata, env environment.Environment) model.Run {
	return rf.runner.Run(ctx, test, rm, env)
}

func (rf runnerFacade) RunTransaction(ctx context.Context, tr model.Transaction, rm model.RunMetadata, env environment.Environment) model.TransactionRun {
	return rf.transactionRunner.Run(ctx, tr, rm, env)
}

func (rf runnerFacade) RunAssertions(ctx context.Context, request executor.AssertionRequest) {
	rf.assertionRunner.RunAssertions(ctx, request)
}

func newRunnerFacades(
	ppRepo *pollingprofile.Repository,
	dsRepo *datastoreresource.Repository,
	lintRepo *lintern_resource.Repository,
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
		executor.NewAssertionExecutor(tracer),
		executor.InstrumentedOutputProcessor(tracer),
		subscriptionManager,
		eventEmitter,
	)

	linternRunner := executor.NewLinternRunner(
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
		tracedb.Factory(testDB),
		dsRepo,
		eventEmitter,
	)

	tracePoller := executor.NewTracePoller(
		pollerExecutor,
		ppRepo,
		execTestUpdater,
		linternRunner,
		subscriptionManager,
		eventEmitter,
	)

	runner := executor.NewPersistentRunner(
		triggerRegistry,
		testDB,
		execTestUpdater,
		tracePoller,
		tracer,
		subscriptionManager,
		tracedb.Factory(testDB),
		dsRepo,
		eventEmitter,
	)

	transactionRunner := executor.NewTransactionRunner(
		runner,
		testDB,
		subscriptionManager,
	)

	return &runnerFacade{
		sm:                subscriptionManager,
		runner:            runner,
		transactionRunner: transactionRunner,
		assertionRunner:   assertionRunner,
		tracePoller:       tracePoller,
		linternRunner:     linternRunner,
	}
}
