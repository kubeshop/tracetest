package assertions

import (
	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

func Assert(defs model.OrderedMap[model.SpanQuery, []model.Assertion], trace traces.Trace) (model.OrderedMap[model.SpanQuery, []model.AssertionResult], bool) {
	testResult := model.OrderedMap[model.SpanQuery, []model.AssertionResult]{}
	allPassed := true
	defs.Map(func(spanQuery model.SpanQuery, asserts []model.Assertion) {
		spans := selector(spanQuery).Filter(trace)
		assertionResults := make([]model.AssertionResult, 0)
		for _, assertion := range asserts {
			res := assert(assertion, spans)
			if !res.AllPassed {
				allPassed = false
			}
			assertionResults = append(assertionResults, res)
		}
		testResult, _ = testResult.Add(spanQuery, assertionResults)
	})

	return testResult, allPassed
}

func assert(a model.Assertion, spans []traces.Span) model.AssertionResult {
	res := make([]model.SpanAssertionResult, len(spans))
	allPassed := true
	for i, span := range spans {
		res[i] = apply(a, span)
		if res[i].CompareErr != nil {
			allPassed = false
		}
	}

	return model.AssertionResult{
		Assertion: a,
		AllPassed: allPassed,
		Results:   res,
	}
}

func apply(a model.Assertion, span traces.Span) model.SpanAssertionResult {
	expected := a.Value
	actual := span.Attributes.Get(a.Attribute.String())
	sid := trace.SpanID(span.ID)
	return model.SpanAssertionResult{
		SpanID:        &sid,
		ObservedValue: actual,
		CompareErr:    a.Comparator.Compare(expected, actual),
	}
}

func selector(sq model.SpanQuery) selectors.Selector {
	sel, _ := selectors.New(string(sq))
	return sel
}
