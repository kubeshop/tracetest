package mappings

import (
	"strconv"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

// out

func (m OpenAPI) Trace(in *model.Trace) openapi.Trace {
	if in == nil {
		return openapi.Trace{}
	}

	flat := map[string]openapi.Span{}
	for id, span := range in.Flat {
		flat[id.String()] = m.Span(*span)
	}

	return openapi.Trace{
		TraceId: in.ID.String(),
		Tree:    m.Span(in.RootSpan),
		Flat:    flat,
	}
}

func (m OpenAPI) Span(in model.Span) openapi.Span {
	parentID := ""
	if in.Parent != nil {
		parentID = in.Parent.ID.String()
	}

	attributes := make(map[string]string, len(in.Attributes))
	for name, value := range in.Attributes {
		attributes[name] = value
		if m.traceConversionConfig.IsTimeField(name) {
			valueAsInt, _ := strconv.Atoi(value)
			attributes[name] = traces.ConvertNanoSecondsIntoProperTimeUnit(valueAsInt)
		}
	}

	kind := string(in.Kind)
	if kind == "" {
		kind = string(model.SpanKindUnespecified)
	}

	return openapi.Span{
		Id:         in.ID.String(),
		ParentId:   parentID,
		Kind:       kind,
		StartTime:  in.StartTime.UnixMilli(),
		EndTime:    in.EndTime.UnixMilli(),
		Attributes: attributes,
		Children:   m.Spans(in.Children),
		Name:       in.Name,
	}
}

func (m OpenAPI) Spans(in []*model.Span) []openapi.Span {
	spans := make([]openapi.Span, len(in))
	for i, s := range in {
		spans[i] = m.Span(*s)
	}

	return spans
}

// in

func (m Model) Trace(in openapi.Trace) *model.Trace {
	tid, _ := trace.TraceIDFromHex(in.TraceId)
	return &model.Trace{
		ID:       tid,
		RootSpan: m.Span(in.Tree, nil),
	}
}

func (m Model) Span(in openapi.Span, parent *model.Span) model.Span {
	sid, _ := trace.SpanIDFromHex(in.Id)
	span := model.Span{
		ID:         sid,
		Attributes: in.Attributes,
		Name:       in.Name,
		StartTime:  time.UnixMilli(int64(in.StartTime)),
		EndTime:    time.UnixMilli(int64(in.EndTime)),
		Parent:     parent,
	}
	span.Children = m.Spans(in.Children, &span)

	return span
}

func (m Model) Spans(in []openapi.Span, parent *model.Span) []*model.Span {
	spans := make([]*model.Span, len(in))
	for i, s := range in {
		span := m.Span(s, parent)
		spans[i] = &span
	}

	return spans
}
