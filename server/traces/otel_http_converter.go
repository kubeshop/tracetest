package traces

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"time"

	v11 "go.opentelemetry.io/proto/otlp/common/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type HttpResourceSpans struct {
	v1.ResourceSpans
	ScopeSpans                  []*httpScopeSpans `json:"scopeSpans"`
	InstrumentationLibrarySpans []*httpScopeSpans `json:"instrumentationLibrarySpans"`
}

type httpScopeSpans struct {
	v1.ScopeSpans
	Spans []*httpSpan `json:"spans"`
}

type httpSpan struct {
	v1.Span
	SpanId            string               `json:"spanId"`
	ParentSpanId      string               `json:"parentSpanId"`
	Kind              string               `json:"kind"`
	StartTimeUnixNano string               `json:"startTimeUnixNano"`
	EndTimeUnixNano   string               `json:"endTimeUnixNano"`
	Attributes        []*httpSpanAttribute `json:"attributes"`
	Events            []*httpSpanEvent     `json:"events"`
	Status            *httpSpanStatus      `json:"status"`
}

type httpSpanStatus struct {
	Code string `json:"code"`
}

type httpSpanEvent struct {
	Name       string               `json:"name"`
	Timestamp  string               `json:"timeUnixNano"`
	Attributes []*httpSpanAttribute `json:"attributes"`
}

type httpSpanAttribute struct {
	v11.KeyValue
	Value map[string]interface{} `json:"value"`
}

func FromHttpOtelResourceSpans(resourceSpans []*HttpResourceSpans) Trace {
	flattenSpans := make([]*httpSpan, 0)
	for _, resource := range resourceSpans {
		for _, scopeSpans := range resource.ScopeSpans {
			flattenSpans = append(flattenSpans, scopeSpans.Spans...)
		}

		for _, librarySpans := range resource.InstrumentationLibrarySpans {
			flattenSpans = append(flattenSpans, librarySpans.Spans...)
		}
	}

	traceID := ""
	spans := make([]Span, 0)
	for _, span := range flattenSpans {
		newSpan := convertHttpOtelSpanIntoSpan(span)
		traceID = hex.EncodeToString(span.TraceId)
		spans = append(spans, *newSpan)
	}

	return NewTrace(traceID, spans)
}

func convertHttpOtelSpanIntoSpan(span *httpSpan) *Span {
	attributes := NewAttributes()
	for _, attribute := range span.Attributes {
		attributes.Set(attribute.Key, getHttpAttributeValue(attribute.Value))
	}

	var startTime, endTime time.Time

	if span.StartTimeUnixNano != "" {
		startTimeNs, _ := strconv.ParseInt(span.StartTimeUnixNano, 10, 64)
		startTime = time.Unix(0, startTimeNs)
	}

	if span.EndTimeUnixNano != "" {
		endTimeNs, _ := strconv.ParseInt(span.EndTimeUnixNano, 10, 64)
		endTime = time.Unix(0, endTimeNs)
	}

	spanID := createSpanID([]byte(span.SpanId))
	attributes.Set(TracetestMetadataFieldParentID, createSpanID([]byte(span.ParentSpanId)).String())

	return &Span{
		ID:         spanID,
		Name:       span.Name,
		StartTime:  startTime,
		EndTime:    endTime,
		Parent:     nil,
		Children:   make([]*Span, 0),
		Attributes: attributes,
		Events:     extractEventsFromHttpSpan(span),
	}
}

func extractEventsFromHttpSpan(span *httpSpan) []SpanEvent {
	output := make([]SpanEvent, 0, len(span.Events))
	for _, event := range span.Events {
		attributes := NewAttributes()
		for _, attribute := range event.Attributes {
			attributes.Set(attribute.Key, getHttpAttributeValue(attribute.Value))
		}

		var timestamp time.Time
		if event.Timestamp != "" {
			timestampNs, _ := strconv.ParseInt(span.StartTimeUnixNano, 10, 64)
			timestamp = time.Unix(0, timestampNs)
		}

		output = append(output, SpanEvent{
			Name:       event.Name,
			Timestamp:  timestamp,
			Attributes: attributes,
		})
	}

	return output
}

func getHttpAttributeValue(value map[string]interface{}) string {
	if value["stringValue"] != nil {
		return value["stringValue"].(string)
	}

	if value["intValue"] != nil {
		return fmt.Sprintf("%s", value["intValue"])
	}

	if value["doubleValue"] != nil {
		val, ok := value["doubleValue"].(float64)
		if ok {
			isFloatingPoint := math.Abs(val-math.Abs(val)) > 0.0
			if isFloatingPoint {
				return fmt.Sprintf("%f", val)
			}

			return fmt.Sprintf("%.0f", val)
		}
	}

	if value["boolValue"] != nil {
		return fmt.Sprintf("%t", value["boolValue"])
	}

	return "unsupported value type"
}
