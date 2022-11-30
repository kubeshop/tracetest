package yamlconvert_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/kubeshop/tracetest/server/model/yaml/yamlconvert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentConverter(t *testing.T) {
	in := model.Environment{
		ID:          "567",
		Name:        "Env Name",
		Description: "Env Description",
		Values: []model.EnvironmentValue{{
			Key:   "HOST",
			Value: "http://localhost:8080",
		}},
	}

	expected := `type: Environment
spec:
  id: "567"
  name: Env Name
  description: Env Description
  values:
  - key: HOST
    value: http://localhost:8080
`

	mapped := yamlconvert.Environment(in)
	actual, err := yaml.Encode(mapped)
	require.NoError(t, err)

	assert.Equal(t, expected, string(actual))
}
