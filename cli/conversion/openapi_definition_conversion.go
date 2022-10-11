package conversion

import (
	"github.com/kubeshop/tracetest/cli/definition"
	"github.com/kubeshop/tracetest/cli/openapi"
)

func ConvertOpenapiStringIntoString(in *string) string {
	if in == nil {
		return ""
	}

	return *in
}

func ConvertOpenAPITestIntoSpecObject(test openapi.Test) (definition.Test, error) {
	trigger := convertServiceUnderTestIntoTrigger(test.ServiceUnderTest)
	testSpec := convertOpenAPITestSpecIntoSpecArray(test.Specs)
	description := ""
	if test.Description != nil {
		description = *test.Description
	}

	return definition.Test{
		Id:          *test.Id,
		Name:        *test.Name,
		Description: description,
		Trigger:     trigger,
		Specs:       testSpec,
	}, nil
}

func convertServiceUnderTestIntoTrigger(trigger *openapi.Trigger) definition.TestTrigger {
	if trigger == nil || trigger.TriggerSettings == nil {
		return definition.TestTrigger{}
	}

	return definition.TestTrigger{
		Type:        *trigger.TriggerType,
		HTTPRequest: convertHTTPRequestOpenAPIIntoDefinition(trigger.TriggerSettings.Http),
		GRPC:        convertGRPCOpenAPIIntoDefinition(trigger.TriggerSettings.Grpc),
	}
}

func convertGRPCOpenAPIIntoDefinition(request *openapi.GRPCRequest) definition.GrpcRequest {
	if request == nil {
		return definition.GrpcRequest{}
	}

	metadata := make([]definition.GRPCHeader, 0, len(request.Metadata))
	for _, meta := range request.Metadata {
		metadata = append(metadata, definition.GRPCHeader{
			Key:   ConvertOpenapiStringIntoString(meta.Key),
			Value: ConvertOpenapiStringIntoString(meta.Value),
		})
	}

	return definition.GrpcRequest{
		ProtobufFile: ConvertOpenapiStringIntoString(request.ProtobufFile),
		Address:      ConvertOpenapiStringIntoString(request.Address),
		Method:       ConvertOpenapiStringIntoString(request.Method),
		Metadata:     metadata,
		Auth:         getAuthDefinition(request.Auth),
		Request:      ConvertOpenapiStringIntoString(request.Request),
	}
}

func convertHTTPRequestOpenAPIIntoDefinition(request *openapi.HTTPRequest) definition.HttpRequest {
	if request == nil {
		return definition.HttpRequest{}
	}

	headers := make([]definition.HTTPHeader, 0, len(request.Headers))
	for _, header := range request.Headers {
		headers = append(headers, definition.HTTPHeader{
			Key:   ConvertOpenapiStringIntoString(header.Key),
			Value: ConvertOpenapiStringIntoString(header.Value),
		})
	}

	return definition.HttpRequest{
		URL:            ConvertOpenapiStringIntoString(request.Url),
		Method:         ConvertOpenapiStringIntoString(request.Method),
		Headers:        headers,
		Body:           ConvertOpenapiStringIntoString(request.Body),
		Authentication: getAuthDefinition(request.Auth),
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
			Basic: definition.HTTPBasicAuth{
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

func convertOpenAPITestSpecIntoSpecArray(testSpec *openapi.TestSpecs) []definition.TestSpec {
	if testSpec == nil {
		return []definition.TestSpec{}
	}

	definitionArray := make([]definition.TestSpec, 0, len(testSpec.Specs))
	for _, def := range testSpec.Specs {
		newDefinition := definition.TestSpec{
			Selector:   *def.Selector.Query,
			Assertions: def.Assertions,
		}
		definitionArray = append(definitionArray, newDefinition)
	}

	return definitionArray
}
