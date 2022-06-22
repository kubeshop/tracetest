package executor

import (
	"fmt"
	"math"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
)

type DefaultPollerExecutor struct {
	updater           RunUpdater
	traceDB           tracedb.TraceDB
	maxTracePollRetry int
}

var _ PollerExecutor = &DefaultPollerExecutor{}

func NewPollerExecutor(
	updater RunUpdater,
	traceDB tracedb.TraceDB,
	retryDelay time.Duration,
	maxWaitTimeForTrace time.Duration,
) PollerExecutor {
	maxTracePollRetry := int(math.Ceil(float64(maxWaitTimeForTrace) / float64(retryDelay)))
	return &DefaultPollerExecutor{
		updater:           updater,
		traceDB:           traceDB,
		maxTracePollRetry: maxTracePollRetry,
	}
}

func (pe DefaultPollerExecutor) ExecuteRequest(request PollingRequest) (bool, model.Run, error) {
	run := request.run
	otelTrace, err := pe.traceDB.GetTraceByID(request.ctx, run.TraceID.String())
	if err != nil {
		return false, model.Run{}, err
	}

	trace := traces.FromOtel(otelTrace)
	trace.ID = run.TraceID
	run.Trace = &trace
	request.run = run

	if !pe.donePollingTraces(request, trace) {
		return false, model.Run{}, nil
	}

	run = run.SuccessfullyPolledTraces(augmentData(&trace, run.Response))

	fmt.Printf("completed polling result %s after %d times, number of spans: %d \n", run.ID, request.count, len(run.Trace.Flat))

	err = pe.updater.Update(request.ctx, run)
	if err != nil {
		return false, model.Run{}, nil
	}

	return true, run, nil
}

func (pe DefaultPollerExecutor) donePollingTraces(job PollingRequest, trace traces.Trace) bool {
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
