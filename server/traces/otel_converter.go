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
	flattenSpans := make([]*v1.Span, 0)
	for _, resource := range input.ResourceSpans {
		for _, librarySpans := range resource.InstrumentationLibrarySpans {
			flattenSpans = append(flattenSpans, librarySpans.Spans...)
		}
	}

	traceID := ""
	spans := make([]Span, 0)
	for _, span := range flattenSpans {
		newSpan := convertOtelSpanIntoSpan(span)
		traceID = hex.EncodeToString(span.TraceId)
		spans = append(spans, *newSpan)
	}

	return New(traceID, spans)
}

func convertOtelSpanIntoSpan(span *v1.Span) *Span {
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

	spanID := createSpanID(span.SpanId)
	attributes["parent_id"] = createSpanID(span.ParentSpanId).String()
	return &Span{
		ID:         spanID,
		Name:       span.Name,
		StartTime:  startTime,
		EndTime:    endTime,
		Parent:     nil,
		Children:   make([]*Span, 0),
		Attributes: attributes,
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
