package yaml

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

type TestSpecs []TestSpec

func (ts TestSpecs) Model() model.OrderedMap[model.SpanQuery, model.NamedAssertions] {
	mts := model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}
	for _, spec := range ts {
		assertions := make([]model.Assertion, 0, len(spec.Assertions))
		for _, a := range spec.Assertions {
			assertions = append(assertions, model.Assertion(a))
		}

		mts, _ = mts.Add(model.SpanQuery(spec.Selector), model.NamedAssertions{
			Name:       spec.Name,
			Assertions: assertions,
		})
	}
	return mts
}

type Outputs []Output

func (outs Outputs) Model() model.OrderedMap[string, model.Output] {
	mos := model.OrderedMap[string, model.Output]{}
	for _, output := range outs {
		mos, _ = mos.Add(output.Name, model.Output{
			Selector: model.SpanQuery(output.Selector),
			Value:    output.Value,
		})
	}
	return mos
}

type Test struct {
	ID          string      `mapstructure:"id"`
	Name        string      `mapstructure:"name"`
	Description string      `mapstructure:"description" yaml:",omitempty"`
	Trigger     TestTrigger `mapstructure:"trigger"`
	Specs       TestSpecs   `mapstructure:"specs" yaml:",omitempty"`
	Outputs     Outputs     `mapstructure:"outputs,omitempty" yaml:",omitempty"`
}

func (t Test) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("test name cannot be empty")
	}

	if err := t.Trigger.Validate(); err != nil {
		return fmt.Errorf("test trigger must be valid: %w", err)
	}

	return nil
}

type TestTrigger struct {
	Type        string      `mapstructure:"type"`
	HTTPRequest HTTPRequest `mapstructure:"httpRequest" yaml:"httpRequest,omitempty"`
	GRPC        GRPC        `mapstructure:"grpc" yaml:"grpc,omitempty"`
}

func (t TestTrigger) Model() model.Trigger {
	mt := model.Trigger{
		Type: model.TriggerType(t.Type),
	}

	switch t.Type {
	case "http":
		hr := t.HTTPRequest
		mt.HTTP = &model.HTTPRequest{
			Method:  model.HTTPMethod(hr.Method),
			URL:     hr.URL,
			Headers: hr.Headers.Model(),
			Body:    hr.Body,
			Auth:    hr.Authentication.Model(),
		}
	}
	return mt
}

func (t TestTrigger) Validate() error {
	switch t.Type {
	case "http":
		if err := t.HTTPRequest.Validate(); err != nil {
			return fmt.Errorf("http request must be valid: %w", err)
		}
	case "grpc":
		if err := t.GRPC.Validate(); err != nil {
			return fmt.Errorf("grpc request must be valid: %w", err)
		}
	case "":
		return fmt.Errorf("type cannot be empty")
	default:
		return fmt.Errorf("type \"%s\" is not supported", t.Type)
	}

	return nil
}

type Output struct {
	Name     string `mapstructure:"name"`
	Selector string `mapstructure:"selector"`
	Value    string `mapstructure:"value"`
}

type TestSpec struct {
	Name       string   `mapstructure:"name" yaml:",omitempty"`
	Selector   string   `mapstructure:"selector"`
	Assertions []string `mapstructure:"assertions"`
}

func (t Test) Model() model.Test {
	mt := model.Test{
		ID:               id.ID(t.ID),
		Name:             t.Name,
		Description:      t.Description,
		ServiceUnderTest: t.Trigger.Model(),
		Specs:            t.Specs.Model(),
		Outputs:          t.Outputs.Model(),
	}

	return mt
}

func GetTestFromOpenapiObject(test openapi.Test) (Test, error) {
	return Test{
		ID:          test.Id,
		Name:        test.Name,
		Description: test.Description,
		Trigger:     convertServiceUnderTestIntoTrigger(test.ServiceUnderTest),
		Specs:       convertOpenAPITestSpecIntoSpecArray(test.Specs),
	}, nil
}

func convertServiceUnderTestIntoTrigger(trigger openapi.Trigger) TestTrigger {

	return TestTrigger{
		Type:        trigger.TriggerType,
		HTTPRequest: convertHTTPRequestOpenAPIIntoDefinition(trigger.TriggerSettings.Http),
		GRPC:        convertGRPCOpenAPIIntoDefinition(trigger.TriggerSettings.Grpc),
	}
}

func convertGRPCOpenAPIIntoDefinition(request openapi.GrpcRequest) GRPC {
	metadata := make([]GRPCHeader, 0, len(request.Metadata))
	for _, meta := range request.Metadata {
		metadata = append(metadata, GRPCHeader{
			Key:   meta.Key,
			Value: meta.Value,
		})
	}

	return GRPC{
		ProtobufFile: (request.ProtobufFile),
		Address:      (request.Address),
		Method:       (request.Method),
		Metadata:     metadata,
		Auth:         getAuthDefinition(request.Auth),
		Request:      (request.Request),
	}
}

func convertHTTPRequestOpenAPIIntoDefinition(request openapi.HttpRequest) HTTPRequest {
	headers := make([]HTTPHeader, 0, len(request.Headers))
	for _, header := range request.Headers {
		headers = append(headers, HTTPHeader{
			Key:   (header.Key),
			Value: (header.Value),
		})
	}

	return HTTPRequest{
		URL:            (request.Url),
		Method:         (request.Method),
		Headers:        headers,
		Body:           (request.Body),
		Authentication: getAuthDefinition(request.Auth),
	}
}

func getAuthDefinition(auth openapi.HttpAuth) HTTPAuthentication {
	switch strings.ToLower(auth.Type) {
	case "basic":
		return HTTPAuthentication{
			Type: "basic",
			Basic: HTTPBasicAuth{
				User:     auth.Basic.Username,
				Password: auth.Basic.Password,
			},
		}
	case "apikey":
		return HTTPAuthentication{
			Type: "apikey",
			ApiKey: HTTPAPIKeyAuth{
				Key:   auth.ApiKey.Key,
				Value: auth.ApiKey.Value,
				In:    auth.ApiKey.In,
			},
		}
	case "bearer":
		return HTTPAuthentication{
			Type: "bearer",
			Bearer: HTTPBearerAuth{
				Token: auth.Bearer.Token,
			},
		}
	default:
		return HTTPAuthentication{}
	}
}

func convertOpenAPITestSpecIntoSpecArray(testSpec openapi.TestSpecs) []TestSpec {
	definitionArray := make([]TestSpec, 0, len(testSpec.Specs))
	for _, def := range testSpec.Specs {
		name := ""
		if def.Name != nil {
			name = *def.Name
		}
		newDefinition := TestSpec{
			Name:       name,
			Selector:   def.Selector.Query,
			Assertions: def.Assertions,
		}
		definitionArray = append(definitionArray, newDefinition)
	}

	return definitionArray
}
