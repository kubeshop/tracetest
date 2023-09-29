package config

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/openapi"
)

const defaultVersionExtension = "json"

func GetVersion(ctx context.Context, cfg Config) (string, bool) {
	result := fmt.Sprintf(`CLI: %s`, Version)

	if cfg.UIEndpoint != "" {
		scheme, endpoint, path, _ := ParseServerURL(cfg.UIEndpoint)
		cfg.Scheme = scheme
		cfg.Endpoint = endpoint
		cfg.ServerPath = path
	}
	client := GetAPIClient(cfg)

	if cfg.IsEmpty() {
		return result + `
Server: Not Configured`, false
	}

	meta, err := getVersionMetadata(ctx, client)
	if err != nil {
		return result + fmt.Sprintf(`
Server: Failed to get the server version - %s`, err.Error()), false
	}

	version := meta.GetVersion()
	isVersionMatch := version == Version
	if isVersionMatch {
		version += `
✔️ Version match`
	}

	return result + fmt.Sprintf(`
Server: %s`, version), isVersionMatch
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
