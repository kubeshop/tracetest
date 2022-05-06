package executor

import (
	"github.com/kubeshop/tracetest/assertions"
	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/openapi"
)

func convertAssertionsIntoTestDefinition(openapiAssertions []openapi.Assertion) assertions.TestDefinition {
	testDefinition := make(assertions.TestDefinition, 0)
	selector := ""
	for _, assertion := range openapiAssertions {
		selector = assertion.Selector
		newAssertions := make([]assertions.Assertion, 0, len(openapiAssertions))
		for _, spanAssertion := range assertion.SpanAssertions {
			newAssertion := assertions.Assertion{
				ID:         assertion.AssertionId,
				Attribute:  spanAssertion.PropertyName,
				Comparator: getComparator(spanAssertion.Operator),
				Value:      spanAssertion.ComparisonValue,
			}

			newAssertions = append(newAssertions, newAssertion)
		}
		testDefinition[assertions.SpanQuery(selector)] = newAssertions
	}
	return testDefinition
}

func getComparator(operatorName string) comparator.Comparator {
	switch operatorName {
	case "EQUALS":
		return comparator.Eq
	case "LESSTHAN":
		return comparator.Lt
	case "GREATERTHAN":
		return comparator.Gt
	case "CONTAINS":
		return comparator.Contains
	default:
		// TODO: implement comparators "NOTEQUALS", "GREATEROREQUALS", "LESSOREQUAL"
		return comparator.Eq
	}
}
