package executor

import (
	"fmt"
	"log"
	"math"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DefaultPollerExecutor struct {
	config            config.Config
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
		attribute.String("tracetest.run.trace_poller.test_id", string(request.test.ID)),
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
	config config.Config,
	tracer trace.Tracer,
	updater RunUpdater,
	traceDB tracedb.TraceDB,
) PollerExecutor {
	retryDelay := config.PoolingRetryDelay()
	maxWaitTimeForTrace := config.MaxWaitTimeForTraceDuration()

	maxTracePollRetry := int(math.Ceil(float64(maxWaitTimeForTrace) / float64(retryDelay)))
	pollerExecutor := &DefaultPollerExecutor{
		config:            config,
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
	log.Printf("[PollerExecutor] Test %s Run %d: ExecuteRequest\n", request.test.ID, request.run.ID)
	run := request.run
	traceID := run.TraceID.String()

	trace, err := pe.traceDB.GetTraceByID(request.ctx, traceID)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: GetTraceByID (traceID %s) error: %s\n", request.test.ID, request.run.ID, traceID, err.Error())
		return false, model.Run{}, err
	}

	trace.ID = run.TraceID

	if !pe.donePollingTraces(request, trace) {
		log.Printf("[PollerExecutor] Test %s Run %d: Not done polling\n", request.test.ID, request.run.ID)
		run.Trace = &trace
		request.run = run
		return false, run, nil
	}

	log.Printf("[PollerExecutor] Test %s Run %d: Start Sorting\n", request.test.ID, request.run.ID)
	trace = trace.Sort()
	log.Printf("[PollerExecutor] Test %s Run %d: Sorting complete\n", request.test.ID, request.run.ID)
	run.Trace = &trace
	request.run = run

	run = run.SuccessfullyPolledTraces(augmentData(&trace, run.TriggerResult))

	fmt.Printf("completed polling result %d after %d times, number of spans: %d \n", run.ID, request.count, len(run.Trace.Flat))

	log.Printf("[PollerExecutor] Test %s Run %d: Start updating\n", request.test.ID, request.run.ID)
	err = pe.updater.Update(request.ctx, run)
	if err != nil {
		log.Printf("[PollerExecutor] Test %s Run %d: Update error: %s\n", request.test.ID, request.run.ID, err.Error())
		return false, model.Run{}, err
	}

	return true, run, nil
}

func (pe DefaultPollerExecutor) donePollingTraces(job *PollingRequest, trace traces.Trace) bool {
	// we're done if we have the same amount of spans after polling or `maxTracePollRetry` times
	if job.count == pe.maxTracePollRetry {
		log.Printf("[PollerExecutor] Test %s Run %d: Done polling. Hit MaxRetry of %d\n", job.test.ID, job.run.ID, pe.maxTracePollRetry)
		return true
	}

	if job.run.Trace == nil {
		return false
	}

	minimalNumberOfSpans := 0
	applicationExporter, _ := pe.config.ApplicationExporter()
	if applicationExporter != nil {
		// The triggering span will be sent to the application data storage, so we need to
		// expect at least 1 span in the trace.
		minimalNumberOfSpans = 1
	}

	if len(trace.Flat) > minimalNumberOfSpans && len(trace.Flat) == len(job.run.Trace.Flat) {
		log.Printf("[PollerExecutor] Test %s Run %d: Done polling. Condition met\n", job.test.ID, job.run.ID)
		return true
	}

	return false
}
