package filters

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/expression/types"
)

type valueType int

const (
	SingleItem valueType = iota
	Array
)

type Value struct {
	Items []types.TypedValue
	Type  valueType
}

func NewValueFromString(input string) Value {
	typedValue := types.TypedValue{Value: input, Type: types.GetType(input)}
	return NewValue(typedValue)
}

func NewValue(value types.TypedValue) Value {
	return Value{Items: []types.TypedValue{value}, Type: SingleItem}
}

func NewArrayValueFromStrings(inputs []string) Value {
	typedValues := make([]types.TypedValue, 0, len(inputs))
	for _, str := range inputs {
		typedValues = append(typedValues, types.TypedValue{Value: str, Type: types.GetType(str)})
	}

	return NewArrayValue(typedValues)
}

func NewArrayValue(values []types.TypedValue) Value {
	return Value{Items: values, Type: Array}
}

func (v Value) Len() int {
	return len(v.Items)
}

func (v Value) IsArray() bool {
	return v.Type == Array
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
