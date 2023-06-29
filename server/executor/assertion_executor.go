package executor

import (
	"context"

	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/test"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AssertionExecutor interface {
	Assert(context.Context, test.Specs, model.Trace, []expression.DataStore) (maps.Ordered[test.SpanQuery, []test.AssertionResult], bool)
}

type defaultAssertionExecutor struct{}

func (e defaultAssertionExecutor) Assert(_ context.Context, specs test.Specs, trace model.Trace, ds []expression.DataStore) (maps.Ordered[test.SpanQuery, []test.AssertionResult], bool) {
	testResult := maps.Ordered[test.SpanQuery, []test.AssertionResult]{}
	allPassed := true
	for _, spec := range specs {
		spans := selector(spec.Selector.Query).Filter(trace)
		assertionResults := make([]test.AssertionResult, 0)
		for _, assertion := range spec.Assertions {
			res := e.assert(assertion, spans, ds)
			if !res.AllPassed {
				allPassed = false
			}
			assertionResults = append(assertionResults, res)
		}
		testResult, _ = testResult.Add(spec.Selector.Query, assertionResults)
	}

	return testResult, allPassed
}

func (e defaultAssertionExecutor) assert(assertion test.Assertion, spans model.Spans, ds []expression.DataStore) test.AssertionResult {
	ds = append([]expression.DataStore{
		expression.MetaAttributesDataStore{SelectedSpans: spans},
		expression.VariableDataStore{},
	}, ds...)

	allPassed := true
	spanResults := make([]test.SpanAssertionResult, 0, len(spans))
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

	return test.AssertionResult{
		Assertion: assertion,
		AllPassed: allPassed,
		Results:   spanResults,
	}
}

func (e defaultAssertionExecutor) assertSpan(span model.Span, ds []expression.DataStore, assertion string) test.SpanAssertionResult {
	ds = append([]expression.DataStore{expression.AttributeDataStore{Span: span}}, ds...)
	expressionExecutor := expression.NewExecutor(ds...)

	actualValue, _, err := expressionExecutor.Statement(assertion)

	sar := test.SpanAssertionResult{
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

func (e instrumentedAssertionExecutor) Assert(ctx context.Context, defs test.Specs, trace model.Trace, ds []expression.DataStore) (maps.Ordered[test.SpanQuery, []test.AssertionResult], bool) {
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

func selector(sq test.SpanQuery) selectors.Selector {
	sel, _ := selectors.New(string(sq))
	return sel
}
