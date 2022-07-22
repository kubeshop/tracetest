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

	spanCount := 0
	if run.Trace != nil {
		spanCount = len(run.Trace.Flat)
	}

	attrs := []attribute.KeyValue{
		attribute.String("tracetest.run.trace_poller.trace_id", request.run.TraceID.String()),
		attribute.String("tracetest.run.trace_poller.span_id", request.run.SpanID.String()),
		attribute.Bool("tracetest.run.trace_poller.succesful", finished),
		attribute.String("tracetest.run.trace_poller.test_id", request.test.ID.String()),
		attribute.Int("tracetest.run.trace_poller.amount_retrieved_spans", spanCount),
	}

	if err != nil {
		attrs = append(attrs, attribute.String("tracetest.run.trace_poller.error", err.Error()))
		span.RecordError(err)
	}

	span.SetAttributes(attrs...)
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
	traceId := run.TraceID.String()

	otelTrace, err := pe.traceDB.GetTraceByID(request.ctx, traceId)
	if err != nil {
		return false, model.Run{}, err
	}

	trace := traces.FromOtel(otelTrace)
	trace.ID = run.TraceID

	if !pe.donePollingTraces(request, trace) {
		run.Trace = &trace
		request.run = run
		return false, run, nil
	}

	trace = trace.Sort()
	run.Trace = &trace
	request.run = run

	run = run.SuccessfullyPolledTraces(augmentData(&trace, run.TriggerResult))

	fmt.Printf("completed polling result %s after %d times, number of spans: %d \n", run.ID, request.count, len(run.Trace.Flat))

	err = pe.updater.Update(request.ctx, run)
	if err != nil {
		return false, model.Run{}, err
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

	// We always expect to get the tracetest trigger span, so it has to have more than 1 span
	if len(trace.Flat) > 1 && len(trace.Flat) == len(job.run.Trace.Flat) {
		return true
	}

	return false
}
