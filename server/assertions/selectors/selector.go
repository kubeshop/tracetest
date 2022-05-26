package selectors

import (
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
)

func FromSpanQuery(sq model.SpanQuery) Selector {
	sel, _ := New(string(sq))
	return sel
}

type Selector struct {
	spanSelectors []spanSelector
}

func (s Selector) Filter(trace traces.Trace) []traces.Span {
	if len(s.spanSelectors) == 0 {
		// empty selector should select everything
		return getAllSpans(trace)
	}

	allFilteredSpans := make([]traces.Span, 0)
	for _, spanSelector := range s.spanSelectors {
		spans := filterSpans(trace.RootSpan, spanSelector)
		allFilteredSpans = append(allFilteredSpans, spans...)
	}

	return allFilteredSpans
}

func getAllSpans(trace traces.Trace) []traces.Span {
	var allSpans = make([]traces.Span, 0)
	traverseTree(trace.RootSpan, func(span traces.Span) {
		allSpans = append(allSpans, span)
	})

	return allSpans
}

type spanSelector struct {
	Filters       []filter
	PsedoClass    PseudoClass
	ChildSelector *spanSelector
}

func (ss spanSelector) MatchesFilters(span traces.Span) bool {
	for _, filter := range ss.Filters {
		if err := filter.Filter(span); err != nil {
			return false
		}
	}

	return true
}

type filterFunction func(traces.Span, string, Value) error

type filter struct {
	Property  string
	Operation filterFunction
	Value     Value
}

func (f filter) Filter(span traces.Span) error {
	return f.Operation(span, f.Property, f.Value)
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
