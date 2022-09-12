package conversion_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/encoding/yaml/conversion"
	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjectToYamlConversion(t *testing.T) {
	def := definition.Test{
		Id:          "abcdeef",
		Name:        "This is a value",
		Description: "", // should be ommited
		Trigger: definition.TestTrigger{
			Type: "http",
			HTTPRequest: definition.HTTPRequest{
				URL:    "http://localhost:8080",
				Method: "GET",
			},
		},
	}

	expectedYaml := `id: abcdeef
name: This is a value
trigger:
  httpRequest:
    method: GET
    url: http://localhost:8080
  type: http
`

	yamlBytes, err := conversion.GetYamlFileFromDefinition(def)
	require.NoError(t, err)

	yamlContent := string(yamlBytes)

	assert.Equal(t, expectedYaml, yamlContent)
}
