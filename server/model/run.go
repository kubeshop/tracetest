package model

import (
	"math"
	"time"

	"github.com/kubeshop/tracetest/server/traces"
)

var (
	Now = func() time.Time {
		return time.Now()
	}
)

func (r Run) Copy(testVersion int) Run {
	r.Results = nil
	r.Trace = nil
	r.TestVersion = testVersion

	return r
}

func (r Run) ExecutionTime() int {
	var endDate time.Time
	if !r.CompletedAt.IsZero() {
		endDate = r.CompletedAt
	} else {
		endDate = Now()
	}

	et := math.Ceil(endDate.Sub(r.CreatedAt).Seconds())

	return int(et)
}

func (r Run) Start() Run {
	r.State = RunStateExecuting
	r.ServiceTriggeredAt = Now()
	return r
}

func (r Run) SuccessfullyExecuted() Run {
	r.State = RunStateAwaitingTrace
	r.ServiceTriggerCompletedAt = Now()
	return r
}

func (r Run) SuccessfullyPolledTraces(t *traces.Trace) Run {
	r.State = RunStateAwaitingTestResults
	r.Trace = t
	r.ObtainedTraceAt = time.Now()
	return r
}

func (r Run) SuccessfullyAsserted(res Results, allPassed bool) Run {
	r.Results = &RunResults{
		AllPassed: allPassed,
		Results:   res,
	}
	r.State = RunStateFinished
	return r.Finish()
}

func (r Run) Finish() Run {
	r.CompletedAt = Now()
	return r
}

func (r Run) Failed(err error) Run {
	r.State = RunStateFailed
	r.LastError = err
	return r.Finish()
}
