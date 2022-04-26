package openapi_test

import (
	"strconv"
	"testing"

	openapi "github.com/kubeshop/tracetest/server/go"
	"github.com/stretchr/testify/assert"
	v11 "go.opentelemetry.io/proto/otlp/common/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

var wantAttributes = []*v11.KeyValue{
	{
		Key: "tracetest.response.status",
		Value: &v11.AnyValue{
			Value: &v11.AnyValue_StringValue{strconv.FormatInt(int64(200), 10)},
		},
	},
	{
		Key: "tracetest.response.body",
		Value: &v11.AnyValue{
			Value: &v11.AnyValue_StringValue{"body"},
		},
	},
}

func TestFixParent(t *testing.T) {
	td := &v1.TracesData{
		ResourceSpans: []*v1.ResourceSpans{
			{
				InstrumentationLibrarySpans: []*v1.InstrumentationLibrarySpans{
					{
						Spans: []*v1.Span{
							{
								TraceId:      []byte("traceid"),
								ParentSpanId: []byte("parentSpanID"),
								Attributes: []*v11.KeyValue{
									{
										Key: "test", Value: &v11.AnyValue{
											Value: &v11.AnyValue_StringValue{StringValue: "string"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	resp := openapi.HttpResponse{StatusCode: 200, Body: "body"}
	out, err := openapi.FixParent(td, resp)

	assert.NoError(t, err)
	assert.Equal(t, []uint8([]byte(nil)), out.ResourceSpans[0].InstrumentationLibrarySpans[0].Spans[0].ParentSpanId)
	//TODO check if returned attributes contain response body and status
}
