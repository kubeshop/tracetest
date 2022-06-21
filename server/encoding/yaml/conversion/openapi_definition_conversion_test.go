package conversion_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/encoding/yaml/conversion"
	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"github.com/kubeshop/tracetest/server/openapi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAPIToDefinitionConversion(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          openapi.Test
		ExpectedOutput definition.Test
	}{
		{
			Name: "Should_convert_basic_test_information",
			Input: openapi.Test{
				Id:               "624a8dea-f152-48d4-a742-30b210094959",
				Name:             "my test",
				Description:      "my test description",
				Version:          3,
				ServiceUnderTest: openapi.TestServiceUnderTest{},
				Definition:       openapi.TestDefinition{},
			},
			ExpectedOutput: definition.Test{
				Id:          "624a8dea-f152-48d4-a742-30b210094959",
				Name:        "my test",
				Description: "my test description",
				Trigger: definition.TestTrigger{
					Type: "http",
				},
				TestDefinition: []definition.TestDefinition{},
			},
		},
		{
			Name: "should_convert_service_under_test_information_with_no_authentication",
			Input: openapi.Test{
				Id:          "624a8dea-f152-48d4-a742-30b210094959",
				Name:        "my test",
				Description: "my test description",
				Version:     3,
				ServiceUnderTest: openapi.TestServiceUnderTest{
					Request: openapi.HttpRequest{
						Url:    "http://localhost:1234",
						Method: "POST",
						Headers: []openapi.HttpHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: `{ "id": 52 }`,
						Auth: openapi.HttpAuth{},
					},
				},
			},
			ExpectedOutput: definition.Test{
				Id:          "624a8dea-f152-48d4-a742-30b210094959",
				Name:        "my test",
				Description: "my test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: openapi.HttpRequest{
						Url:    "http://localhost:1234",
						Method: "POST",
						Headers: []openapi.HttpHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Auth: openapi.HttpAuth{},
						Body: `{ "id": 52 }`,
					},
				},
				TestDefinition: []definition.TestDefinition{},
			},
		},
		{
			Name: "should_convert_service_under_test_information_with_no_body",
			Input: openapi.Test{
				Id:          "624a8dea-f152-48d4-a742-30b210094959",
				Name:        "my test",
				Description: "my test description",
				Version:     3,
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
			ExpectedOutput: definition.Test{
				Id:          "624a8dea-f152-48d4-a742-30b210094959",
				Name:        "my test",
				Description: "my test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: openapi.HttpRequest{
						Url:    "http://localhost:1234",
						Method: "POST",
						Headers: []openapi.HttpHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Auth: openapi.HttpAuth{},
					},
				},
				TestDefinition: []definition.TestDefinition{},
			},
		},
		{
			Name: "should_convert_service_under_test_information_with_basic_authentication",
			Input: openapi.Test{
				Id:          "624a8dea-f152-48d4-a742-30b210094959",
				Name:        "my test",
				Description: "my test description",
				Version:     3,
				ServiceUnderTest: openapi.TestServiceUnderTest{
					Request: openapi.HttpRequest{
						Url:    "http://localhost:1234",
						Method: "POST",
						Headers: []openapi.HttpHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: `{ "id": 52 }`,
						Auth: openapi.HttpAuth{
							Type: "basic",
							Basic: openapi.HttpAuthBasic{
								Username: "my username",
								Password: "my password",
							},
						},
					},
				},
			},
			ExpectedOutput: definition.Test{
				Id:          "624a8dea-f152-48d4-a742-30b210094959",
				Name:        "my test",
				Description: "my test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: openapi.HttpRequest{
						Url:    "http://localhost:1234",
						Method: "POST",
						Headers: []openapi.HttpHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Auth: openapi.HttpAuth{
							Type: "basic",
							Basic: openapi.HttpAuthBasic{
								Username: "my username",
								Password: "my password",
							},
						},
						Body: `{ "id": 52 }`,
					},
				},
				TestDefinition: []definition.TestDefinition{},
			},
		},
		{
			Name: "should_convert_service_under_test_information_with_apikey_authentication",
			Input: openapi.Test{
				Id:          "624a8dea-f152-48d4-a742-30b210094959",
				Name:        "my test",
				Description: "my test description",
				Version:     3,
				ServiceUnderTest: openapi.TestServiceUnderTest{
					Request: openapi.HttpRequest{
						Url:    "http://localhost:1234",
						Method: "POST",
						Headers: []openapi.HttpHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: `{ "id": 52 }`,
						Auth: openapi.HttpAuth{
							Type: "apiKey",
							ApiKey: openapi.HttpAuthApiKey{
								Key:   "X-Key",
								Value: "my-key",
								In:    "header",
							},
						},
					},
				},
			},
			ExpectedOutput: definition.Test{
				Id:          "624a8dea-f152-48d4-a742-30b210094959",
				Name:        "my test",
				Description: "my test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: openapi.HttpRequest{
						Url:    "http://localhost:1234",
						Method: "POST",
						Headers: []openapi.HttpHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Auth: openapi.HttpAuth{
							Type: "apiKey",
							ApiKey: openapi.HttpAuthApiKey{
								Key:   "X-Key",
								Value: "my-key",
								In:    "header",
							},
						},
						Body: `{ "id": 52 }`,
					},
				},
				TestDefinition: []definition.TestDefinition{},
			},
		},
		{
			Name: "should_convert_service_under_test_information_with_bearer_authentication",
			Input: openapi.Test{
				Id:          "624a8dea-f152-48d4-a742-30b210094959",
				Name:        "my test",
				Description: "my test description",
				Version:     3,
				ServiceUnderTest: openapi.TestServiceUnderTest{
					Request: openapi.HttpRequest{
						Url:    "http://localhost:1234",
						Method: "POST",
						Headers: []openapi.HttpHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: `{ "id": 52 }`,
						Auth: openapi.HttpAuth{
							Type: "bearer",
							Bearer: openapi.HttpAuthBearer{
								Token: "my token",
							},
						},
					},
				},
			},
			ExpectedOutput: definition.Test{
				Id:          "624a8dea-f152-48d4-a742-30b210094959",
				Name:        "my test",
				Description: "my test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: openapi.HttpRequest{
						Url:    "http://localhost:1234",
						Method: "POST",
						Headers: []openapi.HttpHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Auth: openapi.HttpAuth{
							Type: "bearer",
							Bearer: openapi.HttpAuthBearer{
								Token: "my token",
							},
						},
						Body: `{ "id": 52 }`,
					},
				},
				TestDefinition: []definition.TestDefinition{},
			},
		},
		{
			Name: "Should_convert_test_definition",
			Input: openapi.Test{
				Id:               "624a8dea-f152-48d4-a742-30b210094959",
				Name:             "my test",
				Description:      "my test description",
				Version:          3,
				ServiceUnderTest: openapi.TestServiceUnderTest{},
				Definition: openapi.TestDefinition{
					Definitions: []openapi.TestDefinitionDefinitions{
						{
							Selector: openapi.Selector{
								Query: `span[name = "my span name"]`,
							},
							Assertions: []openapi.Assertion{
								{
									Attribute:  "tracetest.span.duration",
									Comparator: "<=",
									Expected:   "200",
								},
								{
									Attribute:  "db.operation",
									Comparator: "=",
									Expected:   "create",
								},
							},
						},
					},
				},
			},
			ExpectedOutput: definition.Test{
				Id:          "624a8dea-f152-48d4-a742-30b210094959",
				Name:        "my test",
				Description: "my test description",
				Trigger: definition.TestTrigger{
					Type: "http",
				},
				TestDefinition: []definition.TestDefinition{
					{
						Selector: `span[name = "my span name"]`,
						Assertions: []string{
							"tracetest.span.duration <= 200",
							`db.operation = "create"`,
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := conversion.ConvertOpenAPITestIntoDefinitionObject(testCase.Input)

			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedOutput, output)
		})
	}
}
