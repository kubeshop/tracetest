package conversion

import (
	"fmt"
	"strconv"
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
		assertions := make([]string, 0, len(def.Assertions))
		for _, assertion := range def.Assertions {
			assertionFormat := `%s %s "%s"`
			if isNumber(assertion.Expected) {
				assertionFormat = "%s %s %s"
			}
			assertionString := fmt.Sprintf(assertionFormat, assertion.Attribute, assertion.Comparator, assertion.Expected)
			assertions = append(assertions, assertionString)
		}

		newDefinition := definition.TestSpec{
			Selector:   def.Selector.Query,
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
