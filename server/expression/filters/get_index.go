package filters

import (
	"fmt"
	"strconv"
)

func GetIndex(input Value, args ...string) (Value, error) {
	if len(args) != 1 {
		return Value{}, fmt.Errorf("wrong number of args. Expected 1, got %d", len(args))
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		return Value{}, fmt.Errorf("index must be an integer: %w", err)
	}

	if index < 0 || index >= len(input) {
		return Value{}, fmt.Errorf("index out of boundaries: %d out of %d", index, len(input))
	}

	value := input.ValueAt(index)
	return NewValue(value), nil
}
