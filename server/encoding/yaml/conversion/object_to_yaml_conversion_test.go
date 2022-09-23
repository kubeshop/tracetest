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
				URL:    "http://localhost:11633",
				Method: "GET",
			},
		},
		Specs: []definition.TestSpec{
			{
				Selector: `span[name="My span"]`,
				Assertions: []string{
					"assertion 1",
					"assertion 2",
				},
			},
			{
				Selector: `span[name="My other span"]`,
				Assertions: []string{
					"assertion 3",
					"assertion 4",
				},
			},
		},
	}

	expectedYaml := `id: abcdeef
name: This is a value
trigger:
  type: http
  httpRequest:
    url: http://localhost:11633
    method: GET
specs:
- selector: span[name="My span"]
  assertions:
  - assertion 1
  - assertion 2
- selector: span[name="My other span"]
  assertions:
  - assertion 3
  - assertion 4
`

	yamlBytes, err := conversion.GetYamlFileFromDefinition(def)
	require.NoError(t, err)

	yamlContent := string(yamlBytes)

	assert.Equal(t, expectedYaml, yamlContent)
}
