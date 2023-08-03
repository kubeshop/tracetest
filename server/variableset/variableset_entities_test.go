package variableset_test

import (
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/variableset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariableSetMerge(t *testing.T) {
	env1 := variableset.VariableSet{
		ID:          "my-env",
		Name:        "my var set",
		Description: "my description",
		CreatedAt:   time.Now().String(),
		Values: []variableset.VariableSetValue{
			{Key: "URL", Value: "http://localhost"},
			{Key: "PORT", Value: "8085"},
		},
	}

	env2 := variableset.VariableSet{
		Values: []variableset.VariableSetValue{
			{Key: "apiKey", Value: "abcdef"},
			{Key: "apiKeyLocation", Value: "header"},
			{Key: "URL", Value: "http://my-api.com"},
		},
	}

	newEnv := env1.Merge(env2)

	assert.Equal(t, env1.ID, newEnv.ID)
	assert.Equal(t, env1.Name, newEnv.Name)
	assert.Equal(t, env1.Description, newEnv.Description)
	assert.Equal(t, env1.CreatedAt, newEnv.CreatedAt)

	require.Len(t, newEnv.Values, 4)
	assert.Contains(t, newEnv.Values, variableset.VariableSetValue{Key: "URL", Value: "http://my-api.com"})
	assert.Contains(t, newEnv.Values, variableset.VariableSetValue{Key: "PORT", Value: "8085"})
	assert.Contains(t, newEnv.Values, variableset.VariableSetValue{Key: "apiKey", Value: "abcdef"})
	assert.Contains(t, newEnv.Values, variableset.VariableSetValue{Key: "apiKeyLocation", Value: "header"})
}

func TestVariableSetMergeWithEmptyVariableSet(t *testing.T) {
	env1 := variableset.VariableSet{
		ID:          "my-env",
		Name:        "my var set",
		Description: "my description",
		CreatedAt:   time.Now().String(),
		Values: []variableset.VariableSetValue{
			{Key: "URL", Value: "http://localhost"},
			{Key: "PORT", Value: "8085"},
		},
	}

	env2 := variableset.VariableSet{
		Values: []variableset.VariableSetValue{},
	}

	newEnv := env1.Merge(env2)
	require.Len(t, newEnv.Values, 2)

}
