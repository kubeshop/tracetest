package conversion_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/encoding/yaml/conversion"
	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefinitionToOpenAPIConversion(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          definition.Test
		ExpectedOutput openapi.Test
	}{
		{
			Name: "Should_be_able_to_convert_request_with_no_authentication_and_no_body",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: openapi.HttpRequest{
						Url:    "http://localhost:1234",
						Method: "POST",
						Headers: []openapi.HttpHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Auth: openapi.HttpAuth{},
						Body: "",
					},
				},
			},
			ExpectedOutput: openapi.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: openapi.TestServiceUnderTest{
					Request: openapi.HttpRequest{
						Url:    "http://localhost:1234",
						Method: "POST",
						Headers: []openapi.HttpHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: "",
						Auth: openapi.HttpAuth{},
					},
				},
			},
		},
		{
			Name: "Should_be_able_to_convert_request_with_basic_auth",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: openapi.HttpRequest{
						Url:     "http://localhost:1234",
						Method:  "POST",
						Headers: []openapi.HttpHeader{},
						Body:    "",
						Auth: openapi.HttpAuth{
							Type: "basic",
							Basic: openapi.HttpAuthBasic{
								Username: "matheus",
								Password: "pikachu",
							},
						},
					},
				},
			},
			ExpectedOutput: openapi.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: openapi.TestServiceUnderTest{
					Request: openapi.HttpRequest{
						Url:     "http://localhost:1234",
						Method:  "POST",
						Headers: []openapi.HttpHeader{},
						Body:    "",
						Auth: openapi.HttpAuth{
							Type: "basic",
							Basic: openapi.HttpAuthBasic{
								Username: "matheus",
								Password: "pikachu",
							},
						},
					},
				},
			},
		},
		{
			Name: "Should_be_able_to_convert_request_with_api_key_auth",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: openapi.HttpRequest{
						Url:     "http://localhost:1234",
						Method:  "POST",
						Headers: []openapi.HttpHeader{},
						Body:    "",
						Auth: openapi.HttpAuth{
							Type: "apiKey",
							ApiKey: openapi.HttpAuthApiKey{
								Key:   "X-Key",
								Value: "my-api-key",
								In:    "header",
							},
						},
					},
				},
			},
			ExpectedOutput: openapi.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: openapi.TestServiceUnderTest{
					Request: openapi.HttpRequest{
						Url:     "http://localhost:1234",
						Method:  "POST",
						Headers: []openapi.HttpHeader{},
						Body:    "",
						Auth: openapi.HttpAuth{
							Type: "apiKey",
							ApiKey: openapi.HttpAuthApiKey{
								Key:   "X-Key",
								Value: "my-api-key",
								In:    "header",
							},
						},
					},
				},
			},
		},
		{
			Name: "Should_be_able_to_convert_request_with_bearer_auth",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: openapi.HttpRequest{
						Url:     "http://localhost:1234",
						Method:  "POST",
						Headers: []openapi.HttpHeader{},
						Body:    "",
						Auth: openapi.HttpAuth{
							Type: "bearer",
							Bearer: openapi.HttpAuthBearer{
								Token: "my-token",
							},
						},
					},
				},
			},
			ExpectedOutput: openapi.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: openapi.TestServiceUnderTest{
					Request: openapi.HttpRequest{
						Url:     "http://localhost:1234",
						Method:  "POST",
						Headers: []openapi.HttpHeader{},
						Body:    "",
						Auth: openapi.HttpAuth{
							Type: "bearer",
							Bearer: openapi.HttpAuthBearer{
								Token: "my-token",
							},
						},
					},
				},
			},
		},
		{
			Name: "Should_be_able_to_convert_request_body",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: openapi.HttpRequest{
						Url:     "http://localhost:1234",
						Method:  "POST",
						Headers: []openapi.HttpHeader{},
						Auth:    openapi.HttpAuth{},
						Body:    `{ "message": "hello" }`,
					},
				},
			},
			ExpectedOutput: openapi.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: openapi.TestServiceUnderTest{
					Request: openapi.HttpRequest{
						Url:     "http://localhost:1234",
						Method:  "POST",
						Headers: []openapi.HttpHeader{},
						Body:    `{ "message": "hello" }`,
						Auth:    openapi.HttpAuth{},
					},
				},
			},
		},
		{
			Name: "Should_be_able_to_convert_test_definitions",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: openapi.HttpRequest{
						Url:     "http://localhost:1234",
						Method:  "POST",
						Headers: []openapi.HttpHeader{},
						Auth:    openapi.HttpAuth{},
						Body:    "",
					},
				},
				TestDefinition: []definition.TestDefinition{
					{
						Selector: `span[tracetest.span.type="http"]`,
						Assertions: []string{
							"tracetest.span.duration <= 200",
							"http.status_code = 200",
						},
					},
				},
			},
			ExpectedOutput: openapi.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: openapi.TestServiceUnderTest{
					Request: openapi.HttpRequest{
						Url:     "http://localhost:1234",
						Method:  "POST",
						Headers: []openapi.HttpHeader{},
						Body:    "",
						Auth:    openapi.HttpAuth{},
					},
				},
				Definition: openapi.TestDefinition{
					Definitions: []openapi.TestDefinitionDefinitions{
						{
							Selector: openapi.Selector{
								Query: `span[tracetest.span.type="http"]`,
							},
							Assertions: []openapi.Assertion{
								{
									Attribute:  "tracetest.span.duration",
									Comparator: "<=",
									Expected:   "200",
								},
								{
									Attribute:  "http.status_code",
									Comparator: "=",
									Expected:   "200",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := conversion.ConvertTestDefinitionIntoOpenAPIObject(testCase.Input)
			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedOutput, output)
		})
	}
}
