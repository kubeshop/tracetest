package trigger

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
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

func (t *instrumentedTriggerer) Type() TriggerType {
	return TriggerType("instrumented")
}

func (t *instrumentedTriggerer) Trigger(ctx context.Context, triggerConfig Trigger, opts *Options) (Response, error) {
	_, span := t.tracer.Start(ctx, "Trigger test")
	defer span.End()

	tracestate, err := trace.ParseTraceState("tracetest=true")
	if err != nil {
		return Response{}, fmt.Errorf("could not create tracestate: %w", err)
	}

	var tf trace.TraceFlags
	spanContextConfig := trace.SpanContextConfig{
		TraceState: tracestate,
		Remote:     true,
		TraceFlags: tf.WithSampled(true),
	}

	if opts != nil {
		spanContextConfig.TraceID = opts.TraceID
		spanContextConfig.SpanID = opts.SpanID
	}

	spanContext := trace.NewSpanContext(spanContextConfig)

	triggerCtx := trace.ContextWithSpanContext(context.Background(), spanContext)

	tid := spanContext.TraceID()

	resp, err := t.triggerer.Trigger(triggerCtx, triggerConfig, opts)

	resp.TraceID = tid

	attrs := []attribute.KeyValue{
		attribute.String("tracetest.run.trace_id", tid.String()),
		attribute.String("tracetest.run.test_id", string(opts.TestID)),
		attribute.String("tracetest.run.type", string(t.triggerer.Type())),
	}

	if err != nil {
		span.RecordError(err)
		attrs = append(attrs, attribute.String("tracetest.run.error", err.Error()))
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
