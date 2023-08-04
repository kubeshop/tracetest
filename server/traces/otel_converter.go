package traces

import (
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"go.opentelemetry.io/otel/trace"
	v11 "go.opentelemetry.io/proto/otlp/common/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func FromOtel(input *v1.TracesData) Trace {
	return FromOtelResourceSpans(input.ResourceSpans)
}

func FromOtelResourceSpans(resourceSpans []*v1.ResourceSpans) Trace {
	flattenSpans := make([]*v1.Span, 0)
	for _, resource := range resourceSpans {
		for _, scopeSpans := range resource.ScopeSpans {
			flattenSpans = append(flattenSpans, scopeSpans.Spans...)
		}
	}

	return FromSpanList(flattenSpans)
}

func FromSpanList(input []*v1.Span) Trace {
	traceID := ""
	spans := make([]Span, 0)
	for _, span := range input {
		newSpan := ConvertOtelSpanIntoSpan(span)
		traceID = CreateTraceID(span.TraceId).String()
		spans = append(spans, *newSpan)
	}

	return NewTrace(traceID, spans)
}

func ConvertOtelSpanIntoSpan(span *v1.Span) *Span {
	attributes := make(Attributes, 0)
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

	var spanStatus *SpanStatus
	if span.Status != nil {
		spanStatus = &SpanStatus{
			Code:        span.Status.Code.String(),
			Description: span.Status.Message,
		}
	}

	spanID := createSpanID(span.SpanId)
	attributes[TracetestMetadataFieldParentID] = createSpanID(span.ParentSpanId).String()
	return &Span{
		ID:         spanID,
		Name:       span.Name,
		Kind:       spanKind(span),
		StartTime:  startTime,
		EndTime:    endTime,
		Parent:     nil,
		Events:     extractEvents(span),
		Status:     spanStatus,
		Children:   make([]*Span, 0),
		Attributes: attributes,
	}
}

func extractEvents(v1 *v1.Span) []SpanEvent {
	output := make([]SpanEvent, 0, len(v1.Events))
	for _, v1Event := range v1.Events {
		attributes := make(Attributes, 0)
		for _, attribute := range v1Event.Attributes {
			attributes[attribute.Key] = getAttributeValue(attribute.Value)
		}
		var timestamp time.Time

		if v1Event.GetTimeUnixNano() != 0 {
			timestamp = time.Unix(0, int64(v1Event.GetTimeUnixNano()))
		}

		output = append(output, SpanEvent{
			Name:       v1Event.Name,
			Timestamp:  timestamp,
			Attributes: attributes,
		})
	}

	return output
}

func spanKind(span *v1.Span) SpanKind {
	switch span.Kind {
	case v1.Span_SPAN_KIND_CLIENT:
		return SpanKindClient
	case v1.Span_SPAN_KIND_SERVER:
		return SpanKindServer
	case v1.Span_SPAN_KIND_INTERNAL:
		return SpanKindInternal
	case v1.Span_SPAN_KIND_PRODUCER:
		return SpanKindProducer
	case v1.Span_SPAN_KIND_CONSUMER:
		return SpanKindConsumer
	default:
		return SpanKindUnespecified
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
