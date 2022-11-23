package yaml_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
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
					Props: map[string]string{
						"username": "user",
						"password": "passwd",
					},
				},
			},
		},
		Specs: (model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}).
			MustAdd(model.SpanQuery(`span[name="Test Span"]`), model.NamedAssertions{
				Name: "count test spans",
				Assertions: []model.Assertion{
					"attr:tracetest.selected_spans.count = 2",
				},
			}),
		Outputs: (model.OrderedMap[string, model.Output]{}).
			MustAdd("user_id", model.Output{
				Selector: model.SpanQuery(`span[name="Create User"]`),
				Value:    "attr:myapp.user_id",
			}),
	}

	expected := `
type: Test
spec:
	id: 123
	name: The Name
	description: Description
	trigger:
		type: http
		httpRequest:
			url: "http://google.com"
			method: POST
			headers:
			- key: Content-Type
				value: application/json
			body: '{"id":123}'
	specs:
	- selector: span[name = "Test Span"]
		assertions:
			- attr:tracetest.selected_spans.count = 2
	outputs:
	- selector: span[name = "Create User"]
		value: attr:myapp.user_id
`

	mapped := yaml.Test(in)
	actual, err := yaml
}
