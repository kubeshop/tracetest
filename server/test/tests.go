package test

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"go.opentelemetry.io/otel/trace"
)

const (
	ResourceName       string = "test"
	ResourceNamePlural string = "tests"
)

type (
	// this struct yaml/json encoding is handled at ./test_json.go for custom encodings
	Test struct {
		ID          id.ID                        `json:"id"`
		CreatedAt   time.Time                    `json:"createdAt"`
		Name        string                       `json:"name"`
		Description string                       `json:"description"`
		Version     int                          `json:"version"`
		Trigger     trigger.Trigger              `json:"trigger"`
		Specs       Specs                        `json:"specs"`
		Outputs     maps.Ordered[string, Output] `json:"outputs"`
		Summary     Summary                      `json:"summary"`
	}

	Specs []TestSpec

	Output struct {
		Selector SpanQuery `json:"selector"`
		Value    string    `json:"value" expr_enabled:"true"`
	}

	TestSpec struct {
		Selector   Selector    `json:"selector"`
		Name       string      `json:"name,omitempty"`
		Assertions []Assertion `json:"assertions"`
	}

	Selector struct {
		Query          SpanQuery    `json:"query"`
		ParsedSelector SpanSelector `json:"parsedSelector"`
	}

	SpanSelector struct {
		Filters       []SelectorFilter     `json:"filters"`
		PseudoClass   *SelectorPseudoClass `json:"pseudoClass,omitempty"`
		ChildSelector *SpanSelector        `json:"childSelector,omitempty"`
	}

	SelectorFilter struct {
		Property string `json:"property"`
		Operator string `json:"operator"`
		Value    string `json:"value"`
	}

	SelectorPseudoClass struct {
		Name     string `json:"name"`
		Argument *int32 `json:"argument,omitempty"`
	}

	Summary struct {
		Runs    int     `json:"runs"`
		LastRun LastRun `json:"lastRun"`
	}

	LastRun struct {
		Time   time.Time `json:"time"`
		Passes int       `json:"passes"`
		Fails  int       `json:"fails"`
	}

	SpanQuery string

	Assertion string

	AssertionExpression struct {
		LiteralValue LiteralValue
		Operation    string
		Expression   *AssertionExpression
	}

	LiteralValue struct {
		Value string
		Type  string
	}

	RunMetadata map[string]string

	Run struct {
		ID          int
		TestID      id.ID
		TestVersion int

		// Timestamps
		CreatedAt                 time.Time
		ServiceTriggeredAt        time.Time
		ServiceTriggerCompletedAt time.Time
		ObtainedTraceAt           time.Time
		CompletedAt               time.Time

		// trigger params
		State   RunState
		TraceID trace.TraceID
		SpanID  trace.SpanID

		// result info
		TriggerResult trigger.TriggerResult
		Results       *RunResults
		Trace         *model.Trace
		Outputs       maps.Ordered[string, RunOutput]
		LastError     error
		Pass          int
		Fail          int

		Metadata RunMetadata

		// environment
		Environment environment.Environment

		// transaction

		TransactionID    string
		TransactionRunID string
		Linter           model.LinterResult
	}

	RunResults struct {
		AllPassed bool
		Results   maps.Ordered[SpanQuery, []AssertionResult]
	}

	RunOutput struct {
		Name     string
		Value    string
		SpanID   string
		Resolved bool
		Error    error
	}

	AssertionResult struct {
		Assertion Assertion
		AllPassed bool
		Results   []SpanAssertionResult
	}

	SpanAssertionResult struct {
		SpanID        *trace.SpanID
		ObservedValue string
		CompareErr    error
	}
)

func (t Test) GetID() id.ID {
	return t.ID
}

func (t Test) Validate() error {
	return nil
}

func (sar SpanAssertionResult) SpanIDString() string {
	if sar.SpanID == nil {
		return ""
	}

	return sar.SpanID.String()
}

func (t Test) HasID() bool {
	return t.ID != ""
}

func (e *AssertionExpression) String() string {
	if e == nil {
		return ""
	}

	if e.Expression != nil {
		return fmt.Sprintf("%s %s %s", e.LiteralValue.Value, e.Operation, e.Expression.String())
	}

	if e.LiteralValue.Type == "attribute" {
		return fmt.Sprintf("attr:%s", e.LiteralValue.Value)
	}

	return e.LiteralValue.Value
}

func (e *AssertionExpression) Type() string {
	if e == nil {
		return "string"
	}

	if e.Expression == nil {
		return e.LiteralValue.Type
	}

	if e.LiteralValue.Type == "attribute" {
		// in case of attributes, we check the rest of the expression to know its type
		return e.Expression.Type()
	}

	return e.LiteralValue.Type
}

type RunState string

const (
	RunStateCreated             RunState = "CREATED"
	RunStateExecuting           RunState = "EXECUTING"
	RunStateAwaitingTrace       RunState = "AWAITING_TRACE"
	RunStateTriggerFailed       RunState = "TRIGGER_FAILED"
	RunStateTraceFailed         RunState = "TRACE_FAILED"
	RunStateAssertionFailed     RunState = "ASSERTION_FAILED"
	RunStateAnalyzingTrace      RunState = "ANALYZING_TRACE"
	RunStateAnalyzingError      RunState = "ANALYZING_ERROR"
	RunStateFinished            RunState = "FINISHED"
	RunStateStopped             RunState = "STOPPED"
	RunStateAwaitingTestResults RunState = "AWAITING_TEST_RESULTS"
)

func (rs RunState) IsFinal() bool {
	return rs.IsError() || rs == RunStateFinished
}

func (rs RunState) IsError() bool {
	return rs == RunStateTriggerFailed ||
		rs == RunStateTraceFailed
}

func (r Run) ResultsCount() (pass, fail int) {
	if r.Results == nil {
		return
	}

	r.Results.Results.ForEach(func(_ SpanQuery, ars []AssertionResult) error {
		for _, ar := range ars {
			for _, rs := range ar.Results {
				if rs.CompareErr == nil {
					pass++
				} else {
					fail++
				}
			}
		}
		return nil
	})

	return
}
