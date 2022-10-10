package filters

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

func JSON_path(input string, args ...string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("wrong number of args. Expected 1, got %d", len(args))
	}

	jsonPathQuery, err := jp.ParseString(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid json_path: %w", err)
	}

	jsonObject, err := oj.ParseString(input)
	if err != nil {
		return "", fmt.Errorf("invalid JSON input: %w", err)
	}

	results := jsonPathQuery.Get(jsonObject)
	if len(results) == 0 {
		return "", nil
	}

	if len(results) == 1 {
		return fmt.Sprintf("%v", results[0]), nil
	}

	return formatArray(toStringSlice(results)), nil
}
