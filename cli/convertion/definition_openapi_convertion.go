package convertion

import (
	"github.com/kubeshop/tracetest/cli/definition"
	"github.com/kubeshop/tracetest/cli/openapi"
)

func ConvertTestDefinitionIntoOpenAPIObject(definition definition.Test) openapi.Test {
	return openapi.Test{
		Name:        &definition.Name,
		Description: &definition.Description,
		ServiceUnderTest: &openapi.TestServiceUnderTest{
			Request: convertHTTPRequestDefinitionIntoOpenAPIObject(definition.Trigger.HTTPRequest),
		},
		Definition: convertTestDefinitionsIntoOpenAPIObject(definition.TestDefinition),
	}
}

func convertHTTPRequestDefinitionIntoOpenAPIObject(request definition.HttpRequest) *openapi.HTTPRequest {
	headers := make([]openapi.HTTPHeader, 0, len(request.Headers))
	for _, header := range request.Headers {
		headers = append(headers, openapi.HTTPHeader{
			Key:   &header.Key,
			Value: &header.Value,
		})
	}

	return &openapi.HTTPRequest{
		Url:     &request.URL,
		Method:  &request.Method,
		Headers: headers,
		Body:    &request.Body.Raw,
		Auth: &openapi.HTTPAuth{
			Type: &request.Authentication.Type,
			ApiKey: &openapi.HTTPAuthApiKey{
				Key:   &request.Authentication.ApiKey.Key,
				Value: &request.Authentication.ApiKey.Value,
				In:    &request.Authentication.ApiKey.In,
			},
			Basic: &openapi.HTTPAuthBasic{
				Username: &request.Authentication.BasicAuth.User,
				Password: &request.Authentication.BasicAuth.Password,
			},
			Bearer: &openapi.HTTPAuthBearer{
				Token: &request.Authentication.Bearer.Token,
			},
		},
	}
}

func convertTestDefinitionsIntoOpenAPIObject(testDefinitions []definition.TestDefinition) *openapi.TestDefinition {
	definitions := make([]openapi.TestDefinitionDefinitions, 0, len(testDefinitions))
	for _, testDefinition := range testDefinitions {
		assertions := make([]openapi.Assertion, 0, len(testDefinition.Assertions))
		for _, assertion := range testDefinition.Assertions {
			assertions = append(assertions, convertStringIntoAssertion(assertion))
		}

		definitions = append(definitions, openapi.TestDefinitionDefinitions{
			Selector:   &testDefinition.Selector,
			Assertions: assertions,
		})
	}

	return &openapi.TestDefinition{
		Definitions: definitions,
	}
}

func convertStringIntoAssertion(assertion string) openapi.Assertion {
	// TODO: convert string into assertion (using a parser?)
	return openapi.Assertion{}
}
