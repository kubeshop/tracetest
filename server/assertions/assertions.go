package assertions

import (
	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
)

func Assert(defs model.OrderedMap[model.SpanQuery, model.NamedAssertions], trace traces.Trace, ds []expression.DataStore) (model.OrderedMap[model.SpanQuery, []model.AssertionResult], bool) {
	testResult := model.OrderedMap[model.SpanQuery, []model.AssertionResult]{}
	allPassed := true
	defs.ForEach(func(spanQuery model.SpanQuery, asserts model.NamedAssertions) error {
		spans := selector(spanQuery).Filter(trace)
		assertionResults := make([]model.AssertionResult, 0)
		for _, assertion := range asserts.Assertions {
			res := assert(assertion, spans, ds)
			if !res.AllPassed {
				allPassed = false
			}
			assertionResults = append(assertionResults, res)
		}
		testResult, _ = testResult.Add(spanQuery, assertionResults)
		return nil
	})

	return testResult, allPassed
}

func assert(assertion model.Assertion, spans traces.Spans, ds []expression.DataStore) model.AssertionResult {
	ds = append([]expression.DataStore{
		expression.MetaAttributesDataStore{SelectedSpans: spans},
		expression.VariableDataStore{},
	}, ds...)

	allPassed := true
	spanResults := make([]model.SpanAssertionResult, 0, len(spans))
	spans.
		ForEach(func(_ int, span traces.Span) bool {
			res := assertSpan(span, ds, string(assertion))
			spanResults = append(spanResults, res)

			if res.CompareErr != nil {
				allPassed = false
			}

			return true
		}).
		OrEmpty(func() {
			res := assertSpan(traces.Span{}, ds, string(assertion))
			spanResults = append(spanResults, res)
			allPassed = res.CompareErr == nil
		})

	return model.AssertionResult{
		Assertion: assertion,
		AllPassed: allPassed,
		Results:   spanResults,
	}
}

func assertSpan(span traces.Span, ds []expression.DataStore, assertion string) model.SpanAssertionResult {
	ds = append([]expression.DataStore{expression.AttributeDataStore{Span: span}}, ds...)
	expressionExecutor := expression.NewExecutor(ds...)

	actualValue, _, err := expressionExecutor.Statement(assertion)

	sar := model.SpanAssertionResult{
		ObservedValue: actualValue,
		CompareErr:    err,
	}

	if span.ID.IsValid() {
		sar.SpanID = &span.ID
	}

	return sar
}

func selector(sq model.SpanQuery) selectors.Selector {
	sel, _ := selectors.New(string(sq))
	return sel
}
