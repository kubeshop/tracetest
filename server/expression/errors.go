package expression

import (
	"errors"
	"fmt"
)

var ErrExpressionResolution error = errors.New("error when resolving expression")

type ResourceType string

var (
	ResourceTypeAttribute           ResourceType = "attribute"
	ResourceTypeEnvironmentVariable ResourceType = "environment variable"
	ResourceTypeFunction            ResourceType = "function"
	ResourceTypeArrayItem           ResourceType = "array item"
	ResourceTypeFunctionArgument    ResourceType = "function argument"
	ResourceTypeFilter              ResourceType = "filter"
	ResourceTypeOperator            ResourceType = "operator"
)

type ResolutionError struct {
	Type     ResourceType
	Name     string
	InnerErr error
}

func resolutionError(typ ResourceType, name string, innerErr error) error {
	newError := fmt.Errorf(`%w: %s "%s"`, ErrExpressionResolution, typ, name)

	if innerErr != nil {
		newError = fmt.Errorf("%w: %s", newError, innerErr.Error())
	}

	return newError
}
