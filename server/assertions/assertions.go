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
	if a.Attribute.IsMeta() {
		return assertMeta(a, spans)
	}

	return assertIndividualSpans(a, spans)
}

func assertMeta(a model.Assertion, spans []traces.Span) model.AssertionResult {
	res := func(res model.SpanAssertionResult, err error) model.AssertionResult {
		return model.AssertionResult{
			Assertion: a,
			AllPassed: err == nil,
			Results:   []model.SpanAssertionResult{res},
		}
	}

	ma, err := metaAttr(a.Attribute.String())
	if err != nil {
		return res(model.SpanAssertionResult{}, err)
	}

	expectedValue, err := ExecuteExpression(*a.Value, traces.Span{})
	actualValue := ma.Value(spans)

	sar := apply(a, expectedValue, actualValue, nil)
	if err != nil {
		sar.CompareErr = err
	}

	return res(sar, sar.CompareErr)
}

func assertIndividualSpans(a model.Assertion, spans []traces.Span) model.AssertionResult {
	res := make([]model.SpanAssertionResult, len(spans))
	allPassed := true
	for i, span := range spans {
		spanID := span.ID
		var err error = nil
		expectedValue, err := ExecuteExpression(*a.Value, span)
		actualValue := span.Attributes.Get(a.Attribute.String())

		res[i] = apply(
			a,
			expectedValue,
			actualValue,
			&spanID,
		)
		if err != nil {
			res[i].CompareErr = err
		}

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

func apply(a model.Assertion, expected, actual string, spanID *trace.SpanID) model.SpanAssertionResult {
	return model.SpanAssertionResult{
		SpanID:        spanID,
		ObservedValue: actual,
		CompareErr:    a.Comparator.Compare(expected, actual),
	}
}

func selector(sq model.SpanQuery) selectors.Selector {
	sel, _ := selectors.New(string(sq))
	return sel
}
