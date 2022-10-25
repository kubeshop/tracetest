package file_test

import (
	"os"
	"testing"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDefinition(t *testing.T) {
	testCases := []struct {
		name          string
		file          string
		expected      yaml.File
		expectSuccess bool
		envVariables  map[string]string
	}{
		{
			name:          "Should_parse_valid_definition_file",
			file:          "../testdata/definitions/valid_http_test_definition.yml",
			expectSuccess: true,
			expected: yaml.File{
				Type: yaml.FileTypeTest,
				Spec: yaml.Test{
					Name:        "POST import pokemon",
					Description: "Import a pokemon using its ID",
					Trigger: yaml.TestTrigger{
						Type: "http",
						HTTPRequest: yaml.HTTPRequest{
							URL:    "http://pokemon-demo.tracetest.io/pokemon/import",
							Method: "POST",
							Headers: []yaml.HTTPHeader{
								{Key: "Content-Type", Value: "application/json"},
							},
							Body: `{ "id": 52 }`,
						},
					},
					Specs: []yaml.TestSpec{
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
			expected: yaml.File{
				Type: yaml.FileTypeTest,
				Spec: yaml.Test{
					ID:          "3fd66887-4ee7-44d5-bad8-9934ab9c1a9a",
					Name:        "POST import pokemon",
					Description: "Import a pokemon using its ID",
					Trigger: yaml.TestTrigger{
						Type: "http",
						HTTPRequest: yaml.HTTPRequest{
							URL:    "http://pokemon-demo.tracetest.io/pokemon/import",
							Method: "POST",
							Headers: []yaml.HTTPHeader{
								{Key: "Content-Type", Value: "application/json"},
							},
							Body: `{ "id": 52 }`,
						},
					},
					Specs: []yaml.TestSpec{
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
		{
			name:          "Should_parse_valid_definition_file_with_id",
			file:          "../testdata/definitions/valid_http_test_definition_with_env_variables.yml",
			expectSuccess: true,
			envVariables: map[string]string{
				"POKEMON_APP_API_KEY": "1234",
			},
			expected: yaml.File{
				Type: yaml.FileTypeTest,
				Spec: yaml.Test{
					Name:        "POST import pokemon",
					Description: "Import a pokemon using its ID",
					Trigger: yaml.TestTrigger{
						Type: "http",
						HTTPRequest: yaml.HTTPRequest{
							URL:    "http://pokemon-demo.tracetest.io/pokemon/import",
							Method: "POST",
							Headers: []yaml.HTTPHeader{
								{Key: "Content-Type", Value: "application/json"},
							},
							Body: `{ "id": 52 }`,
							Authentication: yaml.HTTPAuthentication{
								Type: "apiKey",
								ApiKey: yaml.HTTPAPIKeyAuth{
									Key:   "X-Key",
									Value: "1234",
									In:    "header",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			for envName, envValue := range testCase.envVariables {
				os.Setenv(envName, envValue)
			}

			defer func() {
				for envName, _ := range testCase.envVariables {
					os.Unsetenv(envName)
				}
			}()

			actual, err := file.Read(testCase.file)
			if testCase.expectSuccess {
				require.NoError(t, err)
				err = actual.Definition().Validate()
				assert.NoError(t, err)
				resolvedFile, err := actual.ResolveVariables()
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, resolvedFile.Definition())
			} else {
				require.Error(t, err, "LoadDefinition should fail")
			}
		})
	}
}

func TestSetID(t *testing.T) {
	t.Run("NoID", func(t *testing.T) {
		f, err := file.Read("../testdata/definitions/valid_http_test_definition.yml")
		require.NoError(t, err)

		test, _ := f.Definition().Test()
		assert.Equal(t, "", test.ID)

		f, err = f.SetID("new-id")
		require.NoError(t, err)

		test, _ = f.Definition().Test()
		assert.Equal(t, "new-id", test.ID)
	})

	t.Run("WithID", func(t *testing.T) {
		f, err := file.Read("../testdata/definitions/valid_http_test_definition_with_id.yml")
		require.NoError(t, err)

		_, err = f.SetID("new-id")
		require.ErrorIs(t, err, file.ErrFileHasID)
	})

}
