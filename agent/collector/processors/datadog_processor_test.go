package processors_test

import (
	"testing"

	"github.com/kubeshop/tracetest/agent/collector/processors"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/stretchr/testify/assert"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v11 "go.opentelemetry.io/proto/otlp/common/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func TestTraceIDConcat(t *testing.T) {
	traceID := id.NewRandGenerator().TraceID()

	firstHalf := traceID.String()[0:16]
	lastHalf := traceID[8:]

	request := &pb.ExportTraceServiceRequest{
		ResourceSpans: []*v1.ResourceSpans{
			{
				ScopeSpans: []*v1.ScopeSpans{
					{
						Spans: []*v1.Span{
							{
								TraceId: lastHalf,
								Attributes: []*v11.KeyValue{
									{Key: "_dd.p.tid", Value: &v11.AnyValue{
										Value: &v11.AnyValue_StringValue{
											StringValue: string(firstHalf),
										},
									}},
								},
							},
						},
					},
				},
			},
		},
	}

	originalSpanIds := getSpansIDs(request)

	datadogProcessor := processors.NewDatadogProcessor()
	alteredRequest := datadogProcessor.Process(request)

	newSpanIds := getSpansIDs(alteredRequest)

	for i, originalSpanID := range originalSpanIds {
		newSpanID := newSpanIds[i]
		assert.NotEqual(t, originalSpanID, newSpanID)
		assert.Equal(t, traceID[:], newSpanID)
	}
}

func getSpansIDs(request *pb.ExportTraceServiceRequest) [][]byte {
	ids := make([][]byte, 0)
	for _, resourceSpan := range request.ResourceSpans {
		for _, scopeSpan := range resourceSpan.ScopeSpans {
			for _, span := range scopeSpan.Spans {
				ids = append(ids, span.TraceId)
			}
		}
	}

	return ids
}
