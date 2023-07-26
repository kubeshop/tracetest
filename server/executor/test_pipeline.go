package executor

import (
	"context"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"go.opentelemetry.io/otel/trace"
)

type TestPipeline struct {
	*Pipeline
	updatePublisher updatePublisher
	assertionQueue  Enqueuer
	runs            runsRepo
	trGetter        testRunnerGetter
	ppGetter        defaultPollingProfileGetter
	dsGetter        currentDataStoreGetter
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

type updatePublisher interface {
	PublishUpdate(subscription.Message)
}

func NewTestPipeline(
	pipeline *Pipeline,

	updatePublisher updatePublisher,

	assertionQueue Enqueuer,
	runs runsRepo,
	trGetter testRunnerGetter,
	ppGetter defaultPollingProfileGetter,
	dsGetter currentDataStoreGetter,
) *TestPipeline {
	return &TestPipeline{
		Pipeline:        pipeline,
		updatePublisher: updatePublisher,
		assertionQueue:  assertionQueue,
		runs:            runs,
		trGetter:        trGetter,
		ppGetter:        ppGetter,
		dsGetter:        dsGetter,
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

func (p *TestPipeline) StopTest(ctx context.Context, testID id.ID, runID int) {
	sr := StopRequest{
		TestID: testID,
		RunID:  runID,
	}

	p.updatePublisher.PublishUpdate(subscription.Message{
		ResourceID: sr.ResourceID(),
		Content:    sr,
	})
}

type runCancelHandlerFn func(ctx context.Context, run test.Run) error

func HandleRunCancelation(updater RunUpdater, tracer trace.Tracer, eventEmitter EventEmitter) runCancelHandlerFn {
	return func(ctx context.Context, run test.Run) error {
		ctx, span := tracer.Start(ctx, "User Requested Stop Run")
		defer span.End()

		if run.State == test.RunStateStopped {
			return nil
		}
		err := updater.Update(ctx, run.Stopped())
		if err != nil {
			fmt.Printf("test %s run #%d cancel DB error: %s\n", run.TestID, run.ID, err.Error())
		}

		evt := events.TraceStoppedInfo(run.TestID, run.ID)
		err = eventEmitter.Emit(ctx, evt)
		if err != nil {
			log.Printf("[HandleRunCancelation] Test %s Run %d: fail to emit TraceStoppedInfo event: %s \n", run.TestID, run.ID, err.Error())
			return err
		}

		return nil
	}
}

func (p *TestPipeline) handleDBError(run test.Run, err error) {
	if err != nil {
		fmt.Printf("test %s run #%d DB error: %s\n", run.TestID, run.ID, err.Error())
	}
}
