package value

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/expression/types"
)

type valueType int

const (
	TypeNil valueType = iota
	TypeSingle
	TypeArray
)

var Nil = Value{vType: TypeNil}

type Value struct {
	Items []types.TypedValue
	vType valueType
}

func NewFromString(input string) Value {
	return New(types.GetTypedValue(input))
}

func New(value types.TypedValue) Value {
	return Value{Items: []types.TypedValue{value}, vType: TypeSingle}
}

func NewArrayFromStrings(inputs []string) Value {
	typedValues := make([]types.TypedValue, 0, len(inputs))
	for _, str := range inputs {
		typedValues = append(typedValues, types.TypedValue{Value: str, Type: types.GetType(str)})
	}

	return NewArray(typedValues)
}

func NewArray(values []types.TypedValue) Value {
	return Value{Items: values, vType: TypeArray}
}

func (v Value) Len() int {
	return len(v.Items)
}

func (v Value) IsArray() bool {
	return v.vType == TypeArray
}

func (v Value) Value() types.TypedValue {
	if v.Len() > 0 {
		return v.Items[0]
	}

	return types.TypedValue{}
}

func (v Value) Type() types.Type {
	if v.IsArray() {
		return types.TypeArray
	}

	return v.Value().Type
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

func (v Value) UnescappedString() string {
	output := v.String()
	output = strings.ReplaceAll(output, `\'`, `'`)
	output = strings.ReplaceAll(output, `\"`, `"`)

	return output
}
