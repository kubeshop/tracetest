package config

import (
	"fmt"
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

	config.Scheme = cliConfig.Scheme
	config.Host = strings.TrimSuffix(cliConfig.Endpoint, "/")
	config.Servers = []openapi.ServerConfiguration{
		{
			URL: cliConfig.Path(),
		},
	}

	return openapi.NewAPIClient(config)
}
