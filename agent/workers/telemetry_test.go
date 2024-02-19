package workers_test

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

var propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

func ContextWithTracingEnabled() context.Context {
	ctx, span := trace.NewTracerProvider().Tracer("tracer").Start(context.Background(), "root span")
	defer span.End()

	return ctx
}

func SameTraceID(ctx1, ctx2 context.Context) bool {
	header1 := make(propagation.HeaderCarrier)
	header2 := make(propagation.HeaderCarrier)

	propagator.Inject(ctx1, header1)
	propagator.Inject(ctx2, header2)

	return header1.Get("traceparent") == header2.Get("traceparent")
}
