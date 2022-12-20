package app

import (
	"context"

	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb"
)

type runnerFacade struct {
	newTraceDB        func(ds model.DataStore) (tracedb.TraceDB, error)
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

func (rf runnerFacade) NewTraceDB(ds model.DataStore) (tracedb.TraceDB, error) {
	return rf.newTraceDB(ds)
}
