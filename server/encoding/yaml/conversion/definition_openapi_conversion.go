package conversion

import (
	"fmt"
	"strings"

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
		Outputs:          convertTestOutputsIntoOpenAPIOutputs(definition.Outputs),
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

func convertTestOutputsIntoOpenAPIOutputs(outputs []definition.Output) []openapi.TestOutput {
	if len(outputs) == 0 {
		return nil
	}

	res := make([]openapi.TestOutput, 0, len(outputs))

	for _, out := range outputs {
		res = append(res, openapi.TestOutput{
			Name: out.Name,
			Selector: openapi.Selector{
				Query: out.Selector,
			},
			Value: out.Value,
		})
	}

	return res
}
func convertTestSpecIntoOpenAPIObject(testSpec []definition.TestSpec) (openapi.TestSpecs, error) {
	if len(testSpec) == 0 {
		return openapi.TestSpecs{}, nil
	}

	definitions := make([]openapi.TestSpecsSpecs, 0, len(testSpec))
	for _, testSpec := range testSpec {
		var name *string
		if testSpec.Name != "" {
			name = &testSpec.Name
		}
		definitions = append(definitions, openapi.TestSpecsSpecs{
			Name: name,
			Selector: openapi.Selector{
				Query: testSpec.Selector,
			},
			Assertions: testSpec.Assertions,
		})
	}

	return openapi.TestSpecs{
		Specs: definitions,
	}, nil
}
