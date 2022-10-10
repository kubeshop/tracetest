package filters

import (
	"fmt"
	"regexp"
)

func Regex(input string, args ...string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("wrong number of args. Expected 1, got %d", len(args))
	}

	regex, err := regexp.Compile(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid regex: %w", err)
	}

	results := regex.FindAllString(input, -1)
	if results == nil {
		return "[]", nil
	}

	if len(results) == 1 {
		return results[0], nil
	}

	return formatArray(results), nil
}
