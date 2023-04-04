package executor

import (
	"context"

	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AssertionExecutor interface {
	Assert(context.Context, model.OrderedMap[model.SpanQuery, model.NamedAssertions], model.Trace, []expression.DataStore) (model.OrderedMap[model.SpanQuery, []model.AssertionResult], bool)
}

type defaultAssertionExecutor struct{}

func (e defaultAssertionExecutor) Assert(ctx context.Context, defs model.OrderedMap[model.SpanQuery, model.NamedAssertions], trace model.Trace, ds []expression.DataStore) (model.OrderedMap[model.SpanQuery, []model.AssertionResult], bool) {
	testResult := model.OrderedMap[model.SpanQuery, []model.AssertionResult]{}
	allPassed := true
	defs.ForEach(func(spanQuery model.SpanQuery, asserts model.NamedAssertions) error {
		spans := selector(spanQuery).Filter(trace)
		assertionResults := make([]model.AssertionResult, 0)
		for _, assertion := range asserts.Assertions {
			res := e.assert(assertion, spans, ds)
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

func (e defaultAssertionExecutor) assert(assertion model.Assertion, spans model.Spans, ds []expression.DataStore) model.AssertionResult {
	ds = append([]expression.DataStore{
		expression.MetaAttributesDataStore{SelectedSpans: spans},
		expression.VariableDataStore{},
	}, ds...)

	allPassed := true
	spanResults := make([]model.SpanAssertionResult, 0, len(spans))
	spans.
		ForEach(func(_ int, span model.Span) bool {
			res := e.assertSpan(span, ds, string(assertion))
			spanResults = append(spanResults, res)

			if res.CompareErr != nil {
				allPassed = false
			}

			return true
		}).
		OrEmpty(func() {
			res := e.assertSpan(model.Span{}, ds, string(assertion))
			spanResults = append(spanResults, res)
			allPassed = res.CompareErr == nil
		})

	return model.AssertionResult{
		Assertion: assertion,
		AllPassed: allPassed,
		Results:   spanResults,
	}
}

func (e defaultAssertionExecutor) assertSpan(span model.Span, ds []expression.DataStore, assertion string) model.SpanAssertionResult {
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

type instrumentedAssertionExecutor struct {
	assertionExecutor AssertionExecutor
	tracer            trace.Tracer
}

func (e instrumentedAssertionExecutor) Assert(ctx context.Context, defs model.OrderedMap[model.SpanQuery, model.NamedAssertions], trace model.Trace, ds []expression.DataStore) (model.OrderedMap[model.SpanQuery, []model.AssertionResult], bool) {
	ctx, span := e.tracer.Start(ctx, "Execute assertions")
	defer span.End()

	result, allPassed := e.assertionExecutor.Assert(ctx, defs, trace, ds)
	span.SetAttributes(
		attribute.Bool("tracetest.run.assertion_runner.all_assertions_passed", allPassed),
	)

	return result, allPassed
}

func NewAssertionExecutor(tracer trace.Tracer) AssertionExecutor {
	return &instrumentedAssertionExecutor{
		assertionExecutor: defaultAssertionExecutor{},
		tracer:            tracer,
	}
}

func selector(sq model.SpanQuery) selectors.Selector {
	sel, _ := selectors.New(string(sq))
	return sel
}
