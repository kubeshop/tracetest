package main

import (
	"fmt"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

type Trace struct {
	ID       trace.TraceID
	RootSpan Span
}

func (t Trace) FindSpans(selector Selector) []Span {
	return []Span{}
}

type Attributes map[string]string

func (a Attributes) Get(key string) string {
	if v, ok := a[key]; ok {
		return v
	}

	return ""
}

type Span struct {
	ID         trace.SpanID
	Name       string
	Attributes Attributes
	parent     *Span
	children   []*Span
}

// ******************

var ErrNoMatch = fmt.Errorf("no match")

type Selector string // TODO: how are we actually gonna represent this?

type ComparatorFn func(string, string) error

type Assertion struct {
	Attribute  string
	Comparator ComparatorFn
	Value      string
}

func (a Assertion) Check(spans []Span) AssertionResult {
	results := make([]AssertionSpanResults, len(spans))
	for i, span := range spans {
		results[i] = a.apply(span)
	}
	return AssertionResult{
		Assertion:            a,
		AssertionSpanResults: results,
	}
}

func (a Assertion) apply(span Span) AssertionSpanResults {
	attr := span.Attributes.Get(a.Attribute)
	return AssertionSpanResults{
		Span:        &span,
		ActualValue: attr,
		CompareErr:  a.Comparator(a.Value, attr),
	}
}

type AssertionResult struct {
	Assertion
	AssertionSpanResults []AssertionSpanResults
}

type AssertionSpanResults struct {
	Span        *Span
	ActualValue string
	CompareErr  error
}

type TestDefinition map[Selector][]Assertion

type TestResult map[Selector]AssertionResult

func RunAssertions(trace Trace, defs TestDefinition) TestResult {
	testResult := TestResult{}
	for selector, asserts := range defs {
		spans := trace.FindSpans(selector)
		for _, assertion := range asserts {
			testResult[selector] = assertion.Check(spans)
		}
	}
	return testResult
}

func ComparatorEq(expected, actual string) error {
	if expected == actual {
		return nil
	}

	return ErrNoMatch
}

func ComparatorLt(expected, actual string) error {
	if expected < actual {
		return nil
	}

	return ErrNoMatch
}

func ComparatorGt(expected, actual string) error {
	if expected > actual {
		return nil
	}

	return ErrNoMatch
}

func ComparatorContains(expected, actual string) error {
	if strings.Contains(actual, expected) {
		return nil
	}

	return ErrNoMatch
}

func ComparatorIsDivisibleBy(expected, actual string) error {
	expectedInt, err := strconv.Atoi(expected)
	if err != nil {
		return fmt.Errorf("Expected value \"%s\" is not an int", expected)
	}

	actualInt, err := strconv.Atoi(actual)
	if err != nil {
		return fmt.Errorf("Actual value \"%s\" is not an int", actual)
	}

	if expectedInt%actualInt == 0 {
		return nil
	}

	return ErrNoMatch
}

// actual should be a json, like {"key": "value"}
func ComparatorJSonValue(expected, actual string) bool {
	return false
}

type versionedResult struct {
	ID            string
	TestID        string
	TraceID       string // db ID referenceØ
	AssertResults AssertionResult
}

var sample = TestDefinition{
	"[tracetest.span.type=http,service.name='pokeshop']": []Assertion{
		{
			Attribute:  "tracetest.span.duration",
			Comparator: ComparatorLt,
			Value:      "200", //ms
		},
		{
			Attribute:  "tracetest.span.duration",
			Comparator: ComparatorLt,
			Value:      "200", //ms
		},
	},
}

// ."Pokeshop" > * >[service.type="db"]
// [service.type="db"]
// span[name="Pokeshop", service.type="db"]:nth_child(2)
// span[name="Pokeshop"] span[service.type="http"]

// span[type="http"],span[service.type="rpc"]

// Pokeshop
//  ---> HTTP POST /import
//      ---> db.insert
//      ---> HTTP POST https://blhablah.com/send-pokemon
//          ---> db.insert
//  ---> PokeStore
//    ---> HTTP POST /buy
//        ---> db.insert
//        ---> HTTP POST https://blhablah.com/send-pokemon
//            ---> db.insert

// tracetest.span.type = "http"
// service.name = "pokeshop"
