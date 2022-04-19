package openapi

import (
	"encoding/hex"
	"strconv"

	v11 "go.opentelemetry.io/proto/otlp/common/v1"
	res "go.opentelemetry.io/proto/otlp/resource/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func mapAnyValue(val *v11.AnyValue) V1AnyValue {

	var arrV V1ArrayValue
	if val.GetArrayValue() != nil {
		arrV = V1ArrayValue{
			Values: []V1AnyValue{},
		}
		for _, a := range val.GetArrayValue().GetValues() {
			a := mapAnyValue(a)
			arrV.Values = append(arrV.Values, a)
		}
	}
	var kvListVal V1KeyValueList
	if val.GetKvlistValue() != nil {
		kvListVal = V1KeyValueList{
			Values: []V1KeyValue{},
		}

		for _, kv := range val.GetKvlistValue().GetValues() {
			v := V1KeyValue{
				Key:   kv.GetKey(),
				Value: mapAnyValue(kv.GetValue()),
			}

			kvListVal.Values = append(kvListVal.Values, v)
		}
	}

	intVal := ""
	if i, ok := val.GetValue().(*v11.AnyValue_IntValue); ok {
		intVal = strconv.FormatInt(i.IntValue, 10)
	}
	return V1AnyValue{
		StringValue: val.GetStringValue(),
		BoolValue:   val.GetBoolValue(),
		IntValue:    intVal,
		DoubleValue: val.GetDoubleValue(),
		ArrayValue:  arrV,
		KvlistValue: kvListVal,
		BytesValue:  string(val.GetBytesValue()),
	}
}

func mapAttributes(kvs []*v11.KeyValue) []V1KeyValue {
	var res []V1KeyValue

	for _, kv := range kvs {
		v := V1KeyValue{
			Key:   kv.GetKey(),
			Value: mapAnyValue(kv.GetValue()),
		}
		res = append(res, v)
	}
	return res
}

func mapTrace(tr *v1.TracesData) ApiV3SpansResponseChunk {
	res := ApiV3SpansResponseChunk{
		ResourceSpans: []V1ResourceSpans{},
	}

	for _, t := range tr.GetResourceSpans() {
		var ilsV []V1InstrumentationLibrarySpans

		for _, ils := range t.GetInstrumentationLibrarySpans() {
			var sps []V1Span
			for _, sp := range ils.GetSpans() {
				var kind SpanSpanKind
				switch sp.GetKind() {
				case v1.Span_SPAN_KIND_UNSPECIFIED:
					kind = UNSPECIFIED
				case v1.Span_SPAN_KIND_INTERNAL:
					kind = INTERNAL
				case v1.Span_SPAN_KIND_SERVER:
					kind = SERVER
				case v1.Span_SPAN_KIND_CLIENT:
					kind = CLIENT
				case v1.Span_SPAN_KIND_PRODUCER:
					kind = PRODUCER
				case v1.Span_SPAN_KIND_CONSUMER:
					kind = CONSUMER
				}

				var code StatusStatusCode
				switch sp.GetStatus().GetCode() {
				case v1.Status_STATUS_CODE_UNSET:
					code = UNSET
				case v1.Status_STATUS_CODE_OK:
					code = OK
				case v1.Status_STATUS_CODE_ERROR:
					code = ERROR
				}
				var events []SpanEvent
				for _, ev := range sp.GetEvents() {
					v := SpanEvent{
						TimeUnixNano:           strconv.FormatUint(ev.GetTimeUnixNano(), 10),
						Name:                   ev.GetName(),
						Attributes:             mapAttributes(ev.GetAttributes()),
						DroppedAttributesCount: int64(ev.GetDroppedAttributesCount()),
					}
					events = append(events, v)
				}

				var links []SpanLink
				for _, l := range sp.GetLinks() {
					v := SpanLink{
						TraceId:                hex.EncodeToString(l.GetTraceId()),
						SpanId:                 hex.EncodeToString(l.GetSpanId()),
						TraceState:             l.GetTraceState(),
						Attributes:             mapAttributes(l.GetAttributes()),
						DroppedAttributesCount: int64(l.GetDroppedAttributesCount()),
					}

					links = append(links, v)
				}

				var startTime string
				if sp.GetStartTimeUnixNano() != 0 {
					startTime = strconv.FormatUint(sp.GetStartTimeUnixNano(), 10)
				}
				var endTime string
				if sp.GetEndTimeUnixNano() != 0 {
					endTime = strconv.FormatUint(sp.GetEndTimeUnixNano(), 10)
				}

				attributes := mapAttributes(sp.GetAttributes())

				if sp.GetStartTimeUnixNano() != 0 && sp.GetEndTimeUnixNano() != 0 {
					spanDuration := (sp.GetEndTimeUnixNano() - sp.GetStartTimeUnixNano()) // in nanoseconds

					attributes = append(attributes, V1KeyValue{
						Key: "tracetest.span.duration",
						Value: V1AnyValue{
							IntValue: strconv.FormatUint(spanDuration, 10),
						},
					})
				}

				attributes = append(attributes, V1KeyValue{
					Key: "tracetest.span.type",
					Value: V1AnyValue{
						StringValue: spanType(attributes),
					},
				})

				v := V1Span{
					TraceId:                hex.EncodeToString(sp.GetTraceId()),
					SpanId:                 hex.EncodeToString(sp.GetSpanId()),
					TraceState:             sp.GetTraceState(),
					ParentSpanId:           hex.EncodeToString(sp.GetParentSpanId()),
					Name:                   sp.GetName(),
					Kind:                   kind,
					StartTimeUnixNano:      startTime,
					EndTimeUnixNano:        endTime,
					Attributes:             attributes,
					DroppedAttributesCount: int64(sp.GetDroppedAttributesCount()),
					Events:                 events,
					DroppedEventsCount:     int64(sp.GetDroppedEventsCount()),
					Links:                  links,
					DroppedLinksCount:      int64(sp.GetDroppedLinksCount()),
					Status: V1Status{
						Message: sp.GetStatus().GetMessage(),
						Code:    code,
					},
				}
				sps = append(sps, v)
			}
			v := V1InstrumentationLibrarySpans{
				InstrumentationLibrary: V1InstrumentationLibrary{
					Name:    ils.GetInstrumentationLibrary().GetName(),
					Version: ils.GetInstrumentationLibrary().GetVersion(),
				},
				Spans:     sps,
				SchemaUrl: ils.GetSchemaUrl(),
			}

			ilsV = append(ilsV, v)
		}
		sp := V1ResourceSpans{
			Resource: V1Resource{
				Attributes:             mapAttributes(t.GetResource().GetAttributes()),
				DroppedAttributesCount: int64(t.GetResource().GetDroppedAttributesCount()),
			},
			InstrumentationLibrarySpans: ilsV,
			SchemaUrl:                   t.GetSchemaUrl(),
		}
		res.ResourceSpans = append(res.ResourceSpans, sp)
	}
	return res
}

