package conversion_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/conversion"
	"github.com/kubeshop/tracetest/cli/definition"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openApiInt(in int32) *int32 {
	return &in
}

func TestOpenAPIToDefinitionConversion(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          openapi.Test
		ExpectedOutput definition.Test
	}{
		{
			Name: "Should_convert_basic_test_information",
			Input: openapi.Test{
				Id:               openAPIStr("624a8dea-f152-48d4-a742-30b210094959"),
				Name:             openAPIStr("my test"),
				Description:      openAPIStr("my test description"),
				Version:          openApiInt(3),
				ServiceUnderTest: &openapi.Trigger{},
				Definition:       &openapi.TestDefinition{},
			},
			ExpectedOutput: definition.Test{
				Id:             "624a8dea-f152-48d4-a742-30b210094959",
				Name:           "my test",
				Description:    "my test description",
				Trigger:        definition.TestTrigger{},
				TestDefinition: []definition.TestDefinition{},
			},
		},
		{
			Name: "should_convert_service_under_test_information_with_no_authentication",
			Input: openapi.Test{
				Id:          openAPIStr("624a8dea-f152-48d4-a742-30b210094959"),
				Name:        openAPIStr("my test"),
				Description: openAPIStr("my test description"),
				Version:     openApiInt(3),
				ServiceUnderTest: &openapi.Trigger{
					TriggerType: openAPIStr("http"),
					TriggerSettings: &openapi.TriggerTriggerSettings{
						Http: &openapi.HTTPRequest{
							Url:    openAPIStr("http://localhost:1234"),
							Method: openAPIStr("POST"),
							Headers: []openapi.HTTPHeader{
								{Key: openAPIStr("Content-Type"), Value: openAPIStr("application/json")},
							},
							Body: openAPIStr(`{ "id": 52 }`),
							Auth: &openapi.HTTPAuth{},
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
					HTTPRequest: definition.HttpRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []definition.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Authentication: definition.HTTPAuthentication{},
						Body:           `{ "id": 52 }`,
					},
				},
				TestDefinition: []definition.TestDefinition{},
			},
		},
		{
			Name: "should_convert_service_under_test_information_with_no_body",
			Input: openapi.Test{
				Id:          openAPIStr("624a8dea-f152-48d4-a742-30b210094959"),
				Name:        openAPIStr("my test"),
				Description: openAPIStr("my test description"),
				Version:     openApiInt(3),
				ServiceUnderTest: &openapi.Trigger{
					TriggerType: openAPIStr("http"),
					TriggerSettings: &openapi.TriggerTriggerSettings{
						Http: &openapi.HTTPRequest{
							Url:    openAPIStr("http://localhost:1234"),
							Method: openAPIStr("POST"),
							Headers: []openapi.HTTPHeader{
								{Key: openAPIStr("Content-Type"), Value: openAPIStr("application/json")},
							},
							Body: nil,
							Auth: &openapi.HTTPAuth{},
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
					HTTPRequest: definition.HttpRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []definition.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Authentication: definition.HTTPAuthentication{},
					},
				},
				TestDefinition: []definition.TestDefinition{},
			},
		},
		{
			Name: "should_convert_service_under_test_information_with_basic_authentication",
			Input: openapi.Test{
				Id:          openAPIStr("624a8dea-f152-48d4-a742-30b210094959"),
				Name:        openAPIStr("my test"),
				Description: openAPIStr("my test description"),
				Version:     openApiInt(3),
				ServiceUnderTest: &openapi.Trigger{
					TriggerType: openAPIStr("http"),
					TriggerSettings: &openapi.TriggerTriggerSettings{
						Http: &openapi.HTTPRequest{
							Url:    openAPIStr("http://localhost:1234"),
							Method: openAPIStr("POST"),
							Headers: []openapi.HTTPHeader{
								{Key: openAPIStr("Content-Type"), Value: openAPIStr("application/json")},
							},
							Body: openAPIStr(`{ "id": 52 }`),
							Auth: &openapi.HTTPAuth{
								Type: openAPIStr("basic"),
								Basic: &openapi.HTTPAuthBasic{
									Username: openAPIStr("my username"),
									Password: openAPIStr("my password"),
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
					HTTPRequest: definition.HttpRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []definition.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Authentication: definition.HTTPAuthentication{
							Type: "basic",
							Basic: definition.HTTPBasicAuth{
								User:     "my username",
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
				Id:          openAPIStr("624a8dea-f152-48d4-a742-30b210094959"),
				Name:        openAPIStr("my test"),
				Description: openAPIStr("my test description"),
				Version:     openApiInt(3),
				ServiceUnderTest: &openapi.Trigger{
					TriggerType: openAPIStr("http"),
					TriggerSettings: &openapi.TriggerTriggerSettings{
						Http: &openapi.HTTPRequest{
							Url:    openAPIStr("http://localhost:1234"),
							Method: openAPIStr("POST"),
							Headers: []openapi.HTTPHeader{
								{Key: openAPIStr("Content-Type"), Value: openAPIStr("application/json")},
							},
							Body: openAPIStr(`{ "id": 52 }`),
							Auth: &openapi.HTTPAuth{
								Type: openAPIStr("apikey"),
								ApiKey: &openapi.HTTPAuthApiKey{
									Key:   openAPIStr("X-Key"),
									Value: openAPIStr("my-key"),
									In:    openAPIStr("header"),
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
					HTTPRequest: definition.HttpRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []definition.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Authentication: definition.HTTPAuthentication{
							Type: "apikey",
							ApiKey: definition.HTTPAPIKeyAuth{
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
				Id:          openAPIStr("624a8dea-f152-48d4-a742-30b210094959"),
				Name:        openAPIStr("my test"),
				Description: openAPIStr("my test description"),
				Version:     openApiInt(3),
				ServiceUnderTest: &openapi.Trigger{
					TriggerType: openAPIStr("http"),
					TriggerSettings: &openapi.TriggerTriggerSettings{
						Http: &openapi.HTTPRequest{
							Url:    openAPIStr("http://localhost:1234"),
							Method: openAPIStr("POST"),
							Headers: []openapi.HTTPHeader{
								{Key: openAPIStr("Content-Type"), Value: openAPIStr("application/json")},
							},
							Body: openAPIStr(`{ "id": 52 }`),
							Auth: &openapi.HTTPAuth{
								Type: openAPIStr("bearer"),
								Bearer: &openapi.HTTPAuthBearer{
									Token: openAPIStr("my token"),
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
					HTTPRequest: definition.HttpRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []definition.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Authentication: definition.HTTPAuthentication{
							Type: "bearer",
							Bearer: definition.HTTPBearerAuth{
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
				Id:               openAPIStr("624a8dea-f152-48d4-a742-30b210094959"),
				Name:             openAPIStr("my test"),
				Description:      openAPIStr("my test description"),
				Version:          openApiInt(3),
				ServiceUnderTest: &openapi.Trigger{},
				Definition: &openapi.TestDefinition{
					Definitions: []openapi.TestDefinitionDefinitions{
						{
							Selector: &openapi.Selector{
								Query: openAPIStr(`span[name = "my span name"]`),
							},
							Assertions: []openapi.Assertion{
								{
									Attribute:  openAPIStr("tracetest.span.duration"),
									Comparator: openAPIStr("<="),
									Expected:   openAPIStr("200"),
								},
								{
									Attribute:  openAPIStr("db.operation"),
									Comparator: openAPIStr("="),
									Expected:   openAPIStr("create"),
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
				Trigger:     definition.TestTrigger{},
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
