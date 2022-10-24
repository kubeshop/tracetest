package model

import (
	"fmt"
	"math"
	"time"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/traces"
)

var (
	Now = func() time.Time {
		return time.Now()
	}

	IDGen = id.NewRandGenerator()
)

func NewRun() Run {
	return Run{
		ID:        0,
		TraceID:   IDGen.TraceID(),
		SpanID:    IDGen.SpanID(),
		State:     RunStateCreated,
		CreatedAt: Now(),
	}
}

func (r Run) ResourceID() string {
	return fmt.Sprintf("test/%s/run/%d", r.TestID, r.ID)
}

func (r Run) Copy() Run {
	r.ID = 0
	r.Results = nil
	r.CreatedAt = Now()
	r.ServiceTriggeredAt = time.Time{}
	r.ServiceTriggerCompletedAt = time.Time{}
	r.ObtainedTraceAt = time.Time{}
	r.CompletedAt = time.Time{}

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

func (r Run) SuccessfullyAsserted(outputs OrderedMap[string, string], res OrderedMap[SpanQuery, []AssertionResult], allPassed bool) Run {
	r.Outputs = outputs
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
