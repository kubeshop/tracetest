package executor

import (
	"context"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/variableset"
	"go.opentelemetry.io/otel/trace"
)

type TestPipeline struct {
	*pipeline.Pipeline[Job]
	updatePublisher    updatePublisher
	assertionQueue     pipeline.Enqueuer[Job]
	runs               runsRepo
	trGetter           testRunnerGetter
	ppGetter           defaultPollingProfileGetter
	dsGetter           currentDataStoreGetter
	cancelRunHandlerFn runCancelHandlerFn
}

type runsRepo interface {
	CreateRun(context.Context, test.Test, test.Run) (test.Run, error)
	UpdateRun(context.Context, test.Run) error
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
	pipeline *pipeline.Pipeline[Job],

	updatePublisher updatePublisher,

	assertionQueue pipeline.Enqueuer[Job],
	runs runsRepo,
	trGetter testRunnerGetter,
	ppGetter defaultPollingProfileGetter,
	dsGetter currentDataStoreGetter,

	cancelRunHandlerFn runCancelHandlerFn,
) *TestPipeline {
	return &TestPipeline{
		Pipeline:           pipeline,
		updatePublisher:    updatePublisher,
		assertionQueue:     assertionQueue,
		runs:               runs,
		trGetter:           trGetter,
		ppGetter:           ppGetter,
		dsGetter:           dsGetter,
		cancelRunHandlerFn: cancelRunHandlerFn,
	}
}

func (p *TestPipeline) Run(ctx context.Context, testObj test.Test, metadata test.RunMetadata, variableSet variableset.VariableSet, requiredGates *[]testrunner.RequiredGate) test.Run {
	run := test.NewRun()
	run.Metadata = metadata
	run.VariableSet = variableSet

	// configuring required gates
	if requiredGates == nil {
		rg := p.trGetter.GetDefault(ctx).RequiredGates
		requiredGates = &rg
	}
	run = run.ConfigureRequiredGates(*requiredGates)
	run.SkipTraceCollection = testObj.SkipTraceCollection

	run, err := p.runs.CreateRun(ctx, testObj, run)
	p.handleDBError(run, err)

	datastore, err := p.dsGetter.Current(ctx)
	p.handleDBError(run, err)

	job := NewJob()
	job.TenantID = middleware.TenantIDFromContext(ctx)
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

	err = p.runs.UpdateRun(ctx, newTestRun.SuccessfullyPolledTraces(run.Trace))
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
	sr := UserRequest{
		TenantID: middleware.TenantIDFromContext(ctx),
		TestID:   testID,
		RunID:    runID,
		Type:     string(UserRequestTypeStop),
	}

	p.updatePublisher.PublishUpdate(subscription.Message{
		ResourceID: sr.ResourceID(UserRequestTypeStop),
		Content:    sr,
	})
}

func (p *TestPipeline) SkipTraceCollection(ctx context.Context, testID id.ID, runID int) {
	sr := UserRequest{
		TenantID: middleware.TenantIDFromContext(ctx),
		TestID:   testID,
		RunID:    runID,
		Type:     string(UserRequestTypeSkipTraceCollection),
	}

	p.updatePublisher.PublishUpdate(subscription.Message{
		ResourceID: sr.ResourceID(UserRequestTypeSkipTraceCollection),
		Content:    sr,
	})
}

func (p *TestPipeline) UpdateStoppedTest(ctx context.Context, run test.Run) {
	p.cancelRunHandlerFn(ctx, run)
}

type runCancelHandlerFn func(ctx context.Context, run test.Run) error

var ErrUserCancelled = fmt.Errorf("cancelled by user")

func RunWasUserCancelled(run test.Run) bool {
	// depeending on when the Run was cancelled (which step was being executed)
	// the error might be set on different fields
	return (run.TriggerResult.Error != nil &&
		ErrorMessageIsUserCancelled(run.TriggerResult.Error.ErrorMessage)) ||
		(run.LastError != nil && ErrorMessageIsUserCancelled(run.LastError.Error()))
}

func ErrorMessageIsUserCancelled(msg string) bool {
	return msg == ErrUserCancelled.Error()
}

func ErrorMessageIsSkipTraceCollection(msg string) bool {
	return msg == ErrSkipTraceCollection.Error()
}

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
