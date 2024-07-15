package config

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type ListArgs struct {
	Take          int32
	Skip          int32
	SortDirection string
	SortBy        string
	All           bool
}

func GetAPIClient(cliConfig Config) *openapi.APIClient {
	config := openapi.NewConfiguration()
	config.AddDefaultHeader("x-client-id", analytics.ClientID())
	config.AddDefaultHeader("x-source", "cli")
	config.AddDefaultHeader("x-organization-id", cliConfig.OrganizationID)
	config.AddDefaultHeader("x-environment-id", cliConfig.EnvironmentID)
	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", cliConfig.Jwt))
	if os.Getenv("TRACETEST_DEV_FORCE_URL") == "true" {
		if config.HTTPClient == nil {
			config.HTTPClient = http.DefaultClient
		}
		if config.HTTPClient.Transport == nil {
			config.HTTPClient.Transport = http.DefaultTransport
		}

		if t, ok := config.HTTPClient.Transport.(*http.Transport); ok {
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = &tls.Config{}
			}

			t.TLSClientConfig.InsecureSkipVerify = true
		}
	}

	config.Scheme = cliConfig.Scheme
	config.Host = strings.TrimSuffix(cliConfig.Endpoint, "/")
	config.Servers = []openapi.ServerConfiguration{
		{
			URL: cliConfig.Path(),
		},
	}

	return openapi.NewAPIClient(config)
}
