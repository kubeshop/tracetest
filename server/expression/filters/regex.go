package filters

import (
	"fmt"
	"regexp"

	"github.com/kubeshop/tracetest/server/expression/types"
	"github.com/kubeshop/tracetest/server/expression/value"
)

func Regex(input value.Value, args ...string) (value.Value, error) {
	if len(args) != 1 {
		return value.Value{}, fmt.Errorf("wrong number of args. Expected 1, got %d", len(args))
	}

	if input.IsArray() {
		return value.Value{}, fmt.Errorf("cannot process array of json objects")
	}

	regex, err := regexp.Compile(args[0])
	if err != nil {
		return value.Value{}, fmt.Errorf("invalid regex: %w", err)
	}

	results := regex.FindAllString(input.Value().Value, -1)
	if results == nil {
		return value.NewArray([]types.TypedValue{}), nil
	}

	if len(results) == 1 {
		typedValue := types.GetTypedValue(results[0])
		return value.New(typedValue), nil
	}

	return value.NewArrayFromStrings(results), nil
}
