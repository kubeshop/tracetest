package executor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"go.opentelemetry.io/otel/propagation"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

const PollingRequestStartIteration = 1

type TracePoller interface {
	Poll(context.Context, model.Test, model.Run)
}

type PersistentTracePoller interface {
	TracePoller
	WorkerPool
}

type PollerExecutor interface {
	ExecuteRequest(*PollingRequest) (bool, string, model.Run, error)
}

type TraceFetcher interface {
	GetTraceByID(ctx context.Context, traceID string) (*v1.TracesData, error)
}

type PollingProfileGetter interface {
	GetDefault(ctx context.Context) pollingprofile.PollingProfile
}

func NewTracePoller(
	pe PollerExecutor,
	ppGetter PollingProfileGetter,
	updater RunUpdater,
	linterRunner LinterRunner,
	subscriptionManager *subscription.Manager,
	eventEmitter EventEmitter,
) PersistentTracePoller {
	return tracePoller{
		updater:             updater,
		ppGetter:            ppGetter,
		pollerExecutor:      pe,
		executeQueue:        make(chan PollingRequest, 5),
		exit:                make(chan bool, 1),
		linterRunner:        linterRunner,
		subscriptionManager: subscriptionManager,
		eventEmitter:        eventEmitter,
	}
}

type tracePoller struct {
	updater             RunUpdater
	ppGetter            PollingProfileGetter
	pollerExecutor      PollerExecutor
	linterRunner        LinterRunner
	subscriptionManager *subscription.Manager
	eventEmitter        EventEmitter

	executeQueue chan PollingRequest
	exit         chan bool
}

type PollingRequest struct {
	test           model.Test
	run            model.Run
	pollingProfile pollingprofile.PollingProfile
	count          int
	headers        map[string]string
}

func (r PollingRequest) Context() context.Context {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	return propagator.Extract(context.Background(), propagation.MapCarrier(r.headers))
}

func (pr *PollingRequest) SetHeader(name, value string) {
	pr.headers[name] = value
}

func (pr *PollingRequest) Header(name string) string {
	return pr.headers[name]
}

func (pr PollingRequest) IsFirstRequest() bool {
	return pr.Header("requeued") != "true"
}

func NewPollingRequest(ctx context.Context, test model.Test, run model.Run, count int) *PollingRequest {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	request := &PollingRequest{
		test:    test,
		run:     run,
		headers: make(map[string]string),
		count:   count,
	}

	propagator.Inject(ctx, propagation.MapCarrier(request.headers))

	return request
}

func (tp tracePoller) handleDBError(err error) {
	if err != nil {
		fmt.Printf("DB error when polling traces: %s\n", err.Error())
	}
}

func (tp tracePoller) Start(workers int) {
	for i := 0; i < workers; i++ {
		go func() {
			fmt.Println("tracePoller start goroutine")
			for {
				select {
				case <-tp.exit:
					fmt.Println("tracePoller exit goroutine")
					return
				case job := <-tp.executeQueue:
					log.Printf("[TracePoller] Test %s Run %d: Received job\n", job.test.ID, job.run.ID)
					tp.processJob(job)
				}
			}
		}()
	}
}

func (tp tracePoller) Stop() {
	fmt.Println("tracePoller stopping")
	tp.exit <- true
}

func (tp tracePoller) Poll(ctx context.Context, test model.Test, run model.Run) {
	log.Printf("[TracePoller] Test %s Run %d: Poll\n", test.ID, run.ID)

	job := NewPollingRequest(ctx, test, run, PollingRequestStartIteration)

	tp.enqueueJob(*job)
}

func (tp tracePoller) enqueueJob(job PollingRequest) {
	tp.executeQueue <- job
}

