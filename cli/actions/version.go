package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
)

func GetVersion(ctx context.Context, cfg config.Config, client *openapi.APIClient) (string, bool) {
	result := fmt.Sprintf(`CLI: %s`, config.Version)

	if cfg.IsEmpty() {
		return result + `
Server: Not Configured`, false
	}

	version, err := getServerVersion(ctx, client)
	if err != nil {
		return result + fmt.Sprintf(`
Server: Failed to get the server version - %s`, err.Error()), false
	}

	isVersionMatch := version == config.Version
	if isVersionMatch {
		version += `
✔️ Version match`
	}

	return result + fmt.Sprintf(`
Server: %s`, version), isVersionMatch
}

func getServerVersion(ctx context.Context, client *openapi.APIClient) (string, error) {
	resp, _, err := client.ApiApi.
		GetVersion(ctx).
		Execute()
	if err != nil {
		return "", err
	}

	return resp.GetVersion(), nil
}
