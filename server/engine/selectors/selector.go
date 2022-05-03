package selectors

import "github.com/kubeshop/tracetest/engine/operations"

type Selector struct {
	spanSelectors []spanSelector
}

type spanSelector struct {
	Filters       []filter
	PsedoClass    pseudoClass
	ChildSelector *spanSelector
}

type filter struct {
	Property  string
	Operation operations.OperationFunction
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
