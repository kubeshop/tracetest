package selectors

import (
	"fmt"

	"github.com/kubeshop/tracetest/traces"
)

type Selector struct {
	spanSelectors []spanSelector
}

type spanSelector struct {
	Filters       []filter
	PsedoClass    *pseudoClass
	ChildSelector *spanSelector
}

type filterFunction func(traces.Span, string, Value) error

type filter struct {
	Property  string
	Operation filterFunction
	Value     Value
}

type pseudoClass struct {
	Name     string
	Argument Value
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
	default:
		return "null"
	}
}
