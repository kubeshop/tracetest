package executor

import (
	"fmt"
	"math"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DefaultPollerExecutor struct {
	updater           RunUpdater
	traceDB           tracedb.TraceDB
	maxTracePollRetry int
}

type InstrumentedPollerExecutor struct {
	tracer         trace.Tracer
	pollerExecutor PollerExecutor
}

func (pe InstrumentedPollerExecutor) ExecuteRequest(request *PollingRequest) (bool, model.Run, error) {
	_, span := pe.tracer.Start(request.ctx, "Fetch trace")
	defer span.End()

	finished, run, err := pe.pollerExecutor.ExecuteRequest(request)

	span.SetAttributes(
		attribute.Bool("tracetest.run.trace_poller.succesful", finished),
		attribute.String("tracetest.run.trace_poller.test_id", request.test.ID.String()),
	)

	return finished, run, err
}

func NewPollerExecutor(
	tracer trace.Tracer,
	updater RunUpdater,
	traceDB tracedb.TraceDB,
	retryDelay time.Duration,
	maxWaitTimeForTrace time.Duration,
) PollerExecutor {
	maxTracePollRetry := int(math.Ceil(float64(maxWaitTimeForTrace) / float64(retryDelay)))
	pollerExecutor := &DefaultPollerExecutor{
		updater:           updater,
		traceDB:           traceDB,
		maxTracePollRetry: maxTracePollRetry,
	}

	return &InstrumentedPollerExecutor{
		tracer:         tracer,
		pollerExecutor: pollerExecutor,
	}
}

func (pe DefaultPollerExecutor) ExecuteRequest(request *PollingRequest) (bool, model.Run, error) {
	run := request.run
	otelTrace, err := pe.traceDB.GetTraceByID(request.ctx, run.TraceID.String())
	if err != nil {
		return false, model.Run{}, err
	}

	trace := traces.FromOtel(otelTrace)
	trace.ID = run.TraceID

	if !pe.donePollingTraces(request, trace) {
		run.Trace = &trace
		request.run = run
		return false, model.Run{}, nil
	}

	trace = trace.Sort()
	run.Trace = &trace
	request.run = run

	run = run.SuccessfullyPolledTraces(augmentData(&trace, run.Response))

	fmt.Printf("completed polling result %s after %d times, number of spans: %d \n", run.ID, request.count, len(run.Trace.Flat))

	err = pe.updater.Update(request.ctx, run)
	if err != nil {
		return false, model.Run{}, nil
	}

	return true, run, nil
}

func (pe DefaultPollerExecutor) donePollingTraces(job *PollingRequest, trace traces.Trace) bool {
	// we're done if we have the same amount of spans after polling or `maxTracePollRetry` times
	if job.count == pe.maxTracePollRetry {
		return true
	}

	if job.run.Trace == nil {
		return false
	}

	if len(trace.Flat) > 0 && len(trace.Flat) == len(job.run.Trace.Flat) {
		return true
	}

	return false
}
