package value

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/expression/types"
)

type valueType int

const (
	TypeSingle valueType = iota
	TypeArray
	TypeNil
)

var Nil = Value{Type: TypeNil}

type Value struct {
	Items []types.TypedValue
	Type  valueType
}

func NewFromString(input string) Value {
	return New(types.GetTypedValue(input))
}

func New(value types.TypedValue) Value {
	return Value{Items: []types.TypedValue{value}, Type: TypeSingle}
}

func NewArrayFromStrings(inputs []string) Value {
	typedValues := make([]types.TypedValue, 0, len(inputs))
	for _, str := range inputs {
		typedValues = append(typedValues, types.TypedValue{Value: str, Type: types.GetType(str)})
	}

	return NewArray(typedValues)
}

func NewArray(values []types.TypedValue) Value {
	return Value{Items: values, Type: TypeArray}
}

func (v Value) Len() int {
	return len(v.Items)
}

func (v Value) IsArray() bool {
	return v.Type == TypeArray
}

func (v Value) Value() types.TypedValue {
	if v.Len() > 0 {
		return v.Items[0]
	}

	return types.TypedValue{}
}

func (v Value) ValueAt(index int) types.TypedValue {
	return v.Items[index]
}

func (v Value) String() string {
	if v.IsArray() {
		items := make([]string, 0, v.Len())
		for _, item := range v.Items {
			items = append(items, item.FormattedString())
		}

		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	}

	return v.Value().Value
}
