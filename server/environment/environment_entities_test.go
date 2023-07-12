package environment_test

import (
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentMerge(t *testing.T) {
	env1 := environment.Environment{
		ID:          "my-env",
		Name:        "my environment",
		Description: "my description",
		CreatedAt:   time.Now().String(),
		Values: []environment.EnvironmentValue{
			{Key: "URL", Value: "http://localhost"},
			{Key: "PORT", Value: "8085"},
		},
	}

	env2 := environment.Environment{
		Values: []environment.EnvironmentValue{
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
	assert.Contains(t, newEnv.Values, environment.EnvironmentValue{Key: "URL", Value: "http://my-api.com"})
	assert.Contains(t, newEnv.Values, environment.EnvironmentValue{Key: "PORT", Value: "8085"})
	assert.Contains(t, newEnv.Values, environment.EnvironmentValue{Key: "apiKey", Value: "abcdef"})
	assert.Contains(t, newEnv.Values, environment.EnvironmentValue{Key: "apiKeyLocation", Value: "header"})
}

func TestEnvironmentMergeWithEmptyEnvironment(t *testing.T) {
	env1 := environment.Environment{
		ID:          "my-env",
		Name:        "my environment",
		Description: "my description",
		CreatedAt:   time.Now().String(),
		Values: []environment.EnvironmentValue{
			{Key: "URL", Value: "http://localhost"},
			{Key: "PORT", Value: "8085"},
		},
	}

	env2 := environment.Environment{
		Values: []environment.EnvironmentValue{},
	}

	newEnv := env1.Merge(env2)
	require.Len(t, newEnv.Values, 2)

}
