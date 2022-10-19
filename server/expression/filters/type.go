package filters

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/expression/types"
)

func Type(input Value, args ...string) (Value, error) {
	if len(args) != 0 {
		return Value{}, fmt.Errorf("wrong number of args. Expected 0, got %d", len(args))
	}

	if input.IsArray() {
		return NewValueFromString(types.TypeArray.String()), nil
	}

	if input.Len() == 0 {
		return NewValueFromString(types.TypeNil.String()), nil
	}

	item := input.ValueAt(0)
	return NewValueFromString(item.Type.String()), nil
}
