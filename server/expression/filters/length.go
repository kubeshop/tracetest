package filters

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/expression/types"
)

func Length(input Value, args ...string) (Value, error) {
	if len(args) != 0 {
		return Value{}, fmt.Errorf("wrong number of args. Expected 0, got %d", len(args))
	}

	if len(input) != 1 {
		return Value{}, fmt.Errorf("wrong input: Expected string, not an array")
	}

	if input[0].Type != types.TypeString {
		return Value{}, fmt.Errorf("wrong input type: Expected string, not %s", input[0].Type.String())
	}

	return NewValue(types.TypedValue{
		Type:  types.TypeNumber,
		Value: fmt.Sprintf("%d", len(input[0].Value)),
	}), nil
}
