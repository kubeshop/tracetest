package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/pterm/pterm"
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
	versionMatch, cliOutdated, err := versionMatch(Version, version)
	if err != nil {
		return result + err.Error(), false
	}

	if versionMatch {
		version += `
✔️ Version match`
	}

	if cliOutdated {
		result += pterm.Red(" (outdated)")
	}

	return result + fmt.Sprintf(`
Server: %s`, version), versionMatch
}

func versionMatch(cliVersion, serverVersion string) (bool, bool, error) {
	if !isSemver(serverVersion) || !isSemver(cliVersion) {
		// if either version is not semver, we can't compare them
		// do a basic strict compare
		return serverVersion == cliVersion, false, nil
	}

	serverSemVer, err := semver.NewVersion(serverVersion)
	if err != nil {
		return false, false, fmt.Errorf("server: Failed to parse the server version - %w`", err)
	}

	cliSemVer, err := semver.NewVersion(cliVersion)
	if err != nil {
		return false, false, fmt.Errorf("failed to parse the CLI version - %w", err)
	}

	versionConstrait, err := semver.NewConstraint(fmt.Sprintf(">=%d.%d", cliSemVer.Major(), cliSemVer.Minor()))
	if err != nil {
		return false, false, fmt.Errorf("failed to parse the CLI version constraint - %w", err)
	}

	outdated := false
	serverVersionDifference := serverSemVer.Compare(cliSemVer)
	if serverVersionDifference > 0 {
		outdated = true
	}

	return versionConstrait.Check(serverSemVer), outdated, nil
}

func isSemver(version string) bool {
	_, err := semver.NewVersion(version)
	return !errors.Is(err, semver.ErrInvalidSemVer)
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
