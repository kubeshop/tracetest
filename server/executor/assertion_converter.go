package executor

import (
	"fmt"

	"github.com/kubeshop/tracetest/assertions"
	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/openapi"
)

func ConvertInputTestIntoTestDefinition(test openapi.Test) (assertions.TestDefinition, error) {
	testDefinition := make(assertions.TestDefinition)
	for selector, asserts := range test.Definition.Definitions {
		newAssertions := make([]assertions.Assertion, len(asserts))
		for i, a := range asserts {
			comparator, err := getComparator(a.Comparator)
			if err != nil {
				return assertions.TestDefinition{}, err
			}

			newAssertions[i] = assertions.Assertion{
				ID:         a.Id,
				Attribute:  a.Attribute,
				Comparator: comparator,
				Value:      a.Expected,
			}
		}

		testDefinition[assertions.SpanQuery(selector)] = newAssertions
	}
	return testDefinition, nil
}

func getComparator(operatorName string) (comparator.Comparator, error) {
	operationName := "="
	switch operatorName {
	case "EQUALS":
		operationName = "="
	case "LESSTHAN":
		operationName = "<"
	case "GREATERTHAN":
		operationName = ">"
	case "CONTAINS":
		operationName = "contains"
	}

	registry := comparator.DefaultRegistry()
	comp, err := registry.Get(operationName)

	if err != nil {
		return nil, fmt.Errorf("operator %s is not supported: %w", operatorName, err)
	}

	return comp, nil
}
