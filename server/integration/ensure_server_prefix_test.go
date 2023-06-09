package integration_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/kubeshop/tracetest/server/tracedb/datastoreresource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTracetestApp(options ...testmock.TestingAppOption) (*app.App, error) {
	tracetestApp, err := testmock.GetTestingApp(options...)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		tracetestApp.Start()
		time.Sleep(1 * time.Second)
		wg.Done()
	}()

	wg.Wait()

	return tracetestApp, nil
}

func TestServerPrefix(t *testing.T) {
	_, err := getTracetestApp(
		testmock.WithServerPrefix("/tracetest"),
		testmock.WithHttpPort(8000),
	)
	require.NoError(t, err)

	expectedEndpoint := "http://localhost:8000/tracetest"
	tests := getTests(t, expectedEndpoint)
	assert.NotNil(t, tests)

	dataStores := getDatastores(t, expectedEndpoint)
	assert.NotNil(t, dataStores)
	assert.GreaterOrEqual(t, dataStores.Count, 1)
}

func getTests(t *testing.T, endpoint string) resourcemanager.ResourceList[test.Test] {
	url := fmt.Sprintf("%s/api/tests", endpoint)
	resp, err := http.Get(url)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	bodyJsonBytes, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var tests resourcemanager.ResourceList[test.Test]
	err = json.Unmarshal(bodyJsonBytes, &tests)
	require.NoError(t, err)

	return tests
}

func getDatastores(t *testing.T, endpoint string) resourcemanager.ResourceList[datastoreresource.DataStore] {
	url := fmt.Sprintf("%s/api/datastores", endpoint)
	resp, err := http.Get(url)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	bodyJsonBytes, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var dataStores resourcemanager.ResourceList[datastoreresource.DataStore]
	err = yaml.Unmarshal(bodyJsonBytes, &dataStores)
	require.NoError(t, err)

	return dataStores
}
