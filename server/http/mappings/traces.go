package mappings

import (
	"strconv"

	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/pkg/timing"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

// out

func (m OpenAPI) Trace(in *traces.Trace) openapi.Trace {
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

func (m OpenAPI) Span(in traces.Span) openapi.Span {
	parentID := ""
	if in.Parent != nil {
		parentID = in.Parent.ID.String()
	}

	attributes := make(map[string]string, in.Attributes.Len())
	for name, value := range in.Attributes.Values() {
		attributes[name] = value
		if m.traceConversionConfig.IsTimeField(name) {
			valueAsInt, _ := strconv.Atoi(value)
			attributes[name] = traces.ConvertNanoSecondsIntoProperTimeUnit(valueAsInt)
		}
	}

	kind := string(in.Kind)
	if kind == "" {
		kind = string(traces.SpanKindUnespecified)
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

func (m OpenAPI) Spans(in []*traces.Span) []openapi.Span {
	spans := make([]openapi.Span, len(in))
	for i, s := range in {
		spans[i] = m.Span(*s)
	}

	return spans
}

// in

func (m Model) Trace(in openapi.Trace) *traces.Trace {
	tid, _ := trace.TraceIDFromHex(in.TraceId)

	flat := make(map[trace.SpanID]*traces.Span, len(in.Flat))
	for id, span := range in.Flat {
		sid, _ := trace.SpanIDFromHex(id)
		span := m.Span(span, nil)
		flat[sid] = &span
	}

	tree := m.Span(in.Tree, nil)
	if spanIsEmpty(tree) && len(flat) == 0 {
		return nil
	}

	return &traces.Trace{
		ID:       tid,
		RootSpan: tree,
		Flat:     flat,
	}
}

func spanIsEmpty(span traces.Span) bool {
	return span.ID == trace.SpanID{} && span.Attributes.Len() == 0 && span.Name == "" && len(span.Children) == 0
}

func (m Model) Span(in openapi.Span, parent *traces.Span) traces.Span {
	sid, _ := trace.SpanIDFromHex(in.Id)
	span := traces.Span{
		ID:         sid,
		Attributes: traces.NewAttributes(in.Attributes),
		Name:       in.Name,
		StartTime:  timing.ParseUnix(int64(in.StartTime)),
		EndTime:    timing.ParseUnix(int64(in.EndTime)),
		Parent:     parent,
	}
	span.Children = m.Spans(in.Children, &span)

	return span
}

func (m Model) Spans(in []openapi.Span, parent *traces.Span) []*traces.Span {
	spans := make([]*traces.Span, len(in))
	for i, s := range in {
		span := m.Span(s, parent)
		spans[i] = &span
	}

	return spans
}
