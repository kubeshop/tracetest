package file_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDefinition(t *testing.T) {
	testCases := []struct {
		Name               string
		File               string
		ExpectedDefinition definition.Test
		ShouldSucceed      bool
	}{
		{
			Name:          "Should_parse_valid_definition_file",
			File:          "../testdata/definitions/valid_http_test_definition.yml",
			ShouldSucceed: true,
			ExpectedDefinition: definition.Test{
				Name:        "POST import pokemon",
				Description: "Import a pokemon using its ID",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: definition.HTTPRequest{
						URL:    "http://pokemon-demo.tracetest.io/pokemon/import",
						Method: "POST",
						Headers: []definition.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: `{ "id": 52 }`,
					},
				},
				Specs: []definition.TestSpec{
					{
						Selector: "span[name = \"POST /pokemon/import\"]",
						Assertions: []string{
							"tracetest.span.duration <= 100ms",
							"http.status_code = 200",
						},
					},
					{
						Selector: "span[name = \"send message to queue\"]",
						Assertions: []string{
							"messaging.message.payload contains 52",
						},
					},
					{
						Selector: "span[name = \"consume message from queue\"]:last",
						Assertions: []string{
							"messaging.message.payload contains 52",
						},
					},
					{
						Selector: "span[name = \"consume message from queue\"]:last span[name = \"import pokemon from pokeapi\"]",
						Assertions: []string{
							"http.status_code = 200",
						},
					},
					{
						Selector: "span[name = \"consume message from queue\"]:last span[name = \"save pokemon on database\"]",
						Assertions: []string{
							"db.repository.operation = \"create\"",
							"tracetest.span.duration <= 100ms",
							`tracetest.response.body contains "\"id\": 52"`,
						},
					},
				},
			},
		},
		{
			Name:          "Should_parse_valid_definition_file_with_id",
			File:          "../testdata/definitions/valid_http_test_definition_with_id.yml",
			ShouldSucceed: true,
			ExpectedDefinition: definition.Test{
				Id:          "3fd66887-4ee7-44d5-bad8-9934ab9c1a9a",
				Name:        "POST import pokemon",
				Description: "Import a pokemon using its ID",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: definition.HTTPRequest{
						URL:    "http://pokemon-demo.tracetest.io/pokemon/import",
						Method: "POST",
						Headers: []definition.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Body: `{ "id": 52 }`,
					},
				},
				Specs: []definition.TestSpec{
					{
						Selector: "span[name = \"POST /pokemon/import\"]",
						Assertions: []string{
							"tracetest.span.duration <= 100ms",
							"http.status_code = 200",
						},
					},
					{
						Selector: "span[name = \"send message to queue\"]",
						Assertions: []string{
							"messaging.message.payload contains 52",
						},
					},
					{
						Selector: "span[name = \"consume message from queue\"]:last",
						Assertions: []string{
							"messaging.message.payload contains 52",
						},
					},
					{
						Selector: "span[name = \"consume message from queue\"]:last span[name = \"import pokemon from pokeapi\"]",
						Assertions: []string{
							"http.status_code = 200",
						},
					},
					{
						Selector: "span[name = \"consume message from queue\"]:last span[name = \"save pokemon on database\"]",
						Assertions: []string{
							"db.repository.operation = \"create\"",
							"tracetest.span.duration <= 100ms",
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			definition, err := file.LoadDefinition(testCase.File)
			if testCase.ShouldSucceed {
				require.NoError(t, err, "LoadDefinition should not fail")
				err = definition.Validate()
				assert.NoError(t, err)
				assert.Equal(t, testCase.ExpectedDefinition, definition)
			} else {
				require.Error(t, err, "LoadDefinition should fail")
			}
		})
	}
}
