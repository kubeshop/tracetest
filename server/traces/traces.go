package traces

import (
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

type Trace struct {
	ID       trace.TraceID
	RootSpan Span
	Flat     map[trace.SpanID]*Span `json:"-"`
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
	flattenSpans(t.Flat, &t.RootSpan)
	return nil
}

func (t *Trace) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID       string
		RootSpan Span
	}{
		ID:       t.ID.String(),
		RootSpan: t.RootSpan,
	})
}

func flattenSpans(res map[trace.SpanID]*Span, root *Span) {
	res[root.ID] = root
	for _, child := range root.Children {
		res[child.ID] = child
		if len(child.Children) > 0 {
			flattenSpans(res, child)
		}
	}
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
	Attributes Attributes

	Parent   *Span   `json:"-"`
	Children []*Span `json:"-"`
}

type encodedSpan struct {
	ID         string
	Name       string
	Attributes Attributes
	Children   []encodedSpan
}

func (s Span) MarshalJSON() ([]byte, error) {
	return json.Marshal(&encodedSpan{
		ID:         s.ID.String(),
		Name:       s.Name,
		Attributes: s.Attributes,
		Children:   encodeChildren(s.Children),
	})
}

func encodeChildren(children []*Span) []encodedSpan {
	res := make([]encodedSpan, len(children))
	for i, c := range children {
		res[i] = encodedSpan{
			ID:         c.ID.String(),
			Name:       c.Name,
			Attributes: c.Attributes,
			Children:   encodeChildren(c.Children),
		}
	}
	return res
}

func (s *Span) UnmarshalJSON(data []byte) error {
	aux := encodedSpan{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("unmarshal span: %w", err)
	}

	sid, err := trace.SpanIDFromHex(aux.ID)
	if err != nil {
		return fmt.Errorf("unmarshal span: %w", err)
	}

	children, err := decodeChildren(s, aux.Children)
	if err != nil {
		return fmt.Errorf("unmarshal span: %w", err)
	}

	s.ID = sid
	s.Name = aux.Name
	s.Attributes = aux.Attributes
	s.Children = children
	return nil
}

func decodeChildren(parent *Span, children []encodedSpan) ([]*Span, error) {
	if len(children) == 0 {
		return nil, nil
	}
	res := make([]*Span, len(children))
	for i, c := range children {
		sid, err := trace.SpanIDFromHex(c.ID)
		if err != nil {
			return nil, fmt.Errorf("span decode children: %w", err)
		}

		span := &Span{
			ID:         sid,
			Name:       c.Name,
			Attributes: c.Attributes,
			Parent:     parent,
		}

		children, err := decodeChildren(span, c.Children)
		if err != nil {
			return nil, err
		}

		span.Children = children
		res[i] = span
	}
	return res, nil
}
