package traces

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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

	spansMap := map[trace.SpanID]*Span{}
	for _, span := range flattenSpans {
		newSpan := convertOtelSpanIntoSpan(span)
		spansMap[newSpan.ID] = newSpan
	}

	return createTrace(flattenSpans, spansMap)
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

	attributes["name"] = span.Name
	attributes["kind"] = span.Kind.String()
	attributes["tracetest.span.type"] = spanType(attributes)
	attributes["tracetest.span.duration"] = spanDuration(span)

	spanID := createSpanID(span.SpanId)
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

func spanDuration(span *v1.Span) string {
	if span.GetStartTimeUnixNano() != 0 && span.GetEndTimeUnixNano() != 0 {
		spanDuration := (span.GetEndTimeUnixNano() - span.GetStartTimeUnixNano()) / 1000 / 1000 // in milliseconds
		return strconv.FormatUint(spanDuration, 10)
	}

	return "0"
}

func spanType(attrs Attributes) string {
	// based on https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/trace/semantic_conventions
	// using the first required attribute for each type
	for key := range attrs {
		switch true {
		case strings.HasPrefix(key, "http."):
			return "http"
		case strings.HasPrefix(key, "db."):
			return "database"
		case strings.HasPrefix(key, "rpc."):
			return "rpc"
		case strings.HasPrefix(key, "messaging."):
			return "messaging"
		case strings.HasPrefix(key, "faas."):
			return "faas"
		case strings.HasPrefix(key, "exception."):
			return "exception"
		}
	}
	return "unknown"
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
	var sid [8]byte
	copy(sid[:], id[:8])

	return trace.SpanID(sid)
}

func createTrace(spans []*v1.Span, spansMap map[trace.SpanID]*Span) Trace {
	rootSpanID := trace.SpanID{}
	for _, span := range spans {
		spanID := createSpanID(span.SpanId)
		parentSpanID := createSpanID(span.ParentSpanId)
		parentSpan, hasParent := spansMap[parentSpanID]
		if !hasParent {
			rootSpanID = spanID
		} else {
			thisSpan := spansMap[spanID]
			thisSpan.Parent = parentSpan
			parentSpan.Children = append(parentSpan.Children, thisSpan)
		}
	}

	rootSpan := spansMap[rootSpanID]

	return Trace{
		RootSpan: *rootSpan,
		Flat:     spansMap,
	}
}
