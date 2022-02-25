package openapi

import (
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	v11 "go.opentelemetry.io/proto/otlp/common/v1"
	res "go.opentelemetry.io/proto/otlp/resource/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func EncodeJSONPBResponse(i interface{}, status *int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status != nil {
		w.WriteHeader(*status)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	trace := i.(proto.Message)
	m := jsonpb.Marshaler{}
	return m.Marshal(w, trace)
}

func FixParent(tr *v1.TracesData, traceID, parentSpanID string) *v1.TracesData {
	spans := make(map[string]*v1.Span)
	for _, rs := range tr.ResourceSpans {
		for _, ils := range rs.InstrumentationLibrarySpans {
			for _, sp := range ils.Spans {
				spans[string(sp.SpanId)] = sp
			}
		}
	}

	// Fix parent id
	for _, sp := range spans {
		if sp.ParentSpanId == nil {
			continue
		}
		_, ok := spans[string(sp.ParentSpanId)]
		if !ok {
			sp.ParentSpanId = []byte(parentSpanID)
		}
	}
	rs := &v1.ResourceSpans{
		Resource: &res.Resource{
			Attributes:             []*v11.KeyValue{},
			DroppedAttributesCount: 0,
		},
		InstrumentationLibrarySpans: []*v1.InstrumentationLibrarySpans{
			{
				InstrumentationLibrary: &v11.InstrumentationLibrary{
					Name:    "tracetest",
					Version: "v1",
				},
				Spans: []*v1.Span{
					{TraceId: []byte(traceID),
						SpanId:       []byte(parentSpanID),
						ParentSpanId: nil,
						Name:         "tracetest",
						Kind:         v1.Span_SPAN_KIND_CLIENT,
					},
				},
			},
		},
		SchemaUrl: "",
	}
	tr.ResourceSpans = append(tr.ResourceSpans, rs)

	return tr
}
