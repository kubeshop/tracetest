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
		Definition       Definition
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

type Definition struct {
	list        [][]Assertion
	keyPosition map[SpanQuery]int
	positionKey map[int]SpanQuery
}

func (d *Definition) replace(d2 *Definition) {
	*d = *d2

}

func (d Definition) MarshalJSON() ([]byte, error) {
	type s struct {
		Selector   string
		Assertions []Assertion
	}

	j := []s{}
	d.Map(func(spanQuery SpanQuery, asserts []Assertion) {
		j = append(j, s{string(spanQuery), asserts})
	})

	return json.Marshal(j)
}

func (d *Definition) UnmarshalJSON(data []byte) error {
	aux := []struct {
		Selector   string
		Assertions []Assertion
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	newD := Definition{}
	var err error
	for _, s := range aux {
		newD, err = newD.Add(SpanQuery(s.Selector), s.Assertions)
		if err != nil {
			return err
		}
	}

	d.replace(&newD)

	return nil
}

func (d Definition) Add(key SpanQuery, asserts []Assertion) (Definition, error) {
	if d.keyPosition == nil {
		d.keyPosition = make(map[SpanQuery]int)
	}
	if d.positionKey == nil {
		d.positionKey = make(map[int]SpanQuery)
	}

	if _, exists := d.keyPosition[key]; exists {
		return Definition{}, errors.New("selector already exists")
	}

	d.list = append(d.list, asserts)
	ix := len(d.list) - 1
	d.keyPosition[key] = ix
	d.positionKey[ix] = key

	return d, nil
}

func (d Definition) Len() int {
	return len(d.list)
}

func (d Definition) Get(spanQuery SpanQuery) []Assertion {
	ix, exists := d.keyPosition[spanQuery]
	if !exists {
		return nil
	}

	return d.list[ix]
}

type mapFn func(spanQuery SpanQuery, asserts []Assertion)

func (d *Definition) Map(fn mapFn) {
	for ix, asserts := range d.list {
		spanQuery := d.positionKey[ix]
		fn(spanQuery, asserts)
	}
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
