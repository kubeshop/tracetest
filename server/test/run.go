package test

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/pkg/timing"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/kubeshop/tracetest/server/variableset"
)

var (
	IDGen = id.NewRandGenerator()
)

func (r *Run) MarshalJSON() ([]byte, error) {
	encoded, err := EncodeRun(*r)
	if err != nil {
		return nil, err
	}

	return json.Marshal(encoded)
}

func (r *Run) UnmarshalJSON(data []byte) error {
	var encoded EncodedRun

	if err := json.Unmarshal(data, &encoded); err != nil {
		return err
	}

	decoded, err := encoded.ToRun()
	if err != nil {
		return err
	}

	*r = decoded

	return nil
}

func NewRun() Run {
	return Run{
		ID:        0,
		TraceID:   IDGen.TraceID(),
		SpanID:    IDGen.SpanID(),
		State:     RunStateCreated,
		CreatedAt: timing.Now(),
	}
}

func (r Run) ResourceID() string {
	return fmt.Sprintf("test/%s/run/%d", r.TestID, r.ID)
}

func (r Run) TransactionStepResourceID() string {
	return fmt.Sprintf("transaction_step/%s/run/%d", r.TestID, r.ID)
}

func (r Run) Copy() Run {
	r.ID = 0
	r.Results = nil
	r.CreatedAt = timing.Now()

	return r
}

func (r Run) ExecutionTime() int {
	diff := timing.TimeDiff(r.CreatedAt, r.CompletedAt)
	return int(math.Ceil(diff.Seconds()))
}

func (r Run) TriggerTime() int {
	diff := timing.TimeDiff(r.ServiceTriggeredAt, r.ServiceTriggerCompletedAt)
	return int(diff.Milliseconds())
}

func (r Run) Start() Run {
	r.State = RunStateExecuting
	r.ServiceTriggeredAt = timing.Now()
	return r
}

func (r Run) TriggerCompleted(tr trigger.TriggerResult) Run {
	r.ServiceTriggerCompletedAt = timing.Now()
	r.TriggerResult = tr
	return r
}

func (r Run) SuccessfullyTriggered() Run {
	r.State = RunStateAwaitingTrace
	return r
}

func (r Run) SuccessfullyPolledTraces(t *traces.Trace) Run {
	r.State = RunStateAnalyzingTrace
	r.Trace = t
	r.ObtainedTraceAt = time.Now()
	return r
}

func (r Run) SuccessfullyAsserted(
	outputs maps.Ordered[string, RunOutput],
	variableSet variableset.VariableSet,
	res maps.Ordered[SpanQuery, []AssertionResult],
	allPassed bool,
) Run {
	r.Outputs = outputs
	r.VariableSet = variableSet
	r.Results = &RunResults{
		AllPassed: allPassed,
		Results:   res,
	}
	r.State = RunStateFinished

	if !allPassed {
		r = r.RequiredGateFailed(testrunner.RequiredGateTestSpecs)
	}

	return r.Finish()
}

func (r Run) Finish() Run {
	r.CompletedAt = timing.Now()
	return r
}

func (r Run) Stopped() Run {
	r.State = RunStateStopped
	return r.Finish()
}

func (r Run) TriggerFailed(err error) Run {
	r.State = RunStateTriggerFailed
	r.LastError = err
	return r.Finish()
}

func (r Run) TraceFailed(err error) Run {
	r.State = RunStateTraceFailed
	r.LastError = err
	return r.Finish()
}

func (r Run) AssertionFailed(err error) Run {
	r.State = RunStateAssertionFailed
	r.LastError = err
	return r.Finish()
}

func (r Run) LinterError(err error) Run {
	r.State = RunStateAnalyzingError
	r.LastError = err
	return r.Finish()
}

func (r Run) SuccessfulLinterExecution(linter analyzer.LinterResult) Run {
	r.State = RunStateAwaitingTestResults
	r.Linter = linter

	if !r.Linter.Passed {
		r = r.RequiredGateFailed(testrunner.RequiredGateAnalyzerRules)
	}

	if r.Linter.Score < linter.MinimumScore {
		r = r.RequiredGateFailed(testrunner.RequiredGateAnalyzerScore)
	}

	return r
}

func (r Run) ConfigureRequiredGates(gates []testrunner.RequiredGate) Run {
	r.RequiredGatesResult = testrunner.NewRequiredGatesResult(gates)
	return r
}

func (r Run) RequiredGateFailed(gate testrunner.RequiredGate) Run {
	r.RequiredGatesResult = r.RequiredGatesResult.OnFailed(gate)
	return r
}

func (r Run) GenerateRequiredGateResult(gates []testrunner.RequiredGate) testrunner.RequiredGatesResult {
	requiredGatesResult := testrunner.NewRequiredGatesResult(gates)
	if !r.Results.AllPassed {
		requiredGatesResult = requiredGatesResult.OnFailed(testrunner.RequiredGateTestSpecs)
	}

	if !r.Linter.Passed {
		requiredGatesResult = requiredGatesResult.OnFailed(testrunner.RequiredGateAnalyzerRules)
	}

	if r.Linter.Score < r.Linter.MinimumScore {
		requiredGatesResult = requiredGatesResult.OnFailed(testrunner.RequiredGateAnalyzerScore)
	}

	return requiredGatesResult
}

func NewTracetestRootSpan(run Run) traces.Span {
	return traces.AugmentRootSpan(traces.Span{
		ID:         id.NewRandGenerator().SpanID(),
		Name:       traces.TriggerSpanName,
		StartTime:  run.ServiceTriggeredAt,
		EndTime:    run.ServiceTriggerCompletedAt,
		Attributes: traces.Attributes{},
		Children:   []*traces.Span{},
	}, run.TriggerResult)
}
