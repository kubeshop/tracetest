package assertions

import (
	"github.com/kubeshop/tracetest/assertions/selectors"
	"github.com/kubeshop/tracetest/model"
	"github.com/kubeshop/tracetest/traces"
	"go.opentelemetry.io/otel/trace"
)

func Assert(defs model.Definition, trace traces.Trace) model.Results {
	testResult := model.Results{}
	for spanQuery, asserts := range defs {
		spans := selector(spanQuery).Filter(trace)
		assertionResults := make([]model.AssertionResult, 0)
		for _, assertion := range asserts {
			assertionResults = append(assertionResults, assert(assertion, spans))
		}
		testResult[spanQuery] = assertionResults
	}

	return testResult
}

func assert(a model.Assertion, spans []traces.Span) model.AssertionResult {
	res := make([]model.SpanAssertionResult, len(spans))
	for i, span := range spans {
		res[i] = apply(a, span)
	}

	return model.AssertionResult{
		Assertion: a,
		Results:   res,
	}
}

func apply(a model.Assertion, span traces.Span) model.SpanAssertionResult {
	attr := span.Attributes.Get(a.Attribute)
	return model.SpanAssertionResult{
		SpanID:        trace.SpanID(span.ID),
		ObservedValue: attr,
		CompareErr:    a.Comparator.Compare(attr, a.Value),
	}
}

func selector(sq model.SpanQuery) selectors.Selector {
	sel, _ := selectors.New(string(sq))
	return sel
}
