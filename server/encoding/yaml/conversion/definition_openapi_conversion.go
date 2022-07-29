package conversion

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/encoding/yaml/conversion/parser"
	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"github.com/kubeshop/tracetest/server/openapi"
)

func ConvertTestDefinitionIntoOpenAPIObject(definition definition.Test) (openapi.Test, error) {
	spec, err := convertTestSpecIntoOpenAPIObject(definition.Specs)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not convert test definition: %w", err)
	}
	return openapi.Test{
		Id:               definition.Id,
		Name:             definition.Name,
		Description:      definition.Description,
		ServiceUnderTest: convertTriggerIntoServiceUnderTest(definition.Trigger),
		Specs:            spec,
	}, nil
}

func convertTriggerIntoServiceUnderTest(trigger definition.TestTrigger) openapi.Trigger {

	return openapi.Trigger{
		TriggerType: trigger.Type,
		TriggerSettings: openapi.TriggerTriggerSettings{
			Http: convertDefinitionIntoHTTPRequestOpenAPI(trigger.HTTPRequest),
			Grpc: convertDefinitionIntoGRPCOpenAPI(trigger.GRPC),
		},
	}
}

func convertDefinitionIntoGRPCOpenAPI(request definition.GRPC) openapi.GrpcRequest {
	metadata := make([]openapi.GrpcHeader, 0, len(request.Metadata))
	for _, meta := range request.Metadata {
		metadata = append(metadata, openapi.GrpcHeader{
			Key:   meta.Key,
			Value: meta.Value,
		})
	}

	return openapi.GrpcRequest{
		ProtobufFile: (request.ProtobufFile),
		Address:      (request.Address),
		Method:       (request.Method),
		Metadata:     metadata,
		Auth:         getAuthOpenAPI(request.Auth),
		Request:      (request.Request),
	}
}

func convertDefinitionIntoHTTPRequestOpenAPI(request definition.HTTPRequest) openapi.HttpRequest {
	headers := make([]openapi.HttpHeader, 0, len(request.Headers))
	for _, header := range request.Headers {
		headers = append(headers, openapi.HttpHeader{
			Key:   (header.Key),
			Value: (header.Value),
		})
	}

	return openapi.HttpRequest{
		Url:     (request.URL),
		Method:  (request.Method),
		Headers: headers,
		Body:    (request.Body),
		Auth:    getAuthOpenAPI(request.Authentication),
	}
}

func getAuthOpenAPI(auth definition.HTTPAuthentication) openapi.HttpAuth {
	switch strings.ToLower(auth.Type) {
	case "basic":
		return openapi.HttpAuth{
			Type: "basic",
			Basic: openapi.HttpAuthBasic{
				Username: auth.Basic.User,
				Password: auth.Basic.Password,
			},
		}
	case "apikey":
		return openapi.HttpAuth{
			Type: "apikey",
			ApiKey: openapi.HttpAuthApiKey{
				Key:   auth.ApiKey.Key,
				Value: auth.ApiKey.Value,
				In:    auth.ApiKey.In,
			},
		}
	case "bearer":
		return openapi.HttpAuth{
			Type: "bearer",
			Bearer: openapi.HttpAuthBearer{
				Token: auth.Bearer.Token,
			},
		}
	default:
		return openapi.HttpAuth{}
	}
}

func convertTestSpecIntoOpenAPIObject(testSpec []definition.TestSpec) (openapi.TestSpecs, error) {
	if len(testSpec) == 0 {
		return openapi.TestSpecs{}, nil
	}

	definitions := make([]openapi.TestSpecsSpecs, 0, len(testSpec))
	for _, testSpec := range testSpec {
		assertions := make([]openapi.Assertion, 0, len(testSpec.Assertions))
		for _, assertion := range testSpec.Assertions {
			assertionObject, err := convertStringIntoAssertion(assertion)
			if err != nil {
				return openapi.TestSpecs{}, err
			}
			assertions = append(assertions, assertionObject)
		}

		definitions = append(definitions, openapi.TestSpecsSpecs{
			Selector: openapi.Selector{
				Query: testSpec.Selector,
			},
			Assertions: assertions,
		})
	}

	return openapi.TestSpecs{
		Specs: definitions,
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
		Expected:   parsedAssertion.Value.String(),
	}, nil
}
