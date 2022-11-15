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
		return time.Now().UTC()
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
	return timeDiff(r.CreatedAt, r.CompletedAt)
}

func (r Run) TriggerTime() int {
	return timeDiff(r.ServiceTriggeredAt, r.ServiceTriggerCompletedAt)
}

func timeDiff(start, end time.Time) int {
	var endDate time.Time
	if !dateIsZero(end) {
		endDate = end
	} else {
		endDate = Now()
	}

	et := math.Ceil(endDate.Sub(start).Seconds())

	return int(et)
}

func dateIsZero(in time.Time) bool {
	return in.IsZero() || in.Unix() == 0
}

func (r Run) Start() Run {
	r.State = RunStateExecuting
	r.ServiceTriggeredAt = Now()
	return r
}

func (r Run) TriggerCompleted(tr TriggerResult) Run {
	r.ServiceTriggerCompletedAt = Now()
	r.TriggerResult = tr
	return r
}

func (r Run) SuccessfullyTriggered() Run {
	r.State = RunStateAwaitingTrace
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
