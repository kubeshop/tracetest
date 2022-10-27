package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

type (
	Test struct {
		ID               id.ID
		CreatedAt        time.Time
		Name             string
		Description      string
		Version          int
		ServiceUnderTest Trigger
		Specs            OrderedMap[SpanQuery, NamedAssertions]
		Outputs          OrderedMap[string, Output]
		Summary          Summary
	}

	Output struct {
		Selector SpanQuery
		Value    string
	}

	NamedAssertions struct {
		Name       string
		Assertions []Assertion
	}

	Summary struct {
		Runs    int
		LastRun TestLastRun
	}

	TestLastRun struct {
		Time   time.Time
		Passes int
		Fails  int
	}

	TriggerType string

	Trigger struct {
		Type TriggerType
		HTTP *HTTPRequest
		GRPC *GRPCRequest
	}

	TriggerResult struct {
		Type TriggerType
		HTTP *HTTPResponse
		GRPC *GRPCResponse
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

	Attribute string

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
		Trace         *traces.Trace
		Outputs       OrderedMap[string, string]
		LastError     error
		Pass          int
		Fail          int

		Metadata RunMetadata

		// environment
		EnvSnapshot []EnvironmentValue
	}

	RunResults struct {
		AllPassed bool
		Results   OrderedMap[SpanQuery, []AssertionResult]
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

	Environment struct {
		ID          string
		Name        string
		Description string
		CreatedAt   time.Time
		Values      []EnvironmentValue
	}

	EnvironmentValue struct {
		Key   string
		Value string
	}
)

func (t Test) HasID() bool {
	return t.ID != ""
}

const (
	metaPrefix    = "tracetest.selected_spans."
	metaPrefixLen = len("tracetest.selected_spans.")
)

func (a Attribute) IsMeta() bool {
	return len(a) > metaPrefixLen && a[0:metaPrefixLen] == metaPrefix
}

func (a Attribute) String() string {
	return string(a)
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

func (e Environment) GetSlug() string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(e.Name), " ", "-"))
}

type RunState string

const (
	RunStateCreated             RunState = "CREATED"
	RunStateExecuting           RunState = "EXECUTING"
	RunStateAwaitingTrace       RunState = "AWAITING_TRACE"
	RunStateFailed              RunState = "FAILED"
	RunStateFinished            RunState = "FINISHED"
	RunStateAwaitingTestResults RunState = "AWAITING_TEST_RESULTS"
)
