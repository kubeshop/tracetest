package file_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/definition"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDefinition(t *testing.T) {
	testCases := []struct {
		name          string
		file          string
		expected      definition.File
		expectSuccess bool
	}{
		{
			name:          "Should_parse_valid_definition_file",
			file:          "../testdata/definitions/valid_http_test_definition.yml",
			expectSuccess: true,
			expected: definition.File{
				Type: definition.FileTypeTest,
				Spec: definition.Test{
					Name:        "POST import pokemon",
					Description: "Import a pokemon using its ID",
					Trigger: definition.TestTrigger{
						Type: "http",
						HTTPRequest: definition.HttpRequest{
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
		},
		{
			name:          "Should_parse_valid_definition_file_with_id",
			file:          "../testdata/definitions/valid_http_test_definition_with_id.yml",
			expectSuccess: true,
			expected: definition.File{
				Type: definition.FileTypeTest,
				Spec: definition.Test{
					ID:          "3fd66887-4ee7-44d5-bad8-9934ab9c1a9a",
					Name:        "POST import pokemon",
					Description: "Import a pokemon using its ID",
					Trigger: definition.TestTrigger{
						Type: "http",
						HTTPRequest: definition.HttpRequest{
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
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := file.LoadDefinition(testCase.file)
			if testCase.expectSuccess {
				require.NoError(t, err, "LoadDefinition should not fail")
				err = actual.Validate()

				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, actual)
			} else {
				require.Error(t, err, "LoadDefinition should fail")
			}
		})
	}
}
