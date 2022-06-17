package conversion

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/encoding/yaml/conversion/parser"
	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

func ConvertTestDefinitionIntoOpenAPIObject(definition definition.Test) (openapi.Test, error) {
	testDefinition, err := convertTestDefinitionsIntoOpenAPIObject(definition.TestDefinition)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not convert test definition: %w", err)
	}
	return openapi.Test{
		Id:          definition.Id,
		Name:        definition.Name,
		Description: definition.Description,
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Request: convertHTTPRequestDefinitionIntoOpenAPIObject(definition.Trigger.HTTPRequest),
		},
		Definition: testDefinition,
	}, nil
}

func convertHTTPRequestDefinitionIntoOpenAPIObject(request model.HTTPRequest) openapi.HttpRequest {
	headers := make([]openapi.HttpHeader, 0, len(request.Headers))
	for _, header := range request.Headers {
		headers = append(headers, openapi.HttpHeader{
			Key:   header.Key,
			Value: header.Value,
		})
	}

	return openapi.HttpRequest{
		Url:     request.URL,
		Method:  string(request.Method),
		Headers: headers,
		Body:    request.Body,
		Auth: openapi.HttpAuth{
			Type:   request.Auth.Type,
			ApiKey: getApiKeyAuthFromDefinition(request.Auth),
			Basic:  getBasicAuthFromDefinition(request.Auth),
			Bearer: getBearerAuthFromDefinition(request.Auth),
		},
	}
}

func getApiKeyAuthFromDefinition(in *model.HTTPAuthenticator) openapi.HttpAuthApiKey {
	return openapi.HttpAuthApiKey{
		Key:   in.Props["key"],
		Value: in.Props["value"],
		In:    in.Props["in"],
	}
}

func getBasicAuthFromDefinition(in *model.HTTPAuthenticator) openapi.HttpAuthBasic {
	return openapi.HttpAuthBasic{
		Username: in.Props["username"],
		Password: in.Props["password"],
	}
}

func getBearerAuthFromDefinition(in *model.HTTPAuthenticator) openapi.HttpAuthBearer {
	return openapi.HttpAuthBearer{
		Token: in.Props["token"],
	}
}

func convertTestDefinitionsIntoOpenAPIObject(testDefinitions []definition.TestDefinition) (openapi.TestDefinition, error) {
	if len(testDefinitions) == 0 {
		return openapi.TestDefinition{}, nil
	}

	definitions := make([]openapi.TestDefinitionDefinitions, 0, len(testDefinitions))
	for _, testDefinition := range testDefinitions {
		assertions := make([]openapi.Assertion, 0, len(testDefinition.Assertions))
		for _, assertion := range testDefinition.Assertions {
			assertionObject, err := convertStringIntoAssertion(assertion)
			if err != nil {
				return openapi.TestDefinition{}, err
			}
			assertions = append(assertions, assertionObject)
		}

		definitions = append(definitions, openapi.TestDefinitionDefinitions{
			Selector:   testDefinition.Selector,
			Assertions: assertions,
		})
	}

	return openapi.TestDefinition{
		Definitions: definitions,
	}, nil
}

func convertStringIntoAssertion(assertion string) (openapi.Assertion, error) {
	// TODO: convert string into assertion (using a parser?)
	parsedAssertion, err := parser.ParseAssertion(assertion)
	if err != nil {
		return openapi.Assertion{}, err
	}

	return openapi.Assertion{
		Attribute:  parsedAssertion.Attribute,
		Comparator: parsedAssertion.Operator,
		Expected:   parsedAssertion.Value,
	}, nil
}
