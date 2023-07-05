package model

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"go.opentelemetry.io/otel/trace"
)

type (
	// this struct yaml/json encoding is handled at ./test_json.go for custom encodings
	Test struct {
		ID               id.ID
		CreatedAt        time.Time
		Name             string
		Description      string
		Version          int
		ServiceUnderTest Trigger
		Specs            maps.Ordered[SpanQuery, NamedAssertions]
		Outputs          maps.Ordered[string, Output]
		Summary          Summary
	}

	Output struct {
		Selector SpanQuery
		Value    string `expr_enabled:"true"`
	}

	NamedAssertions struct {
		Name       string      `expr_enabled:"true"`
		Assertions []Assertion `stmt_enabled:"true"`
	}

	Summary struct {
		Runs    int     `json:"runs"`
		LastRun LastRun `json:"lastRun"`
	}

	LastRun struct {
		Time          time.Time `json:"time"`
		Passes        int       `json:"passes"`
		Fails         int       `json:"fails"`
		AnalyzerScore int       `json:"analyzerScore"`
	}

	TriggerType string

	// this struct yaml/json encoding is handled at ./trigger_json.go for custom encodings
	Trigger struct {
		Type    TriggerType
		HTTP    *HTTPRequest
		GRPC    *GRPCRequest
		TraceID *TRACEIDRequest
	}

	TriggerResult struct {
		Type    TriggerType
		HTTP    *HTTPResponse
		GRPC    *GRPCResponse
		TRACEID *TRACEIDResponse
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
		TriggerResult TriggerResult
		Results       *RunResults
		Trace         *Trace
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
		Linter           analyzer.LinterResult
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

func (sar SpanAssertionResult) SafeSpanIDString() string {
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

	if e.Expression == nil {
		if e.LiteralValue.Type == "attribute" {
			return fmt.Sprintf("attr:%s", e.LiteralValue.Value)
		}

		return e.LiteralValue.Value
	}

	return fmt.Sprintf("%s %s %s", e.LiteralValue.Value, e.Operation, e.Expression.String())
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
