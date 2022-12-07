package integration_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerRestartWithNewConfig(t *testing.T) {
	app, err := testfixtures.GetTracetestApp()
	require.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:11633/api/tests", nil)
	response, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	err = app.Stop()
	require.NoError(t, err)

	req, _ = http.NewRequest(http.MethodGet, "http://localhost:11633/api/tests", nil)
	_, err = http.DefaultClient.Do(req)
	assert.Error(t, err)

	config := app.GetConfig()
	config.Server.HttpPort = 12345

	app.SetConfig(config)

	go app.Start()

	time.Sleep(1 * time.Second)

	// should fail because old port is not being used
	req, _ = http.NewRequest(http.MethodGet, "http://localhost:11633/api/tests", nil)
	_, err = http.DefaultClient.Do(req)
	assert.Error(t, err)

	// should pass because it's using new port
	req, _ = http.NewRequest(http.MethodGet, "http://localhost:12345/api/tests", nil)
	response, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
