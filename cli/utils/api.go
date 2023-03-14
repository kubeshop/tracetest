package utils

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
)

func GetAPIClient(cliConfig config.Config) *openapi.APIClient {
	config := openapi.NewConfiguration()
	config.AddDefaultHeader("x-client-id", analytics.ClientID())
	config.Scheme = cliConfig.Scheme
	config.Host = strings.TrimSuffix(cliConfig.Endpoint, "/")
	if cliConfig.ServerPath != nil {
		config.Servers = []openapi.ServerConfiguration{
			{
				URL: *cliConfig.ServerPath,
			},
		}
	}

	return openapi.NewAPIClient(config)
}

type ResourceClient struct {
	Client     http.Client
	BaseUrl    string
	BaseHeader http.Header
}

func GetResourceAPIClient(resourceType string, cliConfig config.Config) ResourceClient {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Second,
		}).DialContext,
	}

	client := http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	baseUrl := fmt.Sprintf("%s://%s/api/%s", cliConfig.Scheme, cliConfig.Endpoint, resourceType)
	baseHeader := http.Header{
		"x-client-id":  []string{analytics.ClientID()},
		"Content-Type": []string{"text/yaml"},
	}

	return ResourceClient{
		Client:     client,
		BaseUrl:    baseUrl,
		BaseHeader: baseHeader,
	}
}

func (resourceClient ResourceClient) NewRequest(url string, method string, body string) (*http.Request, error) {
	var reqBody io.Reader
	if body != "" {
		reqBody = StringToIOReader(body)
	}

	request, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	request.Header = resourceClient.BaseHeader
	return request, err
}
