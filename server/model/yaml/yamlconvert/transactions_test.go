package yamlconvert_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/kubeshop/tracetest/server/model/yaml/yamlconvert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionConverter(t *testing.T) {
	in := model.Transaction{
		ID:          "567",
		Name:        "Transaction Name",
		Description: "Transaction Description",
		Steps: []model.Test{{
			ID:          "1",
			Name:        "Step 1",
			Description: "Step 1 Description",
		}},
	}

	expected := `type: Transaction
spec:
  id: "567"
  name: Transaction Name
  description: Transaction Description
  steps:
  - "1"
`

	mapped := yamlconvert.Transaction(in)
	actual, err := yaml.Encode(mapped)
	require.NoError(t, err)

	assert.Equal(t, expected, string(actual))
}
