package executor

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/kubeshop/tracetest/openapi"
	v11 "go.opentelemetry.io/proto/otlp/common/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func mapAnyValue(val *v11.AnyValue) openapi.V1AnyValue {

	var arrV openapi.V1ArrayValue
	if val.GetArrayValue() != nil {
		arrV = openapi.V1ArrayValue{
			Values: []openapi.V1AnyValue{},
		}
		for _, a := range val.GetArrayValue().GetValues() {
			a := mapAnyValue(a)
			arrV.Values = append(arrV.Values, a)
		}
	}
	var kvListVal openapi.V1KeyValueList
	if val.GetKvlistValue() != nil {
		kvListVal = openapi.V1KeyValueList{
			Values: []openapi.V1KeyValue{},
		}

		for _, kv := range val.GetKvlistValue().GetValues() {
			v := openapi.V1KeyValue{
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
	return openapi.V1AnyValue{
		StringValue: val.GetStringValue(),
		BoolValue:   val.GetBoolValue(),
		IntValue:    intVal,
		DoubleValue: val.GetDoubleValue(),
		ArrayValue:  arrV,
		KvlistValue: kvListVal,
		BytesValue:  string(val.GetBytesValue()),
	}
}

func mapAttributes(kvs []*v11.KeyValue) []openapi.V1KeyValue {
	var res []openapi.V1KeyValue

	for _, kv := range kvs {
		v := openapi.V1KeyValue{
			Key:   kv.GetKey(),
			Value: mapAnyValue(kv.GetValue()),
		}
		res = append(res, v)
	}
	return res
}

func mapTrace(tr *v1.TracesData) openapi.ApiV3SpansResponseChunk {
	res := openapi.ApiV3SpansResponseChunk{
		ResourceSpans: []openapi.V1ResourceSpans{},
	}

	for _, t := range tr.GetResourceSpans() {
		var ilsV []openapi.V1InstrumentationLibrarySpans

		for _, ils := range t.GetInstrumentationLibrarySpans() {
			var sps []openapi.V1Span
			for _, sp := range ils.GetSpans() {
				var kind openapi.SpanSpanKind
				switch sp.GetKind() {
				case v1.Span_SPAN_KIND_UNSPECIFIED:
					kind = openapi.UNSPECIFIED
				case v1.Span_SPAN_KIND_INTERNAL:
					kind = openapi.INTERNAL
				case v1.Span_SPAN_KIND_SERVER:
					kind = openapi.SERVER
				case v1.Span_SPAN_KIND_CLIENT:
					kind = openapi.CLIENT
				case v1.Span_SPAN_KIND_PRODUCER:
					kind = openapi.PRODUCER
				case v1.Span_SPAN_KIND_CONSUMER:
					kind = openapi.CONSUMER
				}

				var code openapi.StatusStatusCode
				switch sp.GetStatus().GetCode() {
				case v1.Status_STATUS_CODE_UNSET:
					code = openapi.UNSET
				case v1.Status_STATUS_CODE_OK:
					code = openapi.OK
				case v1.Status_STATUS_CODE_ERROR:
					code = openapi.ERROR
				}
				var events []openapi.SpanEvent
				for _, ev := range sp.GetEvents() {
					v := openapi.SpanEvent{
						TimeUnixNano:           strconv.FormatUint(ev.GetTimeUnixNano(), 10),
						Name:                   ev.GetName(),
						Attributes:             mapAttributes(ev.GetAttributes()),
						DroppedAttributesCount: int64(ev.GetDroppedAttributesCount()),
					}
					events = append(events, v)
				}

				var links []openapi.SpanLink
				for _, l := range sp.GetLinks() {
					v := openapi.SpanLink{
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
					spanDuration := (sp.GetEndTimeUnixNano() - sp.GetStartTimeUnixNano()) / 1000 / 1000 // in milliseconds

					attributes = append(attributes, openapi.V1KeyValue{
						Key: "tracetest.span.duration",
						Value: openapi.V1AnyValue{
							IntValue: strconv.FormatUint(spanDuration, 10),
						},
					})
				}

				attributes = append(attributes, openapi.V1KeyValue{
					Key: "tracetest.span.type",
					Value: openapi.V1AnyValue{
						StringValue: spanType(attributes),
					},
				})

				v := openapi.V1Span{
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
					Status: openapi.V1Status{
						Message: sp.GetStatus().GetMessage(),
						Code:    code,
					},
				}
				sps = append(sps, v)
			}
			v := openapi.V1InstrumentationLibrarySpans{
				InstrumentationLibrary: openapi.V1InstrumentationLibrary{
					Name:    ils.GetInstrumentationLibrary().GetName(),
					Version: ils.GetInstrumentationLibrary().GetVersion(),
				},
				Spans:     sps,
				SchemaUrl: ils.GetSchemaUrl(),
			}

			ilsV = append(ilsV, v)
		}
		sp := openapi.V1ResourceSpans{
			Resource: openapi.V1Resource{
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

func spanType(attrs []openapi.V1KeyValue) string {
	// based on https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/trace/semantic_conventions
	// using the first required attribute for each type
	for _, attr := range attrs {
		switch true {
		case strings.HasPrefix(attr.Key, "http."):
			return "http"
		case strings.HasPrefix(attr.Key, "db."):
			return "database"
		case strings.HasPrefix(attr.Key, "rpc."):
			return "rpc"
		case strings.HasPrefix(attr.Key, "messaging."):
			return "messaging"
		case strings.HasPrefix(attr.Key, "faas."):
			return "faas"
		case strings.HasPrefix(attr.Key, "exception."):
			return "exception"
		}
	}
	return "unknown"
}

func headersToKeyValueList(headers []openapi.HttpResponseHeaders) *v11.KeyValueList {
	mapped := []*v11.KeyValue{}

	for _, h := range headers {
		mapped = append(mapped, &v11.KeyValue{
			Key: h.Key,
			Value: &v11.AnyValue{
				Value: &v11.AnyValue_StringValue{StringValue: h.Value},
			},
		})
	}

	return &v11.KeyValueList{Values: mapped}
}

func FixParent(tr *v1.TracesData, response openapi.HttpResponse) (*v1.TracesData, error) {
	spans := make(map[string]*v1.Span)
	for _, rs := range tr.ResourceSpans {
		for _, ils := range rs.InstrumentationLibrarySpans {
			for _, sp := range ils.Spans {
				spans[string(sp.SpanId)] = sp
			}
		}
	}

	// Find parent span
	var parentSpan *v1.Span
	for _, sp := range spans {
		if sp.ParentSpanId == nil {
			continue
		}
		_, ok := spans[string(sp.ParentSpanId)]
		if !ok {
			parentSpan = sp
		}
	}
	if parentSpan == nil {
		return nil, errors.New("no parentspan")
	}

	tracetestAttrs := []*v11.KeyValue{
		{
			Key: "tracetest.response.status",
			Value: &v11.AnyValue{
				Value: &v11.AnyValue_StringValue{StringValue: strconv.FormatInt(int64(response.StatusCode), 10)},
			},
		},
		{
			Key: "tracetest.response.body",
			Value: &v11.AnyValue{
				Value: &v11.AnyValue_StringValue{StringValue: response.Body},
			},
		},
		{
			Key: "tracetest.response.headers",
			Value: &v11.AnyValue{
				Value: &v11.AnyValue_KvlistValue{KvlistValue: headersToKeyValueList(response.Headers)},
			},
		},
	}

	parentSpan.ParentSpanId = nil
	parentSpan.Attributes = append(parentSpan.Attributes, tracetestAttrs...)

	return tr, nil
}
