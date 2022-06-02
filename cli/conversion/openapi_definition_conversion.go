package conversion

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/definition"
	"github.com/kubeshop/tracetest/cli/openapi"
)

func ConvertOpenAPITestIntoDefinitionObject(test openapi.Test) (definition.Test, error) {
	trigger := convertServiceUnderTestIntoTrigger(test.ServiceUnderTest)
	testDefinition := convertOpenAPITestDefinitionIntoDefinitionArray(test.Definition)
	description := ""
	if test.Description != nil {
		description = *test.Description
	}

	return definition.Test{
		Id:             *test.Id,
		Name:           *test.Name,
		Description:    description,
		Trigger:        trigger,
		TestDefinition: testDefinition,
	}, nil
}

func convertServiceUnderTestIntoTrigger(serviceUnderTest *openapi.TestServiceUnderTest) definition.TestTrigger {
	if serviceUnderTest == nil || serviceUnderTest.Request == nil {
		return definition.TestTrigger{}
	}

	headers := make([]definition.HTTPHeader, 0, len(serviceUnderTest.Request.Headers))
	for _, header := range serviceUnderTest.Request.Headers {
		headers = append(headers, definition.HTTPHeader{
			Key:   *header.Key,
			Value: *header.Value,
		})
	}

	auth := getAuthDefinition(serviceUnderTest.Request.Auth)

	body := definition.HTTPBody{}
	if serviceUnderTest.Request.Body != nil {
		body = definition.HTTPBody{
			// we only support raw for now
			Type: "raw",
			Raw:  *serviceUnderTest.Request.Body,
		}
	}

	return definition.TestTrigger{
		// we only support http for now
		Type: "http",
		HTTPRequest: definition.HttpRequest{
			URL:            *serviceUnderTest.Request.Url,
			Method:         *serviceUnderTest.Request.Method,
			Headers:        headers,
			Body:           body,
			Authentication: auth,
		},
	}
}

func getAuthDefinition(auth *openapi.HTTPAuth) definition.HTTPAuthentication {
	if auth == nil || auth.Type == nil {
		return definition.HTTPAuthentication{}
	}

	switch *auth.Type {
	case "basic":
		return definition.HTTPAuthentication{
			Type: "basic",
			BasicAuth: definition.HTTPBasicAuth{
				User:     *auth.Basic.Username,
				Password: *auth.Basic.Password,
			},
		}
	case "apikey":
		return definition.HTTPAuthentication{
			Type: "apikey",
			ApiKey: definition.HTTPAPIKeyAuth{
				Key:   *auth.ApiKey.Key,
				Value: *auth.ApiKey.Value,
				In:    *auth.ApiKey.In,
			},
		}
	case "bearer":
		return definition.HTTPAuthentication{
			Type: "bearer",
			Bearer: definition.HTTPBearerAuth{
				Token: *auth.Bearer.Token,
			},
		}
	default:
		return definition.HTTPAuthentication{}
	}
}

func convertOpenAPITestDefinitionIntoDefinitionArray(testDefinition *openapi.TestDefinition) []definition.TestDefinition {
	if testDefinition == nil {
		return []definition.TestDefinition{}
	}

	definitionArray := make([]definition.TestDefinition, 0, len(testDefinition.Definitions))
	for _, def := range testDefinition.Definitions {
		assertions := make([]string, 0, len(def.Assertions))
		for _, assertion := range def.Assertions {
			assertionString := fmt.Sprintf("%s %s %s", *assertion.Attribute, *assertion.Comparator, *assertion.Expected)
			assertions = append(assertions, assertionString)
		}

		newDefinition := definition.TestDefinition{
			Selector:   *def.Selector,
			Assertions: assertions,
		}
		definitionArray = append(definitionArray, newDefinition)
	}

	return definitionArray
}
