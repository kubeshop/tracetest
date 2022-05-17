package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/traces"
	"go.opentelemetry.io/otel/trace"
)

type (
	Test struct {
		ID               uuid.UUID
		Name             string
		Description      string
		ServiceUnderTest ServiceUnderTest
		ReferenceRun     Run
		Definition       Definition
	}

	ServiceUnderTest struct {
		Request HTTPRequest
	}

	Definition map[SpanQuery][]Assertion

	SpanQuery string
	Assertion struct {
		ID         string
		Attribute  string
		Comparator comparator.Comparator
		Value      string
	}

	Run struct {
		ID          uuid.UUID
		TraceID     trace.TraceID
		SpanID      trace.SpanID
		State       RunState
		LastError   error
		CreatedAt   time.Time
		CompletedAt time.Time
		Request     HTTPRequest
		Response    HTTPResponse
		Trace       traces.Trace
		Results     RunResults
	}

	RunResults struct {
		AllPassed bool
		Results   Results
	}

	Results map[SpanQuery][]AssertionResult

	AssertionResult struct {
		Assertion Assertion
		AllPassed bool
		Results   []SpanAssertionResult
	}

	SpanAssertionResult struct {
		SpanID        trace.SpanID
		ObservedValue string
		CompareErr    error
	}
)

type RunState string

const (
	RunStateCreated             RunState = "CREATED"
	RunStateExecuting           RunState = "EXECUTING"
	RunStateAwaitingTrace       RunState = "AWAITING_TRACE"
	RunStateFailed              RunState = "FAILED"
	RunStateFinished            RunState = "FINISHED"
	RunStateAwaitingTestResults RunState = "AWAITING_TEST_RESULTS"
)
