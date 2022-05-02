package traces

import (
	"go.opentelemetry.io/otel/trace"
)

type Trace struct {
	ID       trace.TraceID
	RootSpan Span
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

	Parent   *Span
	Children []*Span
}