func (tp tracePoller) processJob(job PollingRequest) {
	ctx := job.Context()
	select {
	default:
	case <-ctx.Done():
		log.Printf("[TracePoller] Context cancelled.")
		return
	}

	if job.IsFirstRequest() {
		err := tp.eventEmitter.Emit(ctx, events.TraceFetchingStart(job.test.ID, job.run.ID))
		if err != nil {
			log.Printf("[TracePoller] Test %s Run %d: fail to emit TracePollingStart event: %s \n", job.test.ID, job.run.ID, err.Error())
		}
	}

	fmt.Println("tracePoller processJob", job.count)

	finished, finishReason, run, err := tp.pollerExecutor.ExecuteRequest(&job)
	if err != nil {
		log.Printf("[TracePoller] Test %s Run %d: ExecuteRequest Error: %s\n", job.test.ID, job.run.ID, err.Error())
		jobFailed, reason := tp.handleTraceDBError(job, err)

		if jobFailed {
			anotherErr := tp.eventEmitter.Emit(ctx, events.TracePollingError(job.test.ID, job.run.ID, reason, err))
			if anotherErr != nil {
				log.Printf("[TracePoller] Test %s Run %d: fail to emit TracePollingError event: %s \n", job.test.ID, job.run.ID, err.Error())
			}

			anotherErr = tp.eventEmitter.Emit(ctx, events.TraceFetchingError(job.test.ID, job.run.ID, err))
			if anotherErr != nil {
				log.Printf("[TracePoller] Test %s Run %d: fail to emit TracePollingError event: %s \n", job.test.ID, job.run.ID, err.Error())
			}
		}

		return
	}

	if !finished {
		job.count += 1
		tp.requeue(job)
		return
	}

	log.Printf("[TracePoller] Test %s Run %d: Done polling (reason: %s). Completed polling after %d iterations, number of spans collected %d\n", job.test.ID, job.run.ID, finishReason, job.count+1, len(run.Trace.Flat))

	err = tp.eventEmitter.Emit(ctx, events.TraceFetchingSuccess(job.test.ID, job.run.ID))
	if err != nil {
		log.Printf("[TracePoller] Test %s Run %d: fail to emit TracePollingSuccess event: %s \n", job.test.ID, job.run.ID, err.Error())
	}

	tp.handleDBError(tp.updater.Update(ctx, run))

	job.run = run
	tp.runAssertions(job)
}

func (tp tracePoller) runAssertions(job PollingRequest) {
	linterRequest := LinterRequest{
		Test: job.test,
		Run:  job.run,
	}

	tp.linterRunner.RunLinter(job.Context(), linterRequest)
}

func (tp tracePoller) handleTraceDBError(job PollingRequest, err error) (bool, string) {
	run := job.run
	ctx := job.Context()

	profile := tp.ppGetter.GetDefault(ctx)
	if profile.Periodic == nil {
		log.Println("[TracePoller] cannot get polling profile.")
		return true, "Cannot get polling profile"
	}

	pp := *profile.Periodic

	// Edge case: the trace still not avaiable on Data Store during polling
	if errors.Is(err, connection.ErrTraceNotFound) && time.Since(run.ServiceTriggeredAt) < pp.TimeoutDuration() {
		log.Println("[TracePoller] Trace not found on Data Store yet. Requeuing...")
		job.count += 1
		tp.requeue(job)
		return false, "Trace not found" // job without fail
	}

	reason := ""

	if errors.Is(err, connection.ErrTraceNotFound) {
		reason = fmt.Sprintf("Timed out without finding trace, trace id \"%s\"", run.TraceID.String())

		err = fmt.Errorf("timed out waiting for traces after %s", pp.Timeout)
		fmt.Println("[TracePoller] Timed-out", err)
	} else {
		reason = "Unexpected error"

		err = fmt.Errorf("cannot fetch trace: %w", err)
		fmt.Println("[TracePoller] Unknown error", err)
	}

	run = run.TraceFailed(err)
	analytics.SendEvent("test_run_finished", "error", "", &map[string]string{
		"finalState": string(run.State),
	})

	tp.handleDBError(tp.updater.Update(ctx, run))

	tp.subscriptionManager.PublishUpdate(subscription.Message{
		ResourceID: run.TransactionStepResourceID(),
		Type:       "update_run",
		Content:    RunResult{Run: run, Err: err},
	})

	return true, reason // job failed
}

func (tp tracePoller) requeue(job PollingRequest) {
	go func() {
		pp := tp.ppGetter.GetDefault(job.Context())
		fmt.Printf("[TracePoller] Requeuing Test Run %d. Current iteration: %d\n", job.run.ID, job.count)
		time.Sleep(pp.Periodic.RetryDelayDuration())

		job.SetHeader("requeued", "true")
		job.pollingProfile = pp
		tp.enqueueJob(job)
	}()
}
