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

	return traces.NewTrace(traceId(data.ID), *getSpanTree(data))
}

func getSpanTree(data traceData) *traces.Span {
	var rootSpan *traces.Span
	spans := make(map[string]*traces.Span, len(data.Attributes.Spans))
	reporters := data.Relationships.Reporters
	for _, span := range data.Attributes.Spans {
		traceSpan, isRootSpan := convertSpan(span, reporters)
		spans[traceSpan.ID.String()] = &traceSpan
		if isRootSpan {
			rootSpan = &traceSpan
		}
	}

	buildSpanTree(spans)

	return rootSpan
}

func convertSpan(span span, reporters []reporter) (traces.Span, bool) {
	attributes := make(map[string]string)
	for key, value := range span.Tags {
		attributes[key] = fmt.Sprintf("%v", value)
	}

	reporterAttributes := getReporterAttributes(reporters, span.ReporterID)
	for key, value := range reporterAttributes {
		attributes[key] = fmt.Sprintf("%v", value)
	}

	_, hasParent := attributes["parent_span_guid"]
	isRootSpan := !hasParent

	attributes["kind"] = attributes["span.kind"]

	return traces.Span{
		ID:         spanId(span.SpanID),
		Name:       span.SpanName,
		StartTime:  time.Unix(0, span.StartTimeMicros*1000),
		EndTime:    time.Unix(0, span.EndTimeMicros*1000),
		Attributes: attributes,
	}, isRootSpan
}

func getReporterAttributes(reporters []reporter, id string) map[string]interface{} {
	for _, reporter := range reporters {
		if reporter.ReporterID == id {
			return reporter.Attributes
		}
	}

	return map[string]interface{}{}
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
