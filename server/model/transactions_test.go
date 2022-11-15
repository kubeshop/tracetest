package model_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentOutputInjection(t *testing.T) {
	environment := model.Environment{
		ID:   "env-id",
		Name: "my environment",
		Values: []model.EnvironmentValue{
			{Key: "HOST", Value: "http://my-api.com"},
			{Key: "PORT", Value: "8081"},
		},
	}

	run := model.TransactionRun{
		CurrentTest: 1,
		StepRuns: []model.TransactionStepRun{
			{
				Environment: environment,
				Outputs: (model.OrderedMap[string, string]{}).
					MustAdd("TEST_ID", "123456").
					MustAdd("RUN_ID", "2"),
			},
		},
	}

	newEnvironment := run.InjectOutputsIntoEnvironment(environment)

	assert.Contains(t, newEnvironment.Values, model.EnvironmentValue{Key: "HOST", Value: "http://my-api.com"})
	assert.Contains(t, newEnvironment.Values, model.EnvironmentValue{Key: "PORT", Value: "8081"})
	assert.Contains(t, newEnvironment.Values, model.EnvironmentValue{Key: "TEST_ID", Value: "123456"})
	assert.Contains(t, newEnvironment.Values, model.EnvironmentValue{Key: "RUN_ID", Value: "2"})
}
