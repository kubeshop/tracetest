package yaml_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestModel(t *testing.T) {
	cases := []struct {
		name     string
		in       yaml.Test
		expected model.Test
	}{
		{
			name: "HTTP",
			in: yaml.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: yaml.TestTrigger{
					Type: "http",
					HTTPRequest: yaml.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []yaml.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Authentication: yaml.HTTPAuthentication{},
						Body:           "",
					},
				},
			},
			expected: model.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: model.Trigger{
					Type: "http",
					HTTP: &model.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []model.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: "",
					},
				},
			},
		},
		{
			name: "HTTPBasicAuth",
			in: yaml.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: yaml.TestTrigger{
					Type: "http",
					HTTPRequest: yaml.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []yaml.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: "",
						Authentication: yaml.HTTPAuthentication{
							Type: "basic",
							Basic: yaml.HTTPBasicAuth{
								User:     "matheus",
								Password: "pikachu",
							},
						},
					},
				},
			},
			expected: model.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: model.Trigger{
					Type: "http",
					HTTP: &model.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []model.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: "",
						Auth: &model.HTTPAuthenticator{
							Type: "basic",
							Props: map[string]string{
								"username": "matheus",
								"password": "pikachu",
							},
						},
					},
				},
			},
		},
		{
			name: "HTTPApiKey",
			in: yaml.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: yaml.TestTrigger{
					Type: "http",
					HTTPRequest: yaml.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []yaml.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: "",
						Authentication: yaml.HTTPAuthentication{
							Type: "apiKey",
							ApiKey: yaml.HTTPAPIKeyAuth{
								Key:   "X-Key",
								Value: "my-api-key",
								In:    "header",
							},
						},
					},
				},
			},
			expected: model.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: model.Trigger{
					Type: "http",
					HTTP: &model.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []model.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: "",
						Auth: &model.HTTPAuthenticator{
							Type: "apiKey",
							Props: map[string]string{
								"key":   "X-Key",
								"value": "my-api-key",
								"in":    "header",
							},
						},
					},
				},
			},
		},
		{
			name: "HTTPBearer",
			in: yaml.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: yaml.TestTrigger{
					Type: "http",
					HTTPRequest: yaml.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []yaml.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: "",
						Authentication: yaml.HTTPAuthentication{
							Type: "bearer",
							Bearer: yaml.HTTPBearerAuth{
								Token: "my-token",
							},
						},
					},
				},
			},
			expected: model.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: model.Trigger{
					Type: "http",
					HTTP: &model.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []model.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: "",
						Auth: &model.HTTPAuthenticator{
							Type: "bearer",
							Props: map[string]string{
								"token": "my-token",
							},
						},
					},
				},
			},
		},
		{
			name: "HTTPRequestBody",
			in: yaml.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: yaml.TestTrigger{
					Type: "http",
					HTTPRequest: yaml.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []yaml.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Authentication: yaml.HTTPAuthentication{},
						Body:           `{ "message": "hello" }`,
					},
				},
			},
			expected: model.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: model.Trigger{
					Type: "http",
					HTTP: &model.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []model.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: `{ "message": "hello" }`,
					},
				},
			},
		},
		{
			name: "Definitions",
			in: yaml.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: yaml.TestTrigger{
					Type: "http",
					HTTPRequest: yaml.HTTPRequest{
						URL:            "http://localhost:1234",
						Method:         "POST",
						Headers:        []yaml.HTTPHeader{},
						Authentication: yaml.HTTPAuthentication{},
						Body:           "",
					},
				},
				Specs: []yaml.TestSpec{
					{
						Selector: `span[tracetest.span.type="http"]`,
						Assertions: []string{
							"attr:tracetest.span.duration <= 200ms",
							"attr:http.status_code = 200",
						},
					},
				},
			},
			expected: model.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: model.Trigger{
					Type: "http",
					HTTP: &model.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
					},
				},
				Specs: (model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}).
					MustAdd(model.SpanQuery(`span[tracetest.span.type="http"]`), model.NamedAssertions{
						Name: "",
						Assertions: []model.Assertion{
							`attr:tracetest.span.duration <= 200ms`,
							`attr:http.status_code = 200`,
						},
					}),
			},
		},
		{
			name: "Outputs",
			in: yaml.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: yaml.TestTrigger{
					Type: "http",
					HTTPRequest: yaml.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
					},
				},
				Outputs: yaml.Outputs{
					{
						Name:     "USER_ID",
						Selector: `span[name = "create user"]`,
						Value:    `attr:myapp.user_id`,
					},
				},
			},
			expected: model.Test{
				Name:        "A test",
				Description: "A test description",
				ServiceUnderTest: model.Trigger{
					Type: "http",
					HTTP: &model.HTTPRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
					},
				},
				Outputs: (model.OrderedMap[string, model.Output]{}).
					MustAdd("USER_ID", model.Output{
						Selector: `span[name = "create user"]`,
						Value:    `attr:myapp.user_id`,
					}),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl := c
			t.Parallel()

			file := yaml.File{
				Type: yaml.FileTypeTest,
				Spec: cl.in,
			}

			test, err := file.Test()
			require.NoError(t, err)

			actual := test.Model()

			assert.Equal(t, cl.expected, actual)
		})
	}
}
