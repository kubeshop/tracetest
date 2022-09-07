package conversion

import (
	"strings"

	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"github.com/kubeshop/tracetest/server/openapi"
)

func ConvertOpenAPITestIntoDefinitionObject(test openapi.Test) (definition.Test, error) {
	return definition.Test{
		Id:          test.Id,
		Name:        test.Name,
		Description: test.Description,
		Trigger:     convertServiceUnderTestIntoTrigger(test.ServiceUnderTest),
		Specs:       convertOpenAPITestSpecIntoSpecArray(test.Specs),
	}, nil
}

func convertServiceUnderTestIntoTrigger(trigger openapi.Trigger) definition.TestTrigger {

	return definition.TestTrigger{
		Type:        trigger.TriggerType,
		HTTPRequest: convertHTTPRequestOpenAPIIntoDefinition(trigger.TriggerSettings.Http),
		GRPC:        convertGRPCOpenAPIIntoDefinition(trigger.TriggerSettings.Grpc),
	}
}

func convertGRPCOpenAPIIntoDefinition(request openapi.GrpcRequest) definition.GRPC {
	metadata := make([]definition.GRPCHeader, 0, len(request.Metadata))
	for _, meta := range request.Metadata {
		metadata = append(metadata, definition.GRPCHeader{
			Key:   meta.Key,
			Value: meta.Value,
		})
	}

	return definition.GRPC{
		ProtobufFile: (request.ProtobufFile),
		Address:      (request.Address),
		Method:       (request.Method),
		Metadata:     metadata,
		Auth:         getAuthDefinition(request.Auth),
		Request:      (request.Request),
	}
}

func convertHTTPRequestOpenAPIIntoDefinition(request openapi.HttpRequest) definition.HTTPRequest {
	headers := make([]definition.HTTPHeader, 0, len(request.Headers))
	for _, header := range request.Headers {
		headers = append(headers, definition.HTTPHeader{
			Key:   (header.Key),
			Value: (header.Value),
		})
	}

	return definition.HTTPRequest{
		URL:            (request.Url),
		Method:         (request.Method),
		Headers:        headers,
		Body:           (request.Body),
		Authentication: getAuthDefinition(request.Auth),
	}
}

func getAuthDefinition(auth openapi.HttpAuth) definition.HTTPAuthentication {
	switch strings.ToLower(auth.Type) {
	case "basic":
		return definition.HTTPAuthentication{
			Type: "basic",
			Basic: definition.HTTPBasicAuth{
				User:     auth.Basic.Username,
				Password: auth.Basic.Password,
			},
		}
	case "apikey":
		return definition.HTTPAuthentication{
			Type: "apikey",
			ApiKey: definition.HTTPAPIKeyAuth{
				Key:   auth.ApiKey.Key,
				Value: auth.ApiKey.Value,
				In:    auth.ApiKey.In,
			},
		}
	case "bearer":
		return definition.HTTPAuthentication{
			Type: "bearer",
			Bearer: definition.HTTPBearerAuth{
				Token: auth.Bearer.Token,
			},
		}
	default:
		return definition.HTTPAuthentication{}
	}
}

func convertOpenAPITestSpecIntoSpecArray(testSpec openapi.TestSpecs) []definition.TestSpec {
	definitionArray := make([]definition.TestSpec, 0, len(testSpec.Specs))
	for _, def := range testSpec.Specs {
		newDefinition := definition.TestSpec{
			Selector:   def.Selector.Query,
			Assertions: def.Assertions,
		}
		definitionArray = append(definitionArray, newDefinition)
	}

	return definitionArray
}
