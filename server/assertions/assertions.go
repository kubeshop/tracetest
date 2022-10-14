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
	defs.Map(func(spanQuery model.SpanQuery, asserts model.NamedAssertions) {
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
	})

	return testResult, allPassed
}

func assert(assertion model.Assertion, spans []traces.Span) model.AssertionResult {
	metaAttributeDataStore := expression.MetaAttributesDataStore{SelectedSpans: spans}
	if len(spans) == 0 {
		// we still have to run the assertion because it migth use a meta attribute
		emptyAttributeDataStore := expression.AttributeDataStore{}
		emptyVariablesDataStore := expression.VariableDataStore{}
		executor := expression.NewExecutor(emptyAttributeDataStore, metaAttributeDataStore, emptyVariablesDataStore)
		result := runAssertion(executor, string(assertion))

		return model.AssertionResult{
			Assertion: assertion,
			AllPassed: result.CompareErr == nil,
			Results:   []model.SpanAssertionResult{result},
		}
	}

	allPassed := true
	spanResults := make([]model.SpanAssertionResult, 0, len(spans))
	for _, span := range spans {
		spanID := span.ID
		attributeDataStore := expression.AttributeDataStore{Span: span}

		// TODO: populate it with variables when they are available
		variablesDataStore := expression.VariableDataStore{}
		expressionExecutor := expression.NewExecutor(attributeDataStore, metaAttributeDataStore, variablesDataStore)

		actualValue, _, err := expressionExecutor.ExecuteStatement(string(assertion))
		if err != nil {
			allPassed = false
		}

		spanResults = append(spanResults, model.SpanAssertionResult{
			SpanID:        &spanID,
			ObservedValue: actualValue,
			CompareErr:    err,
		})
	}

	return model.AssertionResult{
		Assertion: assertion,
		AllPassed: allPassed,
		Results:   spanResults,
	}
}

func runAssertion(executor expression.Executor, assertion string) model.SpanAssertionResult {
	actualValue, _, err := executor.ExecuteStatement(string(assertion))
	return model.SpanAssertionResult{
		SpanID:        nil,
		ObservedValue: actualValue,
		CompareErr:    err,
	}
}

func selector(sq model.SpanQuery) selectors.Selector {
	sel, _ := selectors.New(string(sq))
	return sel
}
