package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

type (
	Test struct {
		ID               uuid.UUID
		Name             string
		Description      string
		Version          int
		ServiceUnderTest ServiceUnderTest
		ReferenceRun     *Run
		Definition       OrderedMap[SpanQuery, []Assertion]
	}

	ServiceUnderTest struct {
		Request HTTPRequest
	}

	SpanQuery string

	Assertion struct {
		Attribute  string
		Comparator comparator.Comparator
		Value      string
	}

	Run struct {
		ID                        uuid.UUID
		TraceID                   trace.TraceID
		SpanID                    trace.SpanID
		State                     RunState
		LastError                 error
		CreatedAt                 time.Time
		ServiceTriggeredAt        time.Time
		ServiceTriggerCompletedAt time.Time
		ObtainedTraceAt           time.Time
		CompletedAt               time.Time
		Request                   HTTPRequest
		Response                  HTTPResponse
		Trace                     *traces.Trace
		Results                   *RunResults
		TestVersion               int
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
