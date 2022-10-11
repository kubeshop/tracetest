package filters

import (
	"fmt"
	"regexp"

	"github.com/kubeshop/tracetest/server/expression/types"
)

func RegexGroup(input Value, args ...string) (Value, error) {
	if len(args) != 1 {
		return Value{}, fmt.Errorf("wrong number of args. Expected 1, got %d", len(args))
	}

	if len(input) != 1 {
		return Value{}, fmt.Errorf("cannot process array of json objects")
	}

	regex, err := regexp.Compile(args[0])
	if err != nil {
		return Value{}, fmt.Errorf("invalid regex: %w", err)
	}

	groups := regex.FindAllStringSubmatch(input.Value().Value, -1)
	if groups == nil {
		return NewArrayValue([]types.TypedValue{}), nil
	}

	output := make([]string, 0)
	for _, group := range groups {
		output = append(output, group[1:]...)
	}

	if len(output) == 1 {
		typedValue := types.GetTypedValue(output[0])
		return NewValue(typedValue), nil
	}

	return NewArrayValueFromStrings(output), nil
}
