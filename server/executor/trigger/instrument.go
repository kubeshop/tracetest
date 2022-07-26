package trigger

import (
	"context"

	"github.com/kubeshop/tracetest/server/model"
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

func (t *instrumentedTriggerer) Type() model.TriggerType {
	return model.TriggerType("instrumented")
}

func (t *instrumentedTriggerer) Trigger(ctx context.Context, triggerSpanCtx context.Context, test model.Test, tid trace.TraceID, sid trace.SpanID) (Response, error) {
	ctx, span := t.tracer.Start(ctx, "Trigger test")
	defer span.End()

	resp, err := t.triggerer.Trigger(ctx, triggerSpanCtx, test, tid, sid)

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
