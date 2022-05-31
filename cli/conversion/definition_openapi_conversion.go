package conversion

import (
	"github.com/kubeshop/tracetest/cli/definition"
	"github.com/kubeshop/tracetest/cli/openapi"
)

func ConvertStringIntoOpenAPIString(in string) *string {
	if in == "" {
		return nil
	}

	return &in
}

func ConvertTestDefinitionIntoOpenAPIObject(definition definition.Test) openapi.Test {
	return openapi.Test{
		Name:        ConvertStringIntoOpenAPIString(definition.Name),
		Description: ConvertStringIntoOpenAPIString(definition.Description),
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
			Key:   ConvertStringIntoOpenAPIString(header.Key),
			Value: ConvertStringIntoOpenAPIString(header.Value),
		})
	}

	return &openapi.HTTPRequest{
		Url:     ConvertStringIntoOpenAPIString(request.URL),
		Method:  ConvertStringIntoOpenAPIString(request.Method),
		Headers: headers,
		Body:    ConvertStringIntoOpenAPIString(request.Body.Raw),
		Auth: &openapi.HTTPAuth{
			Type:   ConvertStringIntoOpenAPIString(request.Authentication.Type),
			ApiKey: getApiKeyAuthFromDefinition(request.Authentication.ApiKey),
			Basic:  getBasicAuthFromDefinition(request.Authentication.BasicAuth),
			Bearer: getBearerAuthFromDefinition(request.Authentication.Bearer),
		},
	}
}

func getApiKeyAuthFromDefinition(in definition.HTTPAPIKeyAuth) *openapi.HTTPAuthApiKey {
	if in.Key == "" && in.Value == "" && in.In == "" {
		return nil
	}

	return &openapi.HTTPAuthApiKey{
		Key:   ConvertStringIntoOpenAPIString(in.Key),
		Value: ConvertStringIntoOpenAPIString(in.Value),
		In:    ConvertStringIntoOpenAPIString(in.In),
	}
}

func getBasicAuthFromDefinition(in definition.HTTPBasicAuth) *openapi.HTTPAuthBasic {
	if in.User == "" && in.Password == "" {
		return nil
	}

	return &openapi.HTTPAuthBasic{
		Username: ConvertStringIntoOpenAPIString(in.User),
		Password: ConvertStringIntoOpenAPIString(in.Password),
	}
}

func getBearerAuthFromDefinition(in definition.HTTPBearerAuth) *openapi.HTTPAuthBearer {
	if in.Token == "" {
		return nil
	}

	return &openapi.HTTPAuthBearer{
		Token: ConvertStringIntoOpenAPIString(in.Token),
	}
}

func convertTestDefinitionsIntoOpenAPIObject(testDefinitions []definition.TestDefinition) *openapi.TestDefinition {
	if len(testDefinitions) == 0 {
		return nil
	}

	definitions := make([]openapi.TestDefinitionDefinitions, 0, len(testDefinitions))
	for _, testDefinition := range testDefinitions {
		assertions := make([]openapi.Assertion, 0, len(testDefinition.Assertions))
		for _, assertion := range testDefinition.Assertions {
			assertions = append(assertions, convertStringIntoAssertion(assertion))
		}

		definitions = append(definitions, openapi.TestDefinitionDefinitions{
			Selector:   ConvertStringIntoOpenAPIString(testDefinition.Selector),
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
