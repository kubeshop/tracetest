package trigger

import (
	"context"

	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/attribute"
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
		attribute.String("tracetest.run.trigger.test_id", test.ID.String()),
		attribute.String("tracetest.run.trigger.type", string(t.triggerer.Type())),
	}

	for k, v := range resp.SpanAttributes {
		attrs = append(attrs, attribute.String(k, v))
	}

	span.SetAttributes(attrs...)

	return resp, err
}
