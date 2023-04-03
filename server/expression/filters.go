package expression

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/expression/filters"
	"github.com/kubeshop/tracetest/server/expression/value"
)

type filterFn func(input value.Value, args ...string) (value.Value, error)

var filterFunctions = map[string]filterFn{
	"json_path":   filters.JSON_path,
	"regex":       filters.Regex,
	"regex_group": filters.RegexGroup,
	"get_index":   filters.GetIndex,
	"length":      filters.Length,
	"type":        filters.Type,
}

func executeFilter(input value.Value, filterName string, args []string) (value.Value, error) {
	fn, found := filterFunctions[filterName]
	if !found {
		return value.Nil, resolutionError(ResourceTypeFilter, filterName, fmt.Errorf("filter not found"))
	}

	output, err := fn(input, args...)
	if err != nil {
		return value.Value{}, fmt.Errorf("%s: %w", filterName, err)
	}

	return output, nil
}
