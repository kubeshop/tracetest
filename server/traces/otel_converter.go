package traces

import (
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/trace"
	v11 "go.opentelemetry.io/proto/otlp/common/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func FromOtel(input *v1.TracesData) model.Trace {
	return fromOtelResourceSpans(input.ResourceSpans)
}

func fromOtelResourceSpans(resourceSpans []*v1.ResourceSpans) model.Trace {
	flattenSpans := make([]*v1.Span, 0)
	for _, resource := range resourceSpans {
		for _, scopeSpans := range resource.ScopeSpans {
			flattenSpans = append(flattenSpans, scopeSpans.Spans...)
		}
	}

	traceID := ""
	spans := make([]model.Span, 0)
	for _, span := range flattenSpans {
		newSpan := ConvertOtelSpanIntoSpan(span)
		traceID = hex.EncodeToString(span.TraceId)
		spans = append(spans, *newSpan)
	}

	return model.NewTrace(traceID, spans)
}

func ConvertOtelSpanIntoSpan(span *v1.Span) *model.Span {
	attributes := make(model.Attributes, 0)
	for _, attribute := range span.Attributes {
		attributes[attribute.Key] = getAttributeValue(attribute.Value)
	}

	var startTime, endTime time.Time

	if span.GetStartTimeUnixNano() != 0 {
		startTime = time.Unix(0, int64(span.GetStartTimeUnixNano()))
	}

	if span.GetEndTimeUnixNano() != 0 {
		endTime = time.Unix(0, int64(span.GetEndTimeUnixNano()))
	}

	spanID := createSpanID(span.SpanId)
	attributes["parent_id"] = createSpanID(span.ParentSpanId).String()
	return &model.Span{
		ID:         spanID,
		Name:       span.Name,
		Kind:       spanKind(span),
		StartTime:  startTime,
		EndTime:    endTime,
		Parent:     nil,
		Children:   make([]*model.Span, 0),
		Attributes: attributes,
	}
}

func spanKind(span *v1.Span) model.SpanKind {
	switch span.Kind {
	case v1.Span_SPAN_KIND_CLIENT:
		return model.SpanKindClient
	case v1.Span_SPAN_KIND_SERVER:
		return model.SpanKindServer
	case v1.Span_SPAN_KIND_INTERNAL:
		return model.SpanKindInternal
	case v1.Span_SPAN_KIND_PRODUCER:
		return model.SpanKindProducer
	case v1.Span_SPAN_KIND_CONSUMER:
		return model.SpanKindConsumer
	default:
		return model.SpanKindUnespecified
	}
}

func getAttributeValue(value *v11.AnyValue) string {
	switch v := value.GetValue().(type) {
	case *v11.AnyValue_StringValue:
		return v.StringValue

	case *v11.AnyValue_IntValue:
		return fmt.Sprintf("%d", v.IntValue)

	case *v11.AnyValue_DoubleValue:
		if v.DoubleValue != 0.0 {
			isFloatingPoint := math.Abs(v.DoubleValue-math.Abs(v.DoubleValue)) > 0.0
			if isFloatingPoint {
				return fmt.Sprintf("%f", v.DoubleValue)
			}

			return fmt.Sprintf("%.0f", v.DoubleValue)
		}

	case *v11.AnyValue_BoolValue:
		return fmt.Sprintf("%t", v.BoolValue)
	}

	return "unsupported value type"
}

func createSpanID(id []byte) trace.SpanID {
	if id == nil {
		return trace.SpanID{}
	}

	var sid [8]byte
	copy(sid[:], id[:8])

	return trace.SpanID(sid)
}

func CreateTraceID(id []byte) trace.TraceID {
	if id == nil {
		return trace.TraceID{}
	}

	var tid [16]byte
	copy(tid[:], id[:16])
	return trace.TraceID(tid)
}

func DecodeTraceID(id string) trace.TraceID {
	bytes, _ := hex.DecodeString(id)
	var tid [16]byte
	copy(tid[:], bytes[:16])
	return trace.TraceID(tid)
}
