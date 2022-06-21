package conversion

import (
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"github.com/kubeshop/tracetest/server/openapi"
)

func ConvertOpenAPITestIntoDefinitionObject(test openapi.Test) (definition.Test, error) {
	trigger := convertServiceUnderTestIntoTrigger(test.ServiceUnderTest)
	testDefinition := convertOpenAPITestDefinitionIntoDefinitionArray(test.Definition)
	description := test.Description

	return definition.Test{
		Id:             test.Id,
		Name:           test.Name,
		Description:    description,
		Trigger:        trigger,
		TestDefinition: testDefinition,
	}, nil
}

func convertServiceUnderTestIntoTrigger(serviceUnderTest openapi.TestServiceUnderTest) definition.TestTrigger {
	return definition.TestTrigger{
		// we only support http for now
		Type:        "http",
		HTTPRequest: serviceUnderTest.Request,
	}
}

func convertOpenAPITestDefinitionIntoDefinitionArray(testDefinition openapi.TestDefinition) []definition.TestDefinition {
	definitionArray := make([]definition.TestDefinition, 0, len(testDefinition.Definitions))
	for _, def := range testDefinition.Definitions {
		assertions := make([]string, 0, len(def.Assertions))
		for _, assertion := range def.Assertions {
			assertionFormat := `%s %s "%s"`
			if isNumber(assertion.Expected) {
				assertionFormat = "%s %s %s"
			}
			assertionString := fmt.Sprintf(assertionFormat, assertion.Attribute, assertion.Comparator, assertion.Expected)
			assertions = append(assertions, assertionString)
		}

		newDefinition := definition.TestDefinition{
			Selector:   def.Selector.Query,
			Assertions: assertions,
		}
		definitionArray = append(definitionArray, newDefinition)
	}

	return definitionArray
}

func isNumber(in string) bool {
	if _, err := strconv.Atoi(in); err == nil {
		return true
	}

	if _, err := strconv.ParseFloat(in, 64); err == nil {
		return true
	}

	return false
}
