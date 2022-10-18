package filters

import (
	"fmt"
	"regexp"

	"github.com/kubeshop/tracetest/server/expression/types"
)

func Regex(input Value, args ...string) (Value, error) {
	if len(args) != 1 {
		return Value{}, fmt.Errorf("wrong number of args. Expected 1, got %d", len(args))
	}

	if input.IsArray() {
		return Value{}, fmt.Errorf("cannot process array of json objects")
	}

	regex, err := regexp.Compile(args[0])
	if err != nil {
		return Value{}, fmt.Errorf("invalid regex: %w", err)
	}

	results := regex.FindAllString(input.Value().Value, -1)
	if results == nil {
		return NewArrayValue([]types.TypedValue{}), nil
	}

	if len(results) == 1 {
		typedValue := types.GetTypedValue(results[0])
		return NewValue(typedValue), nil
	}

	return NewArrayValueFromStrings(results), nil
}
