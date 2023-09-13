package config

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/openapi"
)

const defaultVersionExtension = "json"

func GetVersion(ctx context.Context, cfg Config, client *openapi.APIClient) (string, bool) {
	result := fmt.Sprintf(`CLI: %s`, Version)

	if cfg.IsEmpty() {
		return result + `
Server: Not Configured`, false
	}

	version, err := getServerVersion(ctx, client)
	if err != nil {
		return result + fmt.Sprintf(`
Server: Failed to get the server version - %s`, err.Error()), false
	}

	isVersionMatch := version == Version
	if isVersionMatch {
		version += `
✔️ Version match`
	}

	return result + fmt.Sprintf(`
Server: %s`, version), isVersionMatch
}

func getServerVersion(ctx context.Context, client *openapi.APIClient) (string, error) {
	resp, _, err := client.ApiApi.
		GetVersion(ctx, defaultVersionExtension).
		Execute()
	if err != nil {
		return "", err
	}

	return resp.GetVersion(), nil
}

func getVersionMetadata(ctx context.Context, client *openapi.APIClient) (*openapi.Version, error) {
	resp, _, err := client.ApiApi.
		GetVersion(ctx, defaultVersionExtension).
		Execute()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
