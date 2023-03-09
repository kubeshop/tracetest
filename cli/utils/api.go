package utils

import (
	"bytes"
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

func GetHttpClient(cliConfig config.Config) http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Second,
		}).DialContext,
	}

	client := http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return client
}

func GetResourceRequest(resourceType string, cliConfig config.Config, method string, body string) (*http.Request, error) {
	urlString := fmt.Sprintf("%s://%s/api/%s/", cliConfig.Scheme, strings.TrimSuffix(cliConfig.Endpoint, "/"), resourceType)

	header := http.Header{
		"X-Client-Id":  []string{analytics.ClientID()},
		"content-type": []string{"text/yaml"},
	}

	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	}

	request, err := http.NewRequest(method, urlString, reqBody)

	if err != nil {
		return &http.Request{}, err
	}
	request.Header = header

	return request, nil
}
