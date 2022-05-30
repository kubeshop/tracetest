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
		Name               string
		File               string
		ExpectedDefinition definition.Test
		ShouldSucceed      bool
	}{
		{
			Name:          "Should parse valid definition file",
			File:          "../testdata/definitions/valid_http_test_definition.yml",
			ShouldSucceed: true,
			ExpectedDefinition: definition.Test{
				Name: "POST import pokemon",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: definition.HttpRequest{
						URL:    "http://pokemon-demo.tracetest.io/pokemon/import",
						Method: "POST",
						Headers: []definition.HTTPHeader{
							{Key: "ContentType", Value: "application/json"},
						},
						Authentication: definition.HTTPAuthentication{
							Type:  "Bearer",
							Token: "my-token-123",
						},
						Body: definition.HTTPBody{
							Type: "raw",
							Raw:  `{ "id": 52 }`,
						},
					},
				},
			},
		},
		{
			Name:          "Should parse valid definition file with id",
			File:          "../testdata/definitions/valid_http_test_definition_with_id.yml",
			ShouldSucceed: true,
			ExpectedDefinition: definition.Test{
				Id:   "3fd66887-4ee7-44d5-bad8-9934ab9c1a9a",
				Name: "POST import pokemon",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: definition.HttpRequest{
						URL:    "http://pokemon-demo.tracetest.io/pokemon/import",
						Method: "POST",
						Headers: []definition.HTTPHeader{
							{Key: "ContentType", Value: "application/json"},
						},
						Authentication: definition.HTTPAuthentication{
							Type:  "Bearer",
							Token: "my-token-123",
						},
						Body: definition.HTTPBody{
							Type: "raw",
							Raw:  `{ "id": 52 }`,
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
				assert.Equal(t, testCase.ExpectedDefinition, definition)
			} else {
				require.Error(t, err, "LoadDefinition should fail")
			}
		})
	}
}
