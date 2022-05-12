package traces

import (
	"fmt"
	"math"

	"github.com/kubeshop/tracetest/openapi"
	"go.opentelemetry.io/otel/trace"
)

func FromOtel(trace openapi.ApiV3SpansResponseChunk) (Trace, error) {
	flattenSpans := make([]openapi.V1Span, 0)
	for _, resource := range trace.ResourceSpans {
		for _, librarySpans := range resource.InstrumentationLibrarySpans {
			for _, span := range librarySpans.Spans {
				flattenSpans = append(flattenSpans, span)
			}
		}
	}

	spansMap := make(map[string]*Span, 0)
	for _, span := range flattenSpans {
		newSpan, err := convertOtelSpanIntoSpan(span)
		if err != nil {
			return Trace{}, err
		}
		spansMap[span.SpanId] = newSpan
	}

	return createTrace(flattenSpans, spansMap), nil
}

func convertOtelSpanIntoSpan(span openapi.V1Span) (*Span, error) {
	attributes := make(Attributes, 0)
	for _, attribute := range span.Attributes {
		attributes[attribute.Key] = getAttributeValue(attribute.Value)
	}

	attributes["name"] = span.Name

	spanID, err := createSpanID(span.SpanId)
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

func getAttributeValue(value openapi.V1AnyValue) string {
	if value.StringValue != "" {
		return value.StringValue
	}

	if value.IntValue != "" {
		return value.IntValue
	}

	if value.DoubleValue != 0.0 {
		isFloatingPoint := math.Abs(value.DoubleValue-math.Abs(value.DoubleValue)) > 0.0
		if isFloatingPoint {
			return fmt.Sprintf("%f", value.DoubleValue)
		}

		return fmt.Sprintf("%.0f", value.DoubleValue)
	}

	return fmt.Sprintf("%t", value.BoolValue)
}

func createSpanID(id string) (trace.SpanID, error) {
	spanId, err := trace.SpanIDFromHex(id)
	if err != nil {
		return trace.SpanID{}, fmt.Errorf("could not convert spanID")
	}

	return spanId, nil
}

func createTrace(spans []openapi.V1Span, spansMap map[string]*Span) Trace {
	var rootSpanID string = ""
	for _, span := range spans {
		if span.ParentSpanId == "" {
			rootSpanID = span.SpanId
		} else {
			parent := spansMap[span.ParentSpanId]
			thisSpan := spansMap[span.SpanId]

			thisSpan.Parent = parent
			parent.Children = append(parent.Children, thisSpan)
		}
	}

	rootSpan := spansMap[rootSpanID]

	return Trace{
		RootSpan: *rootSpan,
	}
}
