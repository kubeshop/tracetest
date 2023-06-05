package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

type Trace struct {
	ID       trace.TraceID
	RootSpan Span
	Flat     map[trace.SpanID]*Span `json:"-"`
}

func NewTrace(traceID string, spans []Span) Trace {
	spanMap := make(map[string]*Span, 0)
	for _, span := range spans {
		spanCopy := span.setMetadataAttributes()
		spanID := span.ID.String()
		spanMap[spanID] = &spanCopy
	}

	rootSpans := make([]*Span, 0)
	for _, span := range spanMap {
		parentID := span.Attributes["parent_id"]
		parentSpan, found := spanMap[parentID]
		if !found {
			rootSpans = append(rootSpans, span)
			continue
		}

		parentSpan.Children = append(parentSpan.Children, span)
		span.Parent = parentSpan
	}

	rootSpan := getRootSpan(rootSpans)

	id, _ := trace.TraceIDFromHex(traceID)
	trace := Trace{
		ID:       id,
		RootSpan: *rootSpan,
	}

	trace = trace.Sort()

	return trace
}

func getRootSpan(allRoots []*Span) *Span {
	if len(allRoots) == 1 {
		return allRoots[0]
	}

	var root *Span
	for _, span := range allRoots {
		if span.Name == TriggerSpanName || span.Name == TemporaryRootSpanName {
			// This span should be promoted because it's either a temporary root or the definitive root
			if root != nil && root.Name == TriggerSpanName {
				// We cannot override the root because we already have the definitive root, otherwise,
				// we will replace the definitive root with the temporary root.
				continue
			}

			root = span
		}
	}

	if root == nil {
		root = &Span{ID: IDGen.SpanID(), Name: TemporaryRootSpanName, Attributes: make(Attributes), Children: []*Span{}}
	}

	for _, span := range allRoots {
		if root != span {
			root.Children = append(root.Children, span)
			span.Parent = root
		}
	}

	return root
}

func spanType(attrs Attributes) string {
	// based on https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/trace/semantic_conventions
	// using the first required attribute for each type
	for key := range attrs {
		switch true {
		case strings.HasPrefix(key, "http."):
			return "http"
		case strings.HasPrefix(key, "db."):
			return "database"
		case strings.HasPrefix(key, "rpc."):
			return "rpc"
		case strings.HasPrefix(key, "messaging."):
			return "messaging"
		case strings.HasPrefix(key, "faas."):
			return "faas"
		case strings.HasPrefix(key, "exception."):
			return "exception"
		}
	}
	return "general"
}

func spanDuration(span Span) string {
	timeDifference := span.EndTime.Sub(span.StartTime)
	return fmt.Sprintf("%d", int64(timeDifference))
}

func (t *Trace) Sort() Trace {
	sortedRoot := sortSpanChildren(t.RootSpan)

	trace := Trace{
		ID:       t.ID,
		RootSpan: sortedRoot,
		Flat:     make(map[trace.SpanID]*Span, 0),
	}

	flattenSpans(trace.Flat, sortedRoot)

	return trace
}

const TriggerSpanName = "Tracetest trigger"
const TemporaryRootSpanName = "Temporary Tracetest root span"

func (t *Trace) HasRootSpan() bool {
	return t.RootSpan.Name == TriggerSpanName
}

func (t *Trace) InsertRootSpan(newRoot Span) *Trace {
	if !t.RootSpan.IsZero() {
		newRoot = replaceRoot(t.RootSpan, newRoot)
	}

	trace := Trace{
		ID:       t.ID,
		RootSpan: newRoot,
	}

	sorted := trace.Sort()

	return &sorted
}

func replaceRoot(oldRoot, newRoot Span) Span {
	if oldRoot.Name == TemporaryRootSpanName {
		// Replace the temporary root with the actual root
		newRoot.Children = oldRoot.Children
		for _, span := range oldRoot.Children {
			span.Parent = &newRoot
		}
		return newRoot
	}

	if oldRoot.Attributes == nil {
		oldRoot.Attributes = make(Attributes)
	}
	oldRoot.Parent = &newRoot
	oldRoot.Attributes["parent_id"] = newRoot.ID.String()

	newRoot.Children = append(newRoot.Children, &oldRoot)

	return newRoot
}

func (t *Trace) Spans() []Span {
	if t == nil {
		return []Span{}
	}

	spans := make([]Span, 0, len(t.Flat))
	for _, span := range t.Flat {
		spans = append(spans, *span)
	}

	return spans
}

func sortSpanChildren(span Span) Span {
	sort.SliceStable(span.Children, func(i, j int) bool {
		return span.Children[i].StartTime.Before(span.Children[j].StartTime)
	})

	children := make([]*Span, 0, len(span.Children))
	for _, childSpan := range span.Children {
		newChild := sortSpanChildren(*childSpan)
		children = append(children, &newChild)
	}

	span.Children = children

	return span
}

func (t *Trace) UnmarshalJSON(data []byte) error {
	type Alias Trace
	aux := &struct {
		ID string
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("unmarshal trace: %w", err)
	}
	tid, err := trace.TraceIDFromHex(aux.ID)
	if err != nil {
		return fmt.Errorf("unmarshal trace: %w", err)
	}

	t.ID = tid
	t.Flat = map[trace.SpanID]*Span{}
	flattenSpans(t.Flat, t.RootSpan)
	return nil
}

func (t Trace) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID       string
		RootSpan Span
	}{
		ID:       t.ID.String(),
		RootSpan: t.RootSpan,
	})
}

func flattenSpans(res map[trace.SpanID]*Span, root Span) {
	rootPtr := &root

	res[root.ID] = rootPtr
	for _, child := range root.Children {
		flattenSpans(res, *child)
	}

	// Remove children and parent because they are now part of the flatten structure
	rootPtr.Children = nil
}
