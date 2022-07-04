package traces

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Trace struct {
	ID       trace.TraceID
	RootSpan Span
	Flat     map[trace.SpanID]*Span `json:"-"`
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

	// We don't need the parent in the flat structure
	rootPtr.Parent = nil
	res[root.ID] = rootPtr
	for _, child := range root.Children {
		flattenSpans(res, *child)
	}

	// Remove children and parent because they are now part of the flatten structure
	rootPtr.Children = nil
}

type Attributes map[string]string

func (a Attributes) Get(key string) string {
	if v, ok := a[key]; ok {
		return v
	}

	return ""
}

type Span struct {
	ID         trace.SpanID
	Name       string
	StartTime  time.Time
	EndTime    time.Time
	Attributes Attributes

	Parent   *Span   `json:"-"`
	Children []*Span `json:"-"`
}

type encodedSpan struct {
	ID         string
	Name       string
	StartTime  string
	EndTime    string
	Attributes Attributes
	Children   []encodedSpan
}

func (s Span) MarshalJSON() ([]byte, error) {
	enc := encodeSpan(s)
	return json.Marshal(&enc)
}

func encodeSpan(s Span) encodedSpan {
	return encodedSpan{
		ID:         s.ID.String(),
		Name:       s.Name,
		StartTime:  fmt.Sprintf("%d", s.StartTime.UnixMilli()),
		EndTime:    fmt.Sprintf("%d", s.EndTime.UnixMilli()),
		Attributes: s.Attributes,
		Children:   encodeChildren(s.Children),
	}
}

func encodeChildren(children []*Span) []encodedSpan {
	res := make([]encodedSpan, len(children))
	for i, c := range children {
		res[i] = encodeSpan(*c)
	}
	return res
}

func (s *Span) UnmarshalJSON(data []byte) error {
	aux := encodedSpan{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("unmarshal span: %w", err)
	}

	return s.decodeSpan(aux)
}

func (s *Span) decodeSpan(aux encodedSpan) error {
	sid, err := trace.SpanIDFromHex(aux.ID)
	if err != nil {
		return fmt.Errorf("unmarshal span: %w", err)
	}

	children, err := decodeChildren(s, aux.Children)
	if err != nil {
		return fmt.Errorf("unmarshal span: %w", err)
	}

	startTime, err := getTimeFromString(aux.StartTime)
	if err != nil {
		return fmt.Errorf("unmarshal span: %w", err)
	}

	endTime, err := getTimeFromString(aux.EndTime)
	if err != nil {
		return fmt.Errorf("unmarshal span: %w", err)
	}

	s.ID = sid
	s.Name = aux.Name
	s.StartTime = startTime.UTC()
	s.EndTime = endTime.UTC()
	s.Attributes = aux.Attributes
	s.Children = children

	return nil
}

func getTimeFromString(value string) (time.Time, error) {
	milliseconds, err := strconv.Atoi(value)
	if err != nil {
		// Maybe it is in RFC3339 format. Convert it for compatibility sake
		output, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return time.Time{}, fmt.Errorf("could not convert string (%s) to time: %w", value, err)
		}

		return output, nil
	}

	return time.UnixMilli(int64(milliseconds)), nil
}

func decodeChildren(parent *Span, children []encodedSpan) ([]*Span, error) {
	if len(children) == 0 {
		return nil, nil
	}
	res := make([]*Span, len(children))
	for i, c := range children {
		span := &Span{
			Parent: parent,
		}
		if err := span.decodeSpan(c); err != nil {
			return nil, fmt.Errorf("unmarshal children: %w", err)
		}

		children, err := decodeChildren(span, c.Children)
		if err != nil {
			return nil, fmt.Errorf("unmarshal children: %w", err)
		}

		span.Children = children
		res[i] = span
	}
	return res, nil
}
