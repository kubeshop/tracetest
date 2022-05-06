package assertions

import (
	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/openapi"
)

func convertAssertionsIntoTestDefinition(assertions []openapi.Assertion) TestDefinition {
	testDefinition := make(TestDefinition, 0)
	selector := ""
	for _, assertion := range assertions {
		selector = assertion.Selector
		newAssertions := make([]Assertion, 0, len(assertions))
		for _, spanAssertion := range assertion.SpanAssertions {
			newAssertion := Assertion{
				ID:         assertion.AssertionId,
				Attribute:  spanAssertion.PropertyName,
				Comparator: getComparator(spanAssertion.Operator),
				Value:      spanAssertion.ComparisonValue,
			}

			newAssertions = append(newAssertions, newAssertion)
		}
		testDefinition[SpanQuery(selector)] = newAssertions
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
