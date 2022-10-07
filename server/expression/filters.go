package expression

import (
	"errors"
	"fmt"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

var errJsonPathUnsupportedMultipleValues = errors.New("matching multiple values from JSON is not supported yet")

type filterFn func(input string, args ...string) (string, error)

var filterFunctions = map[string]filterFn{
	"json_path": json_filter,
}

func executeFilter(input string, filterName string, args []string) (string, error) {
	fn, found := filterFunctions[filterName]
	if !found {
		return "", fmt.Errorf("filter %s was not found", filterName)
	}

	output, err := fn(input, args...)
	if err != nil {
		return "", err
	}

	return output, nil
}

func json_filter(input string, args ...string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("json_path: wrong number of args. Expected 1, got %d", len(args))
	}

	jsonPathQuery, err := jp.ParseString(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid json_path: %w", err)
	}

	jsonObject, err := oj.ParseString(input)
	if err != nil {
		return "", fmt.Errorf("json_path: invalid JSON input: %w", err)
	}

	results := jsonPathQuery.Get(jsonObject)
	if len(results) == 0 {
		return "", nil
	}

	if len(results) == 1 {
		if resultMap, ok := results[0].(map[string]any); ok {
			if len(resultMap) > 0 {
				return "", errJsonPathUnsupportedMultipleValues
			}
		}
		return fmt.Sprintf("%v", results[0]), nil
	}

	// TODO: decide what to do in case of multiple matches

	return "", errJsonPathUnsupportedMultipleValues
}
