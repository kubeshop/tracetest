package conversion_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/conversion"
	"github.com/kubeshop/tracetest/cli/definition"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openAPIStr(in string) *string {
	return conversion.ConvertStringIntoOpenAPIString(in)
}

func TestDefinitionToOpenAPIConversion(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          definition.Test
		ExpectedOutput openapi.Test
	}{
		{
			Name: "Should be able to convert request with no authentication and no body",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: definition.HttpRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []definition.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Authentication: definition.HTTPAuthentication{},
						Body:           definition.HTTPBody{},
					},
				},
			},
			ExpectedOutput: openapi.Test{
				Name:        openAPIStr("A test"),
				Description: openAPIStr("A test description"),
				ServiceUnderTest: &openapi.TestServiceUnderTest{
					Request: &openapi.HTTPRequest{
						Url:    openAPIStr("http://localhost:1234"),
						Method: openAPIStr("POST"),
						Headers: []openapi.HTTPHeader{
							{Key: openAPIStr("Content-Type"), Value: openAPIStr("application/json")},
						},
						Body: openAPIStr(""),
						Auth: &openapi.HTTPAuth{},
					},
				},
			},
		},
		{
			Name: "Should be able to convert request with basic auth",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: definition.HttpRequest{
						URL:     "http://localhost:1234",
						Method:  "POST",
						Headers: []definition.HTTPHeader{},
						Body:    definition.HTTPBody{},
						Authentication: definition.HTTPAuthentication{
							BasicAuth: definition.HTTPBasicAuth{
								User:     "matheus",
								Password: "pikachu",
							},
						},
					},
				},
			},
			ExpectedOutput: openapi.Test{
				Name:        openAPIStr("A test"),
				Description: openAPIStr("A test description"),
				ServiceUnderTest: &openapi.TestServiceUnderTest{
					Request: &openapi.HTTPRequest{
						Url:     openAPIStr("http://localhost:1234"),
						Method:  openAPIStr("POST"),
						Headers: []openapi.HTTPHeader{},
						Body:    openAPIStr(""),
						Auth: &openapi.HTTPAuth{
							Basic: &openapi.HTTPAuthBasic{
								Username: openAPIStr("matheus"),
								Password: openAPIStr("pikachu"),
							},
						},
					},
				},
			},
		},
		{
			Name: "Should be able to convert request with api key auth",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: definition.HttpRequest{
						URL:     "http://localhost:1234",
						Method:  "POST",
						Headers: []definition.HTTPHeader{},
						Body:    definition.HTTPBody{},
						Authentication: definition.HTTPAuthentication{
							ApiKey: definition.HTTPAPIKeyAuth{
								Key:   "X-Key",
								Value: "my-api-key",
								In:    "header",
							},
						},
					},
				},
			},
			ExpectedOutput: openapi.Test{
				Name:        openAPIStr("A test"),
				Description: openAPIStr("A test description"),
				ServiceUnderTest: &openapi.TestServiceUnderTest{
					Request: &openapi.HTTPRequest{
						Url:     openAPIStr("http://localhost:1234"),
						Method:  openAPIStr("POST"),
						Headers: []openapi.HTTPHeader{},
						Body:    openAPIStr(""),
						Auth: &openapi.HTTPAuth{
							ApiKey: &openapi.HTTPAuthApiKey{
								Key:   openAPIStr("X-Key"),
								Value: openAPIStr("my-api-key"),
								In:    openAPIStr("header"),
							},
						},
					},
				},
			},
		},
		{
			Name: "Should be able to convert request with bearer auth",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: definition.HttpRequest{
						URL:     "http://localhost:1234",
						Method:  "POST",
						Headers: []definition.HTTPHeader{},
						Body:    definition.HTTPBody{},
						Authentication: definition.HTTPAuthentication{
							Bearer: definition.HTTPBearerAuth{
								Token: "my-token",
							},
						},
					},
				},
			},
			ExpectedOutput: openapi.Test{
				Name:        openAPIStr("A test"),
				Description: openAPIStr("A test description"),
				ServiceUnderTest: &openapi.TestServiceUnderTest{
					Request: &openapi.HTTPRequest{
						Url:     openAPIStr("http://localhost:1234"),
						Method:  openAPIStr("POST"),
						Headers: []openapi.HTTPHeader{},
						Body:    openAPIStr(""),
						Auth: &openapi.HTTPAuth{
							Bearer: &openapi.HTTPAuthBearer{
								Token: openAPIStr("my-token"),
							},
						},
					},
				},
			},
		},
		{
			Name: "Should be able to convert request body",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: definition.HttpRequest{
						URL:            "http://localhost:1234",
						Method:         "POST",
						Headers:        []definition.HTTPHeader{},
						Authentication: definition.HTTPAuthentication{},
						Body: definition.HTTPBody{
							Type: "raw",
							Raw:  `{ "message": "hello" }`,
						},
					},
				},
			},
			ExpectedOutput: openapi.Test{
				Name:        openAPIStr("A test"),
				Description: openAPIStr("A test description"),
				ServiceUnderTest: &openapi.TestServiceUnderTest{
					Request: &openapi.HTTPRequest{
						Url:     openAPIStr("http://localhost:1234"),
						Method:  openAPIStr("POST"),
						Headers: []openapi.HTTPHeader{},
						Body:    openAPIStr(`{ "message": "hello" }`),
						Auth:    &openapi.HTTPAuth{},
					},
				},
			},
		},
		{
			Name: "Should be able to convert test definitions",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: definition.HttpRequest{
						URL:            "http://localhost:1234",
						Method:         "POST",
						Headers:        []definition.HTTPHeader{},
						Authentication: definition.HTTPAuthentication{},
						Body:           definition.HTTPBody{},
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
				Name:        openAPIStr("A test"),
				Description: openAPIStr("A test description"),
				ServiceUnderTest: &openapi.TestServiceUnderTest{
					Request: &openapi.HTTPRequest{
						Url:     openAPIStr("http://localhost:1234"),
						Method:  openAPIStr("POST"),
						Headers: []openapi.HTTPHeader{},
						Body:    openAPIStr(""),
						Auth:    &openapi.HTTPAuth{},
					},
				},
				Definition: &openapi.TestDefinition{
					Definitions: []openapi.TestDefinitionDefinitions{
						{
							Selector: openAPIStr(`span[tracetest.span.type="http"]`),
							Assertions: []openapi.Assertion{
								{
									Attribute:  openAPIStr("tracetest.span.duration"),
									Comparator: openAPIStr("<="),
									Expected:   openAPIStr("100"),
								},
								{
									Attribute:  openAPIStr("http.status_code"),
									Comparator: openAPIStr("="),
									Expected:   openAPIStr("200"),
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
