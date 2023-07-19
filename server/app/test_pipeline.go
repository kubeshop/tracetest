package app

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/test"
)

type TestPipeline struct {
	pipeline *Pipeline
	runs     test.RunRepository
	trGetter TestRunnerGetter
}

type TestRunnerGetter interface {
	GetDefault(ctx context.Context) testrunner.TestRunner
}

func NewTestPipeline(
	pipeline *Pipeline,
	runs test.RunRepository,
	trGetter TestRunnerGetter,
) *TestPipeline {
	return &TestPipeline{
		pipeline: pipeline,
		runs:     runs,
		trGetter: trGetter,
	}
}

func (p *TestPipeline) Run(ctx context.Context, testObj test.Test, metadata test.RunMetadata, environment environment.Environment, requiredGates *[]testrunner.RequiredGate) test.Run {
	run := test.NewRun()
	run.Metadata = metadata
	run.Environment = environment

	// configuring required gates
	if requiredGates == nil {
		rg := p.trGetter.GetDefault(ctx).RequiredGates
		requiredGates = &rg
	}
	run = run.ConfigureRequiredGates(*requiredGates)

	run, err := p.runs.CreateRun(ctx, testObj, run)
	p.handleDBError(run, err)

	// r.listenForStopRequests(ctx, cancelCtx, run)

	p.pipeline.Begin(ctx, executor.Job{
		InitialJobConfigurations: executor.InitialJobConfigurations{
			DataStoreID:      datastore.DataStoreSingleID,
			PollingProfileID: pollingprofile.DefaultPollingProfile.ID,
		},
		Test: testObj,
		Run:  run,
	})

	return run
}

func (p *TestPipeline) handleDBError(run test.Run, err error) {
	if err != nil {
		fmt.Printf("test %s run #%d trigger DB error: %s\n", run.TestID, run.ID, err.Error())
	}
}
