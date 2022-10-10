package filters

import (
	"fmt"
	"regexp"
)

func RegexGroup(input string, args ...string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("wrong number of args. Expected 1, got %d", len(args))
	}

	regex, err := regexp.Compile(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid regex: %w", err)
	}

	groups := regex.FindAllStringSubmatch(input, -1)
	if groups == nil {
		return "[]", nil
	}

	output := make([]string, 0)
	for _, group := range groups {
		output = append(output, group[1:]...)
	}

	if len(output) == 1 {
		return output[0], nil
	}

	return formatArray(output), nil
}
