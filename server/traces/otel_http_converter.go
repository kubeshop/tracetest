package traces

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	v11 "go.opentelemetry.io/proto/otlp/common/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type HttpResourceSpans struct {
	v1.ResourceSpans
	ScopeSpans []*httpScopeSpans `json:"scopeSpans"`
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
	Status            *httpSpanStatus      `json:"status"`
}

type httpSpanStatus struct {
	Code string `json:"code"`
}

type httpSpanAttribute struct {
	v11.KeyValue
	Value map[string]interface{} `json:"value"`
}

func FromHttpOtelResourceSpans(resourceSpans []*HttpResourceSpans) model.Trace {
	flattenSpans := make([]*httpSpan, 0)
	for _, resource := range resourceSpans {
		for _, scopeSpans := range resource.ScopeSpans {
			flattenSpans = append(flattenSpans, scopeSpans.Spans...)
		}
	}

	traceID := ""
	spans := make([]model.Span, 0)
	for _, span := range flattenSpans {
		newSpan := convertHttpOtelSpanIntoSpan(span)
		traceID = hex.EncodeToString(span.TraceId)
		spans = append(spans, *newSpan)
	}

	return model.NewTrace(traceID, spans)
}

func convertHttpOtelSpanIntoSpan(span *httpSpan) *model.Span {
	attributes := make(model.Attributes, 0)
	for _, attribute := range span.Attributes {
		attributes[attribute.Key] = getHttpAttributeValue(attribute.Value)
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
	attributes[string(model.TracetestMetadataFieldParentID)] = createSpanID([]byte(span.ParentSpanId)).String()

	return &model.Span{
		ID:         spanID,
		Name:       span.Name,
		StartTime:  startTime,
		EndTime:    endTime,
		Parent:     nil,
		Children:   make([]*model.Span, 0),
		Attributes: attributes,
	}
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
