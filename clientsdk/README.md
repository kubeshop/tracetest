# Client SDK

This folder contains Client SDKs for Tracetest API, considering our definitions on `api` folder for OpenAPI, and considering our Resource Manager API. 

You can refer to our OpenAPI client with any Golang code by:

```go
import "github.com/kubeshop/tracetest/clientdsdk/openapi"

// somewhere in the code
openapiClient := openapi.NewAPIClient(&openapi.Configuration{
		Scheme:        "http",
		Host:          "address.to.tracetest",
		DefaultHeader: map[string]string{},
		Servers: openapi.ServerConfigurations{
			{URL: "/api"},
		},
	})

// usage example
testRun, resp, err := openapiClient.ApiApi.GetTestRun(ctx, testID, runID).Execute()
if err != nil {
  return nil, fmt.Errorf("could not send request to get test run from tracetest: %w", err)
}

if resp.StatusCode != http.StatusOK {
  return nil, fmt.Errorf("get test run endpoint returned an unexpected status code. Expected 200, got %d", resp.StatusCode)
}

// ...
```

And you can refer to our ResourceManager client with any Golang code by:

```go
import (
	"github.com/kubeshop/tracetest/clientdsdk/openapi"
	"github.com/kubeshop/tracetest/clientdsdk/resourcemanager"
)

// somewhere in the code
httpClient := resourcemanager.NewHTTPClient("address.to.tracetest", nil)
dataStoresClient := resourcemanager.NewClient(
	httpClient, log.Desugar(),
	"datastore", "datastores",
	resourcemanager.WithResourceType("DataStore"),
)

// usage example
dataStoreJSON, err := c.dataStores.Get(ctx, datastoreID, resourcemanager.Formats.Get(resourcemanager.FormatJSON))
if err != nil {
	return datastore.DataStore{}, fmt.Errorf("could not get DataStore: %w", err)
}

var dataStoreResource openapi.DataStoreResource
err = json.Unmarshal([]byte(dataStoreJSON), &dataStoreResource)
if err != nil {
	return datastore.DataStore{}, fmt.Errorf("could not unmarshal DataStore: %w", err)
}

// ...
```
