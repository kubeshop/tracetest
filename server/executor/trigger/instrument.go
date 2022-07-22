package trigger

import (
	"context"
	"io"

	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func Instrument(tracer trace.Tracer, wrapped Triggerer) Triggerer {
	return &instrumentedTriggerer{
		tracer:    tracer,
		triggerer: wrapped,
	}
}

type instrumentedTriggerer struct {
	tracer    trace.Tracer
	triggerer Triggerer
}

func (t *instrumentedTriggerer) Type() model.TriggerType {
	return model.TriggerType("instrumented")
}

func (t *instrumentedTriggerer) Trigger(ctx context.Context, test model.Test, tid trace.TraceID, sid trace.SpanID) (Response, error) {
	ctx, span := t.tracer.Start(ctx, "Trigger test")
	defer span.End()

	resp, err := t.triggerer.Trigger(ctx, test, tid, sid)

	attrs := []attribute.KeyValue{
		attribute.String("tracetest.run.trigger.trace_id", tid.String()),
		attribute.String("tracetest.run.trigger.span_id", sid.String()),
		attribute.String("tracetest.run.trigger.test_id", test.ID.String()),
		attribute.String("tracetest.run.trigger.type", string(t.triggerer.Type())),
	}

	if err != nil {
		span.RecordError(err)
		attrs = append(attrs, attribute.String("tracetest.run.trigger.error", err.Error()))
	}

	for k, v := range resp.SpanAttributes {
		attrs = append(attrs, attribute.String(k, v))
	}

	span.SetAttributes(attrs...)

	return resp, err
}

func propagators() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.Baggage{},
		b3.New(),
		jaeger.Jaeger{},
		ot.OT{},
		xray.Propagator{},
		propagation.TraceContext{})
}

func traceProvider() *sdktrace.TracerProvider {
	// Set standard attributes per semantic conventions
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("tracetest"),
	)

	// this is in fact a noop exporter, so we can ignore errors
	spanExporter, _ := stdouttrace.New(stdouttrace.WithWriter(io.Discard))

	return sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(spanExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	)
}
