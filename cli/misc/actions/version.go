package misc_actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
)

type VersionConfig struct{}

type versionAction struct {
	actionAgs
}

var _ Action[VersionConfig] = &versionAction{}

func NewGetServerVersionAction(options ...ActionArgsOption) versionAction {
	args := NewActionArgs(options...)
	return versionAction{actionAgs: args}
}

func (a versionAction) Run(ctx context.Context, args VersionConfig) error {
	versionText := a.GetVersionText(ctx)

	fmt.Println(versionText)
	return nil
}

func (a versionAction) GetVersionText(ctx context.Context) string {
	result := fmt.Sprintf(`CLI: %s`, config.Version)

	if a.config.IsEmpty() {
		return result + `
Server: Not Configured`
	}

	version, err := a.GetServerVersion(ctx)
	if err != nil {
		return result + fmt.Sprintf(`
Server: Failed to get the server version - %s`, err.Error())
	}

	isVersionMatch := version == config.Version
	if isVersionMatch {
		version += `
✔️ Version match`
	} else {
		version += `
✖️ Version mismatch`
	}

	return result + fmt.Sprintf(`
Server: %s`, version)
}

func (a versionAction) GetServerVersion(ctx context.Context) (string, error) {
	if a.config.IsEmpty() {
		return "", fmt.Errorf("not Configured")
	}

	req := a.client.ApiApi.GetVersion(ctx)
	version, _, err := req.Execute()
	if err != nil {
		return "", err
	}

	return *version.Version, nil
}
