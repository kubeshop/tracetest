package yamlconvert_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/kubeshop/tracetest/server/model/yaml/yamlconvert"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConverter(t *testing.T) {
	in := model.Test{
		ID:          "123",
		Name:        "The Name",
		Description: "Description",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				Method: model.HTTPMethodPOST,
				URL:    "http://google.com",
				Headers: []model.HTTPHeader{
					{Key: "Content-Type", Value: "application/json"},
				},
				Body: `{"id":123}`,
				Auth: &model.HTTPAuthenticator{
					Type: "basic",
					Basic: model.BasicAuthenticator{
						Username: "user",
						Password: "passwd",
					},
				},
			},
		},
		Specs: (maps.Ordered[model.SpanQuery, model.NamedAssertions]{}).
			MustAdd(model.SpanQuery(`span[name="Test Span"]`), model.NamedAssertions{
				Name: "count test spans",
				Assertions: []model.Assertion{
					"attr:tracetest.selected_spans.count = 2",
				},
			}),
		Outputs: (maps.Ordered[string, model.Output]{}).
			MustAdd("user_id", model.Output{
				Selector: model.SpanQuery(`span[name="Create User"]`),
				Value:    "attr:myapp.user_id",
			}),
	}

	expected := `type: Test
spec:
  id: "123"
  name: The Name
  description: Description
  trigger:
    type: http
    httpRequest:
      url: http://google.com
      method: POST
      headers:
      - key: Content-Type
        value: application/json
      authentication:
        type: basic
        basic:
          username: user
          password: passwd
      body: '{"id":123}'
  specs:
  - name: count test spans
    selector: span[name="Test Span"]
    assertions:
    - attr:tracetest.selected_spans.count = 2
  outputs:
  - name: user_id
    selector: span[name="Create User"]
    value: attr:myapp.user_id
`

	mapped := yamlconvert.Test(in)
	actual, err := yaml.Encode(mapped)
	require.NoError(t, err)

	assert.Equal(t, expected, string(actual))
}
