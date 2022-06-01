package conversion

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/conversion/parser"
	"github.com/kubeshop/tracetest/cli/definition"
	"github.com/kubeshop/tracetest/cli/openapi"
)

var availableOperators = map[string]bool{
	"=":        true,
	"<":        true,
	">":        true,
	"!=":       true,
	">=":       true,
	"<=":       true,
	"contains": true,
}

func ConvertStringIntoOpenAPIString(in string) *string {
	if in == "" {
		return nil
	}

	return &in
}

func ConvertTestDefinitionIntoOpenAPIObject(definition definition.Test) (openapi.Test, error) {
	testDefinition, err := convertTestDefinitionsIntoOpenAPIObject(definition.TestDefinition)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not convert test definition: %w", err)
	}
	return openapi.Test{
		Name:        ConvertStringIntoOpenAPIString(definition.Name),
		Description: ConvertStringIntoOpenAPIString(definition.Description),
		ServiceUnderTest: &openapi.TestServiceUnderTest{
			Request: convertHTTPRequestDefinitionIntoOpenAPIObject(definition.Trigger.HTTPRequest),
		},
		Definition: testDefinition,
	}, nil
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

func convertTestDefinitionsIntoOpenAPIObject(testDefinitions []definition.TestDefinition) (*openapi.TestDefinition, error) {
	if len(testDefinitions) == 0 {
		return nil, nil
	}

	definitions := make([]openapi.TestDefinitionDefinitions, 0, len(testDefinitions))
	for _, testDefinition := range testDefinitions {
		assertions := make([]openapi.Assertion, 0, len(testDefinition.Assertions))
		for _, assertion := range testDefinition.Assertions {
			assertionObject, err := convertStringIntoAssertion(assertion)
			if err != nil {
				return nil, err
			}
			assertions = append(assertions, assertionObject)
		}

		definitions = append(definitions, openapi.TestDefinitionDefinitions{
			Selector:   ConvertStringIntoOpenAPIString(testDefinition.Selector),
			Assertions: assertions,
		})
	}

	return &openapi.TestDefinition{
		Definitions: definitions,
	}, nil
}

func convertStringIntoAssertion(assertion string) (openapi.Assertion, error) {
	// TODO: convert string into assertion (using a parser?)
	parsedAssertion, err := parser.ParseAssertion(assertion)
	if err != nil {
		return openapi.Assertion{}, err
	}

	// We have a bug in the parser that doesn't allow us to match operators that have more than one character. To overcome this
	// problem while the bug is not solved, we gonna validate the token here
	if _, ok := availableOperators[parsedAssertion.Operator]; !ok {
		return openapi.Assertion{}, fmt.Errorf("operator \"%s\" is not supported", parsedAssertion.Operator)
	}

	return openapi.Assertion{
		Attribute:  ConvertStringIntoOpenAPIString(parsedAssertion.Attribute),
		Comparator: ConvertStringIntoOpenAPIString(parsedAssertion.Operator),
		Expected:   ConvertStringIntoOpenAPIString(parsedAssertion.Value),
	}, nil
}
