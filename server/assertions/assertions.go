package assertions

import (
	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/assertions/selectors"
	"github.com/kubeshop/tracetest/traces"
)

type SpanQuery string

func (sq SpanQuery) Selector() selectors.Selector {
	sel, _ := selectors.New(string(sq))
	return sel
}

type TestDefinition map[SpanQuery][]Assertion

type Assertion struct {
	ID         string
	Attribute  string
	Comparator comparator.Comparator
	Value      string
}

func (a Assertion) Assert(spans []traces.Span) AssertionResult {
	results := make([]AssertionSpanResults, len(spans))
	for i, span := range spans {
		results[i] = a.apply(span)
	}
	return AssertionResult{
		Assertion:            a,
		AssertionSpanResults: results,
	}
}

func (a Assertion) apply(span traces.Span) AssertionSpanResults {
	attr := span.Attributes.Get(a.Attribute)
	return AssertionSpanResults{
		Span:        &span,
		ActualValue: attr,
		CompareErr:  a.Comparator.Compare(attr, a.Value),
	}
}

type AssertionResult struct {
	Assertion
	AssertionSpanResults []AssertionSpanResults
}

type AssertionSpanResults struct {
	Span        *traces.Span
	ActualValue string
	CompareErr  error
}

type TestResult map[SpanQuery][]AssertionResult

func Assert(trace traces.Trace, defs TestDefinition) TestResult {
	testResult := TestResult{}
	for spanQuery, asserts := range defs {
		spans := spanQuery.Selector().Filter(trace)
		assertionResults := make([]AssertionResult, 0)
		for _, assertion := range asserts {
			assertionResults = append(assertionResults, assertion.Assert(spans))
		}
		testResult[spanQuery] = assertionResults
	}

	return testResult
}
