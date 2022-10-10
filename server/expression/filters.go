package expression

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/expression/filters"
)

type filterFn func(input filters.Value, args ...string) (filters.Value, error)

var filterFunctions = map[string]filterFn{
	"json_path":   filters.JSON_path,
	"regex":       filters.Regex,
	"regex_group": filters.RegexGroup,
	"get_index":   filters.GetIndex,
}

func executeFilter(input filters.Value, filterName string, args []string) (filters.Value, error) {
	fn, found := filterFunctions[filterName]
	if !found {
		return filters.Value{}, fmt.Errorf("filter %s was not found", filterName)
	}

	output, err := fn(input, args...)
	if err != nil {
		return filters.Value{}, fmt.Errorf("%s: %w", filterName, err)
	}

	return output, nil
}
