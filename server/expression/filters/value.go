package filters

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/expression/types"
)

type Value []types.TypedValue

func NewValueFromString(input string) Value {
	typedValue := types.TypedValue{Value: input, Type: types.GetType(input)}
	return NewValue(typedValue)
}

func NewValue(value types.TypedValue) Value {
	return Value([]types.TypedValue{value})
}

func NewArrayValueFromStrings(inputs []string) Value {
	typedValues := make([]types.TypedValue, 0, len(inputs))
	for _, str := range inputs {
		typedValues = append(typedValues, types.TypedValue{Value: str, Type: types.GetType(str)})
	}

	return Value(typedValues)
}

func NewArrayValue(values []types.TypedValue) Value {
	return Value(values)
}

func (v Value) IsArray() bool {
	return len(v) > 1
}

func (v Value) Value() types.TypedValue {
	if len(v) > 0 {
		return v[0]
	}

	return types.TypedValue{}
}

func (v Value) ValueAt(index int) types.TypedValue {
	return v[index]
}

func (v Value) String() string {
	if v.IsArray() {
		items := make([]string, 0, len(v))
		for _, item := range v {
			items = append(items, item.FormattedString())
		}

		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	}

	return v.Value().Value
}
