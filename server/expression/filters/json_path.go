package filters

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/expression/value"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

func JSON_path(input value.Value, args ...string) (value.Value, error) {
	if len(args) != 1 {
		return value.Value{}, fmt.Errorf("wrong number of args. Expected 1, got %d", len(args))
	}

	if input.IsArray() {
		return value.Value{}, fmt.Errorf("cannot process array of json objects")
	}

	jsonPathQuery, err := jp.ParseString(args[0])
	if err != nil {
		return value.Value{}, fmt.Errorf("invalid json_path: %w", err)
	}

	jsonObject, err := oj.ParseString(input.Value().Value)
	if err != nil {
		return value.Value{}, fmt.Errorf("invalid JSON input: %w", err)
	}

	results := jsonPathQuery.Get(jsonObject)
	if len(results) == 0 {
		return value.Value{}, nil
	}

	if len(results) == 1 {
		result := fmt.Sprintf("%v", results[0])
		return value.NewFromString(result), nil
	}

	return value.NewArrayFromStrings(toStringSlice(results)), nil
}
