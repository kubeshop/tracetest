package lightstep

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

func traceId(id string) trace.TraceID {
	traceId, _ := trace.TraceIDFromHex(id)
	return traceId
}

func spanId(id string) trace.SpanID {
	spanId, _ := trace.SpanIDFromHex(id)
	return spanId
}

func ConvertResponseToOtelFormat(response GetTraceResponse) traces.Trace {
	if len(response.Data) == 0 {
		return traces.Trace{}
	}

	data := response.Data[0]

	return traces.Trace{
		ID:       traceId(data.ID),
		RootSpan: *getSpanTree(data),
	}
}

func getSpanTree(data traceData) *traces.Span {
	var rootSpan *traces.Span
	spans := make(map[string]*traces.Span, len(data.Attributes.Spans))
	for _, span := range data.Attributes.Spans {
		traceSpan, isRootSpan := convertSpan(span)
		spans[traceSpan.ID.String()] = &traceSpan
		if isRootSpan {
			rootSpan = &traceSpan
		}
	}

	buildSpanTree(spans)

	return rootSpan
}

func convertSpan(span span) (traces.Span, bool) {
	attributes := make(map[string]string, len(span.Tags))
	for key, value := range span.Tags {
		attributes[key] = fmt.Sprintf("%v", value)
	}

	_, hasParent := attributes["parent_span_guid"]
	isRootSpan := !hasParent

	return traces.Span{
		ID:         spanId(span.SpanID),
		Name:       span.SpanName,
		StartTime:  time.UnixMicro(span.StartTimeMicros),
		EndTime:    time.UnixMicro(span.EndTimeMicros),
		Attributes: attributes,
	}, isRootSpan
}

func buildSpanTree(spans map[string]*traces.Span) {
	for _, span := range spans {
		parentId := span.Attributes["parent_span_guid"]
		if parentId == "" {
			continue
		}

		parent := spans[parentId]
		parent.Children = append(parent.Children, span)
		for _, child := range parent.Children {
			child.Parent = parent
		}
	}
}
