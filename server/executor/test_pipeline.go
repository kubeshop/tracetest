package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
)

type TestPipeline struct {
	*Pipeline
	assertionQueue Enqueuer
	runs           runsRepo
	trGetter       testRunnerGetter
	ppGetter       defaultPollingProfileGetter
	dsGetter       currentDataStoreGetter
}

type runsRepo interface {
	CreateRun(context.Context, test.Test, test.Run) (test.Run, error)
	GetRun(_ context.Context, testID id.ID, runID int) (test.Run, error)
}

type testRunnerGetter interface {
	GetDefault(ctx context.Context) testrunner.TestRunner
}

type defaultPollingProfileGetter interface {
	GetDefault(ctx context.Context) pollingprofile.PollingProfile
}

func NewTestPipeline(
	pipeline *Pipeline,
	assertionQueue Enqueuer,
	runs runsRepo,
	trGetter testRunnerGetter,
	ppGetter defaultPollingProfileGetter,
	dsGetter currentDataStoreGetter,
) *TestPipeline {
	return &TestPipeline{
		Pipeline:       pipeline,
		assertionQueue: assertionQueue,
		runs:           runs,
		trGetter:       trGetter,
		ppGetter:       ppGetter,
		dsGetter:       dsGetter,
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

	datastore, err := p.dsGetter.Current(ctx)
	p.handleDBError(run, err)

	job := NewJob()
	job.Test = testObj
	job.Run = run
	job.PollingProfile = p.ppGetter.GetDefault(ctx)
	job.DataStore = datastore

	p.Pipeline.Begin(ctx, job)

	return run
}

func (p *TestPipeline) Rerun(ctx context.Context, testObj test.Test, runID int) test.Run {
	run, err := p.runs.GetRun(ctx, testObj.ID, runID)
	p.handleDBError(run, err)

	newTestRun, err := p.runs.CreateRun(ctx, testObj, run.Copy())
	p.handleDBError(run, err)

	ds, err := p.dsGetter.Current(ctx)
	p.handleDBError(run, err)

	p.assertionQueue.Enqueue(ctx, Job{
		Test:           testObj,
		Run:            newTestRun,
		PollingProfile: p.ppGetter.GetDefault(ctx),
		DataStore:      ds,
	})

	return newTestRun
}

func (p *TestPipeline) handleDBError(run test.Run, err error) {
	if err != nil {
		fmt.Printf("test %s run #%d DB error: %s\n", run.TestID, run.ID, err.Error())
	}
}
