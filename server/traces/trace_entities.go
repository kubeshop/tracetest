package traces

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/agent/workers/trigger"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/timing"
	"go.opentelemetry.io/otel/trace"
)

type Trace struct {
	ID       trace.TraceID
	RootSpan Span
	Flat     map[trace.SpanID]*Span `json:"-"`
}

func nonNilTraces(traces []*Trace) []*Trace {
	nonNil := make([]*Trace, 0)
	for _, trace := range traces {
		if trace == nil {
			continue
		}

		nonNil = append(nonNil, trace)
	}
	return nonNil
}

func MergeTraces(traces ...*Trace) *Trace {
	traces = nonNilTraces(traces)
	if len(traces) == 0 {
		return nil
	}

	traceID := traces[0].ID
	spans := make([]Span, 0)
	for _, trace := range traces {
		if trace == nil {
			continue
		}

		spans = append(spans, trace.Spans()...)
	}

	trace := NewTrace(traceID.String(), spans)

	return &trace
}

func NewTrace(traceID string, spans []Span) Trace {
	// remove all temporary root spans
	sanitizedSpans := make([]Span, 0)
	for _, span := range spans {
		if span.Name != TemporaryRootSpanName {
			sanitizedSpans = append(sanitizedSpans, span)
		}
	}
	spans = sanitizedSpans

	spanMap := make(map[string]*Span, 0)
	for _, span := range spans {
		spanCopy := span.setMetadataAttributes()
		spanID := span.ID.String()
		spanMap[spanID] = &spanCopy
	}

	rootSpans := make([]*Span, 0)
	for _, span := range spanMap {
		span.injectEventsIntoAttributes()

		parentID := span.Attributes.Get(TracetestMetadataFieldParentID)
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
		root = &Span{ID: id.NewRandGenerator().SpanID(), StartTime: time.Now(), EndTime: time.Now(), Name: TemporaryRootSpanName, Attributes: NewAttributes(), Children: []*Span{}}
	}

	for _, span := range allRoots {
		if root != span {
			root.Children = append(root.Children, span)
			span.Parent = root
		}
	}

	return root
}

// TODO: this is temp while we decide what to do with browser spans and how to handle them
func isBrowserSpan(attrs Attributes) bool {
	return attrs.Get("event_type") != "" || attrs.Get(TracetestMetadataFieldName) == "documentLoad"
}

func spanType(attrs Attributes) string {
	if isBrowserSpan(attrs) {
		return "general"
	}

	// based on https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/trace/semantic_conventions
	// using the first required attribute for each type
	for key := range attrs.Values() {
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
	timeDifference := timing.TimeDiff(span.StartTime, span.EndTime)
	return strconv.FormatInt(timeDifference.Nanoseconds(), 10)
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
	if !t.RootSpan.IsZero() || t.RootSpan.Name == TemporaryRootSpanName {
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

	if oldRoot.Attributes.values == nil {
		oldRoot.Attributes = NewAttributes()
	}
	oldRoot.Parent = &newRoot
	oldRoot.Attributes.Set(TracetestMetadataFieldParentID, newRoot.ID.String())

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
	resetCache()
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

func AugmentRootSpan(span Span, result trigger.TriggerResult) Span {
	return span.
		setMetadataAttributes().
		setTriggerResultAttributes(result)
}
