package selectors

import (
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/test"
)

func FromSpanQuery(sq test.SpanQuery) Selector {
	sel, _ := New(string(sq))
	return sel
}

type Selector struct {
	SpanSelectors []SpanSelector
}

func (s Selector) Filter(trace model.Trace) model.Spans {
	if len(s.SpanSelectors) == 0 {
		// empty selector should select everything
		return getAllSpans(trace)
	}

	allFilteredSpans := make([]model.Span, 0)
	for _, spanSelector := range s.SpanSelectors {
		spans := filterSpans(trace.RootSpan, spanSelector)
		allFilteredSpans = append(allFilteredSpans, spans...)
	}

	return allFilteredSpans
}

func getAllSpans(trace model.Trace) model.Spans {
	var allSpans = make(model.Spans, 0)
	traverseTree(trace.RootSpan, func(span model.Span) {
		allSpans = append(allSpans, span)
	})

	return allSpans
}

type SpanSelector struct {
	Filters       []filter
	PseudoClass   PseudoClass
	ChildSelector *SpanSelector
}

func (ss SpanSelector) MatchesFilters(span model.Span) bool {
	for _, filter := range ss.Filters {
		if err := filter.Filter(span); err != nil {
			return false
		}
	}

	return true
}

type FilterFunction struct {
	Filter func(model.Span, string, Value) error
	Name   string
}

type filter struct {
	Property  string
	Operation FilterFunction
	Value     Value
}

func (f filter) Filter(span model.Span) error {
	return f.Operation.Filter(span, f.Property, f.Value)
}

var (
	ValueNull    = "null"
	ValueInt     = "int"
	ValueString  = "string"
	ValueFloat   = "float"
	ValueBoolean = "boolean"
)

type Value struct {
	Type    string
	String  string
	Int     int64
	Float   float64
	Boolean bool
}

func (v Value) AsString() string {
	switch v.Type {
	case ValueInt:
		return fmt.Sprintf("%d", v.Int)
	case ValueBoolean:
		return fmt.Sprintf("%t", v.Boolean)
	case ValueFloat:
		return fmt.Sprintf("%.2f", v.Float)
	case ValueString:
		unquotedString, _ := strconv.Unquote(v.String)
		return unquotedString
	default:
		return "null"
	}
}