func spanType(attrs []V1KeyValue) string {
	// based on https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/trace/semantic_conventions
	// using the first required attribute for each type
	for _, attr := range attrs {
		switch true {
		case attr.Key == "http.method":
			return "http"
		case attr.Key == "db.system":
			return "database"
		case attr.Key == "rpc.system":
			return "rpc"
		case attr.Key == "messaging.system":
			return "messaging"
		case attr.Key == "faas.trigger", attr.Key == "faas.execution":
			// faas has no required attr, so anyone works
			return "faas"
		case attr.Key == "exception.type", attr.Key == "exception.message":
			// at least one of the two must be present
			return "exception"
		}
	}
	return "unknown"
}

func headersToKeyValueList(headers []HttpResponseHeaders) *v11.KeyValueList {
	mapped := []*v11.KeyValue{}

	for _, h := range headers {
		mapped = append(mapped, &v11.KeyValue{
			Key: h.Key,
			Value: &v11.AnyValue{
				Value: &v11.AnyValue_StringValue{h.Value},
			},
		})
	}

	return &v11.KeyValueList{Values: mapped}
}

func FixParent(tr *v1.TracesData, traceID, parentSpanID string, response HttpResponse) *v1.TracesData {
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

	attributes := []*v11.KeyValue{
		{
			Key: "tracetest.response.status",
			Value: &v11.AnyValue{
				Value: &v11.AnyValue_StringValue{strconv.FormatInt(int64(response.StatusCode), 10)},
			},
		},
		{
			Key: "tracetest.response.body",
			Value: &v11.AnyValue{
				Value: &v11.AnyValue_StringValue{response.Body},
			},
		},
		{
			Key: "tracetest.response.headers",
			Value: &v11.AnyValue{
				Value: &v11.AnyValue_KvlistValue{headersToKeyValueList(response.Headers)},
			},
		},
	}

	// this is the parent span
	rs := &v1.ResourceSpans{
		Resource: &res.Resource{
			Attributes:             attributes,
			DroppedAttributesCount: 0,
		},
		InstrumentationLibrarySpans: []*v1.InstrumentationLibrarySpans{
			{
				InstrumentationLibrary: &v11.InstrumentationLibrary{
					Name:    "tracetest",
					Version: "v1",
				},
				Spans: []*v1.Span{
					{
						TraceId:      []byte(traceID),
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
