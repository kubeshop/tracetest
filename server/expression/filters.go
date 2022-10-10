package expression

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/expression/filters"
)

type filterFn func(input string, args ...string) (string, error)

var filterFunctions = map[string]filterFn{
	"json_path":   filters.JSON_path,
	"regex":       filters.Regex,
	"regex_group": filters.RegexGroup,
}

func executeFilter(input string, filterName string, args []string) (string, error) {
	fn, found := filterFunctions[filterName]
	if !found {
		return "", fmt.Errorf("filter %s was not found", filterName)
	}

	output, err := fn(input, args...)
	if err != nil {
		return "", fmt.Errorf("%s: %w", filterName, err)
	}

	return output, nil
}
