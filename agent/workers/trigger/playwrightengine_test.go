package trigger_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubeshop/tracetest/agent/workers/trigger"
	triggerer "github.com/kubeshop/tracetest/agent/workers/trigger"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/stretchr/testify/assert"
)

func createOptions() *trigger.Options {
	return &trigger.Options{
		TraceID: id.NewRandGenerator().TraceID(),
		SpanID:  id.NewRandGenerator().SpanID(),
	}
}

var script = `
	const { expect } = require('@playwright/test');
	async function basicTest(page) {
		await expect(page.getByText('OK')).toBeVisible();
	}

	module.exports = { basicTest };
`

func TestTrigger(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, "GET", req.Method)

		rw.WriteHeader(200)
		_, err := rw.Write([]byte(`OK`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	triggerConfig := trigger.Trigger{
		Type: trigger.TriggerTypeHTTP,
		PlaywrightEngine: &trigger.PlaywrightEngineRequest{
			Target: server.URL,
			Method: "basicTest",
			Script: script,
		},
	}

	ex := triggerer.PLAYWRIGHTENGINE()

	resp, err := ex.Trigger(createContext(), triggerConfig, createOptions())
	assert.NoError(t, err)

	assert.Equal(t, true, resp.Result.PlaywrightEngine.Success)
}

var scriptFail = `
	const { expect } = require('@playwright/test');
	async function basicTest(page) {
		await expect(page.getByText('NOT FOUND YORCH')).toBeVisible();
	}

	module.exports = { basicTest };
`

func TestTriggerFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, "GET", req.Method)

		rw.WriteHeader(200)
		_, err := rw.Write([]byte(`OK`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	triggerConfig := trigger.Trigger{
		Type: trigger.TriggerTypeHTTP,
		PlaywrightEngine: &trigger.PlaywrightEngineRequest{
			Target: server.URL,
			Method: "basicTest",
			Script: scriptFail,
		},
	}

	ex := triggerer.PLAYWRIGHTENGINE()

	resp, err := ex.Trigger(createContext(), triggerConfig, createOptions())
	assert.NotNil(t, err)

	assert.Equal(t, false, resp.Result.PlaywrightEngine.Success)
}
