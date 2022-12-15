package executor

import (
	"context"

	"github.com/kubeshop/tracetest/server/assertions"
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
	return assertions.Assert(defs, trace, ds)
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
