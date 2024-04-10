package processors

import (
	"encoding/binary"
	"fmt"
	"strconv"

	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type Processor interface {
	Process(*pb.ExportTraceServiceRequest) (*pb.ExportTraceServiceRequest, error)
}

type datadogProcessor struct{}

func NewDatadogProcessor() Processor {
	return &datadogProcessor{}
}

func (p *datadogProcessor) Process(request *pb.ExportTraceServiceRequest) (*pb.ExportTraceServiceRequest, error) {
	spans, err := fixDatadogTraceID(request.ResourceSpans)
	request.ResourceSpans = spans
	return request, err
}

func fixDatadogTraceID(in []*v1.ResourceSpans) ([]*v1.ResourceSpans, error) {
	traceIDCache := make(map[string][]byte)
	resourceSpans := make([]*v1.ResourceSpans, 0, len(in))
	for _, resourceSpan := range in {
		scopeSpans := make([]*v1.ScopeSpans, 0, len(resourceSpan.ScopeSpans))
		for _, scopeSpan := range resourceSpan.ScopeSpans {
			spans := make([]*v1.Span, 0, len(scopeSpan.Spans))
			for _, span := range scopeSpan.Spans {
				originalTraceID := string(span.TraceId)
				traceID, err := getTraceIDFromDatadogSpan(span)
				if err != nil {
					return in, fmt.Errorf("cannot get trace-id from datadog span: %w", err)
				}
				traceIDCache[originalTraceID] = traceID
				span.TraceId = traceID
				spans = append(spans, span)
			}
			scopeSpan.Spans = spans
			scopeSpans = append(scopeSpans, scopeSpan)
		}

		resourceSpan.ScopeSpans = scopeSpans
		resourceSpans = append(resourceSpans, resourceSpan)
	}

	for _, resourceSpan := range resourceSpans {
		scopeSpans := make([]*v1.ScopeSpans, 0, len(resourceSpan.ScopeSpans))
		for _, scopeSpan := range resourceSpan.ScopeSpans {
			spans := make([]*v1.Span, 0, len(scopeSpan.Spans))
			for _, span := range scopeSpan.Spans {
				if realParentID, ok := traceIDCache[string(span.ParentSpanId)]; ok {
					span.ParentSpanId = realParentID
				}

				spans = append(spans, span)
			}
			scopeSpan.Spans = spans
			scopeSpans = append(scopeSpans, scopeSpan)
		}

		resourceSpan.ScopeSpans = scopeSpans
		resourceSpans = append(resourceSpans, resourceSpan)
	}

	return resourceSpans, nil
}

func getTraceIDFromDatadogSpan(span *v1.Span) ([]byte, error) {
	firstHalfTraceID := ""
	for _, attr := range span.Attributes {
		if attr.Key == "_dd.p.tid" {
			firstHalfTraceID = attr.Value.GetStringValue()
		}
	}

	if len(firstHalfTraceID) == 0 {
		// not a datadog span
		return span.TraceId, nil
	}

	b64, err := strconv.ParseUint(firstHalfTraceID, 16, 64)
	if err != nil {
		return span.TraceId, fmt.Errorf("cannot parse UInt from span traceID: %w", err)
	}

	binary.BigEndian.PutUint64(span.TraceId[:8], b64)

	return span.TraceId[:], nil
}
