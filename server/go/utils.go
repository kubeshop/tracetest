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

				v := V1Span{
					TraceId:                hex.EncodeToString(sp.GetTraceId()),
					SpanId:                 hex.EncodeToString(sp.GetSpanId()),
					TraceState:             sp.GetTraceState(),
					ParentSpanId:           hex.EncodeToString(sp.GetParentSpanId()),
					Name:                   sp.GetName(),
					Kind:                   kind,
					StartTimeUnixNano:      strconv.FormatUint(sp.GetStartTimeUnixNano(), 10),
					EndTimeUnixNano:        strconv.FormatUint(sp.GetEndTimeUnixNano(), 10),
					Attributes:             mapAttributes(sp.GetAttributes()),
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
