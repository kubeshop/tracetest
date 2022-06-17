package conversion

import (
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"github.com/kubeshop/tracetest/server/model"
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
	headers := make([]model.HTTPHeader, 0, len(serviceUnderTest.Request.Headers))
	for _, header := range serviceUnderTest.Request.Headers {
		headers = append(headers, model.HTTPHeader{
			Key:   header.Key,
			Value: header.Value,
		})
	}

	auth := getAuthDefinition(serviceUnderTest.Request.Auth)

	body := serviceUnderTest.Request.Body

	if serviceUnderTest.Request.Url == "" {
		// Probably the request is empty, so return an empty trigger
		// TODO: this has to be refactored when new triggering methods are added
		return definition.TestTrigger{}
	}

	return definition.TestTrigger{
		// we only support http for now
		Type: "http",
		HTTPRequest: model.HTTPRequest{
			URL:     serviceUnderTest.Request.Url,
			Method:  model.HTTPMethod(serviceUnderTest.Request.Method),
			Headers: headers,
			Body:    body,
			Auth:    &auth,
		},
	}
}

func getAuthDefinition(auth openapi.HttpAuth) model.HTTPAuthenticator {
	switch auth.Type {
	case "basic":
		return model.HTTPAuthenticator{
			Type: "basic",
			Props: map[string]string{
				"username": auth.Basic.Username,
				"password": auth.Basic.Password,
			},
		}
	case "apiKey":
		return model.HTTPAuthenticator{
			Type: "apiKey",
			Props: map[string]string{
				"key":   auth.ApiKey.Key,
				"value": auth.ApiKey.Value,
				"in":    auth.ApiKey.In,
			},
		}
	case "bearer":
		return model.HTTPAuthenticator{
			Type: "bearer",
			Props: map[string]string{
				"token": auth.Bearer.Token,
			},
		}
	default:
		return model.HTTPAuthenticator{}
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
			Selector:   def.Selector,
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
