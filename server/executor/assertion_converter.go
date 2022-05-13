package executor

import (
	"fmt"

	"github.com/kubeshop/tracetest/assertions"
	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/openapi"
)

func ConvertAssertionsIntoTestDefinition(openapiAssertions []openapi.Assertion) (assertions.TestDefinition, error) {
	testDefinition := make(assertions.TestDefinition, 0)
	selector := ""
	for _, assertion := range openapiAssertions {
		selector = assertion.Selector
		newAssertions := make([]assertions.Assertion, 0, len(openapiAssertions))
		for _, spanAssertion := range assertion.SpanAssertions {
			comparator, err := getComparator(spanAssertion.Operator)
			if err != nil {
				return assertions.TestDefinition{}, err
			}
			newAssertion := assertions.Assertion{
				ID:         assertion.AssertionId,
				Attribute:  spanAssertion.PropertyName,
				Comparator: comparator,
				Value:      spanAssertion.ComparisonValue,
			}

			newAssertions = append(newAssertions, newAssertion)
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
