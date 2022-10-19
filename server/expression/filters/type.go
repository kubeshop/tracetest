package filters

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/expression/types"
	"github.com/kubeshop/tracetest/server/expression/value"
)

func Type(input value.Value, args ...string) (value.Value, error) {
	if len(args) != 0 {
		return value.Value{}, fmt.Errorf("wrong number of args. Expected 0, got %d", len(args))
	}

	if input.IsArray() {
		return value.NewFromString(types.TypeArray.String()), nil
	}

	if input.Len() == 0 {
		return value.NewFromString(types.TypeNil.String()), nil
	}

	item := input.ValueAt(0)
	return value.NewFromString(item.Type.String()), nil
}
