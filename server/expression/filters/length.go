package filters

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/expression/types"
)

func Length(input Value, args ...string) (Value, error) {
	if len(args) != 0 {
		return Value{}, fmt.Errorf("wrong number of args. Expected 0, got %d", len(args))
	}

	length, err := getLength(input)
	if err != nil {
		return Value{}, err
	}

	return NewValue(types.TypedValue{
		Type:  types.TypeNumber,
		Value: fmt.Sprintf("%d", length),
	}), nil
}

func getLength(input Value) (int, error) {
	if input.IsArray() {
		return input.Len(), nil
	}

	if input.Len() == 0 {
		return 0, nil
	}

	singleItem := input.Items[0]
	if singleItem.Type != types.TypeString {
		// we don't support length on types that are not string or array
		return -1, fmt.Errorf("unsupported type: expected array or string, got %s", singleItem.Type.String())
	}

	return len(singleItem.Value), nil
}
