package filters

import (
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/server/expression/value"
)

func GetIndex(input value.Value, args ...string) (value.Value, error) {
	if len(args) != 1 {
		return value.Value{}, fmt.Errorf("wrong number of args. Expected 1, got %d", len(args))
	}

	index, err := getIndex(input, args...)
	if err != nil {
		return value.Value{}, err
	}

	if index < 0 || index >= input.Len() {
		return value.Value{}, fmt.Errorf("index out of boundaries: %d out of %d", index, input.Len())
	}

	v := input.ValueAt(index)
	return value.New(v), nil
}

func getIndex(input value.Value, args ...string) (int, error) {
	if args[0] == "last" {
		return input.Len() - 1, nil
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		return -1, fmt.Errorf("index must be an integer")
	}

	return index, nil
}
