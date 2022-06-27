package integration_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testfixtures"
	"github.com/stretchr/testify/assert"
)

const endpointUrl = "http://localhost:8080"

func TestExecutorIntegration(t *testing.T) {
	testRun, err := testfixtures.GetPokemonTestRun(testfixtures.WithCacheDisabled())
	assert.NoError(t, err)

	assert.Equal(t, string(model.RunStateFinished), testRun.State)
	assert.Greater(t, len(testRun.Result.Results), 0)
	assert.True(t, testRun.Result.AllPassed)

	count := 0
	for _, res := range testRun.Result.Results {
		for _, assertionResult := range res.Results {
			for _, spanRes := range assertionResult.SpanResults {
				assert.True(t, spanRes.Passed)
				count = count + 1
			}
		}
	}

	assert.Equal(t, 2, count)
	assert.Equal(t, 10, len(testRun.Trace.Flat))
}
