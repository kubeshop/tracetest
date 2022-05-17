package traces

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/trace"
	v11 "go.opentelemetry.io/proto/otlp/common/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func FromOtel(input *v1.TracesData) (Trace, error) {
	flattenSpans := make([]*v1.Span, 0)
	for _, resource := range input.ResourceSpans {
		for _, librarySpans := range resource.InstrumentationLibrarySpans {
			flattenSpans = append(flattenSpans, librarySpans.Spans...)
		}
	}

	spansMap := make(map[trace.SpanID]*Span, 0)
	for _, span := range flattenSpans {
		newSpan, err := convertOtelSpanIntoSpan(span)
		if err != nil {
			return Trace{}, err
		}
		spansMap[newSpan.ID] = newSpan
	}

	return createTrace(flattenSpans, spansMap), nil
}

func convertOtelSpanIntoSpan(span *v1.Span) (*Span, error) {
	attributes := make(Attributes, 0)
	for _, attribute := range span.Attributes {
		attributes[attribute.Key] = getAttributeValue(attribute.Value)
	}

	attributes["name"] = span.Name
	attributes["tracetest.span.type"] = spanType(attributes)
	attributes["tracetest.span.duration"] = spanDuration(span)

	spanID, err := createSpanID(string(span.SpanId))
	if err != nil {
		return nil, err
	}

	return &Span{
		ID:         spanID,
		Name:       span.Name,
		Parent:     nil,
		Children:   make([]*Span, 0),
		Attributes: attributes,
	}, nil
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

func createSpanID(id string) (trace.SpanID, error) {
	spanId, err := trace.SpanIDFromHex(id)
	if err != nil {
		return trace.SpanID{}, fmt.Errorf("could not convert spanID")
	}

	return spanId, nil
}

func createTrace(spans []*v1.Span, spansMap map[trace.SpanID]*Span) Trace {
	rootSpanID := trace.SpanID{}
	for _, span := range spans {
		spanID, _ := createSpanID(string(span.SpanId))
		if string(span.ParentSpanId) == "" {
			rootSpanID = spanID
		} else {
			parentID, _ := createSpanID(string(span.ParentSpanId))
			parent := spansMap[parentID]
			thisSpan := spansMap[spanID]

			thisSpan.Parent = parent
			parent.Children = append(parent.Children, thisSpan)
		}
	}

	rootSpan := spansMap[rootSpanID]

	return Trace{
		RootSpan: *rootSpan,
		Flat:     spansMap,
	}
}
