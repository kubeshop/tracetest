package processors

import (
	"go.opentelemetry.io/otel/trace"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type Processor interface {
	Process(*pb.ExportTraceServiceRequest) *pb.ExportTraceServiceRequest
}

type datadogProcessor struct{}

func NewDatadogProcessor() Processor {
	return &datadogProcessor{}
}

func (p *datadogProcessor) Process(request *pb.ExportTraceServiceRequest) *pb.ExportTraceServiceRequest {
	request.ResourceSpans = fixDatadogTraceID(request.ResourceSpans)
	return request
}

func fixDatadogTraceID(in []*v1.ResourceSpans) []*v1.ResourceSpans {
	resourceSpans := make([]*v1.ResourceSpans, 0, len(in))
	for _, resourceSpan := range in {
		scopeSpans := make([]*v1.ScopeSpans, 0, len(resourceSpan.ScopeSpans))
		for _, scopeSpan := range resourceSpan.ScopeSpans {
			spans := make([]*v1.Span, 0, len(scopeSpan.Spans))
			for _, span := range scopeSpan.Spans {
				span.TraceId = getTraceIDFromDatadogSpan(span)
				spans = append(spans, span)
			}
			scopeSpan.Spans = spans
			scopeSpans = append(scopeSpans, scopeSpan)
		}

		resourceSpan.ScopeSpans = scopeSpans
		resourceSpans = append(resourceSpans, resourceSpan)
	}

	return resourceSpans
}

func getTraceIDFromDatadogSpan(span *v1.Span) []byte {
	firstHalfTraceID := ""
	for _, attr := range span.Attributes {
		if attr.Key == "_dd.p.tid" {
			firstHalfTraceID = attr.Value.GetStringValue()
		}
	}

	if len(firstHalfTraceID) == 0 {
		// not a datadog span
		return span.TraceId
	}

	filledSecondHalf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	filledSecondHalf = append(filledSecondHalf, span.TraceId...)

	secondHalfTraceID := trace.TraceID(filledSecondHalf).String()[16:]

	traceID := firstHalfTraceID + secondHalfTraceID

	realTraceID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return span.TraceId
	}

	return realTraceID[:]
}
