package models

import (
	"context"
	"net/http"

	"github.com/kubeshop/tracetest/extensions/k6/utils"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type PropagatorName string

type Propagator struct {
	propagators []PropagatorName
}

const (
	PropagatorW3C    PropagatorName = "w3c"
	HeaderNameW3C    PropagatorName = "traceparent"
	PropagatorB3     PropagatorName = "b3"
	PropagatorJaeger PropagatorName = "jaeger"
	HeaderNameJaeger PropagatorName = "uber-trace-id"
)

func NewPropagator(propagators []PropagatorName) Propagator {
	return Propagator{
		propagators: propagators,
	}
}

func (p Propagator) GenerateHeaders(traceID string) http.Header {
	ctx := context.Background()
	spanContext := NewSpanContext(traceID)
	ctx = trace.ContextWithSpanContext(ctx, spanContext)
	header := http.Header{}

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	propagators().Inject(ctx, propagation.HeaderCarrier(header))

	return header
}

func NewSpanContext(traceID string) trace.SpanContext {
	parsedTraceID, _ := trace.TraceIDFromHex(traceID)
	var tf trace.TraceFlags
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    parsedTraceID,
		SpanID:     utils.SpanID(),
		TraceFlags: tf.WithSampled(true),
		Remote:     true,
	})
}

func propagators() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.Baggage{},
		b3.New(),
		jaeger.Jaeger{},
		ot.OT{},
		xray.Propagator{},
		propagation.TraceContext{})
}
