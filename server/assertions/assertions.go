package assertions

import (
	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
)

func Assert(defs model.OrderedMap[model.SpanQuery, model.NamedAssertions], trace traces.Trace) (model.OrderedMap[model.SpanQuery, []model.AssertionResult], bool) {
	testResult := model.OrderedMap[model.SpanQuery, []model.AssertionResult]{}
	allPassed := true
	defs.Map(func(spanQuery model.SpanQuery, asserts model.NamedAssertions) bool {
		spans := selector(spanQuery).Filter(trace)
		assertionResults := make([]model.AssertionResult, 0)
		for _, assertion := range asserts.Assertions {
			res := assert(assertion, spans)
			if !res.AllPassed {
				allPassed = false
			}
			assertionResults = append(assertionResults, res)
		}
		testResult, _ = testResult.Add(spanQuery, assertionResults)
		return true
	})

	return testResult, allPassed
}

func assert(assertion model.Assertion, spans traces.Spans) model.AssertionResult {
	ds := []expression.DataStore{
		expression.MetaAttributesDataStore{SelectedSpans: spans},
		expression.VariableDataStore{},
	}
	allPassed := true
	spanResults := make([]model.SpanAssertionResult, 0, len(spans))
	spans.MapIfZeroItems(
		func() {
			res := assertSpan(traces.Span{}, ds, string(assertion))
			spanResults = append(spanResults, res)
			allPassed = res.CompareErr == nil
		},
		func(_ int, span traces.Span) bool {
			res := assertSpan(span, ds, string(assertion))
			spanResults = append(spanResults, res)

			if res.CompareErr != nil {
				allPassed = false
			}

			return true
		},
	)

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
