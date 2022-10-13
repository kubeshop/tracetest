package filters

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/expression/types"
)

func Count(input Value, args ...string) (Value, error) {
	if len(args) != 0 {
		return Value{}, fmt.Errorf("wrong number of args. Expected 0, got %d", len(args))
	}

	return NewValue(types.TypedValue{
		Type:  types.TypeNumber,
		Value: fmt.Sprintf("%d", len(input)),
	}), nil
}
