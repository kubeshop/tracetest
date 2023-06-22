package yaml_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/transaction"
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
						Body: "",
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
						Authentication: &yaml.HTTPAuthentication{
							Type: "basic",
							Basic: &yaml.HTTPBasicAuth{
								Username: "matheus",
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
							Basic: model.BasicAuthenticator{
								Username: "matheus",
								Password: "pikachu",
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
						Authentication: &yaml.HTTPAuthentication{
							Type: "apiKey",
							APIKey: &yaml.HTTPAPIKeyAuth{
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
							APIKey: model.APIKeyAuthenticator{
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
						Authentication: &yaml.HTTPAuthentication{
							Type: "bearer",
							Bearer: &yaml.HTTPBearerAuth{
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
							Bearer: model.BearerAuthenticator{
								Bearer: "my-token",
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
						Body: `{ "message": "hello" }`,
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
						URL:    "http://localhost:1234",
						Method: "POST",
						Body:   "",
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
				Specs: (maps.Ordered[model.SpanQuery, model.NamedAssertions]{}).
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
				Outputs: (maps.Ordered[string, model.Output]{}).
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

func TestTransactionModel(t *testing.T) {
	cases := []struct {
		name     string
		in       yaml.Transaction
		expected transaction.Transaction
	}{
		{
			name: "Basic",
			in: yaml.Transaction{
				ID:          "123",
				Name:        "Transaction",
				Description: "Some transaction",
				Steps:       []string{"345"},
			},
			expected: transaction.Transaction{
				ID:          id.ID("123"),
				Name:        "Transaction",
				Description: "Some transaction",
				StepIDs:     []id.ID{"345"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl := c
			t.Parallel()

			file := yaml.File{
				Type: yaml.FileTypeTransaction,
				Spec: cl.in,
			}

			transaction, err := file.Transaction()
			require.NoError(t, err)

			actual := transaction.Model()

			assert.Equal(t, cl.expected, actual)
		})
	}
}
