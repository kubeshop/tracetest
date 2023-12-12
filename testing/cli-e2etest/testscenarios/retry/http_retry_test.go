package testscenarios

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestHTTPRetryConfiguration(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup test server environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	// Make sure version doesn't fail due to version mismatch
	// The mocked server will probably never return the right
	// version number.
	t.Setenv("TRACETEST_DEV", "true")

	fakeHTTPServer, configPath, getNumberServerCalls, err := getFakeVersionServer()
	require.NoError(err)

	defer fakeHTTPServer.Close()

	// Given I am a Tracetest CLI user
	// When I try to check the tracetest version
	// Then I should receive a version string with success
	result := tracetestcli.Exec(t, "version", tracetestcli.WithCLIConfig(configPath))

	helpers.RequireExitCodeEqual(t, result, 0)
	require.Greater(len(result.StdOut), 0)
	require.Greater(getNumberServerCalls(), 1)
}

func getFakeVersionServer() (*httptest.Server, string, func() int, error) {
	callNumber := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if callNumber < 3 {
			callNumber++
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		version := openapi.Version{
			Version: openapi.PtrString("0.14.9"),
			Type:    openapi.PtrString("core"),
		}
		bytes, _ := json.Marshal(version)

		w.Write(bytes)
		w.WriteHeader(http.StatusOK)
	}))

	time.Sleep(time.Second)

	config := config.Config{
		Scheme:   "http",
		Endpoint: strings.ReplaceAll(server.URL, "http://", ""),
	}

	configName := fmt.Sprintf("config_%d.yaml", time.Now().UnixMicro())
	configPath := path.Join(os.TempDir(), configName)

	bytes, err := yaml.Marshal(config)
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not marshal config object: %w", err)
	}

	err = os.WriteFile(configPath, bytes, 0644)
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not save config file: %w", err)
	}

	getNumberInvocations := func() int {
		return callNumber
	}

	return server, configPath, getNumberInvocations, nil
}
