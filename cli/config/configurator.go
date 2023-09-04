package config

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/pkg/oauth"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	cliUI "github.com/kubeshop/tracetest/cli/ui"
)

type onFinishFn func(context.Context, Config)

type Configurator struct {
	resources *resourcemanager.Registry
	ui        cliUI.UI
	onFinish  onFinishFn
}

func NewConfigurator(resources *resourcemanager.Registry) Configurator {
	ui := cliUI.DefaultUI
	onFinish := func(_ context.Context, _ Config) {
		ui.Success("Successfully configured Tracetest CLI")
		ui.Finish()
	}

	return Configurator{resources, ui, onFinish}
}

func (c Configurator) WithOnFinish(onFinish onFinishFn) Configurator {
	c.onFinish = onFinish
	return c
}

func (c Configurator) Start(ctx context.Context, prev Config, flags ConfigFlags) error {
	var serverURL string
	if flags.Endpoint != "" {
		serverURL = flags.Endpoint
	} else {
		path := ""
		if prev.ServerPath != nil {
			path = *prev.ServerPath
		}
		serverURL = c.ui.TextInput("Enter your Tracetest server URL", fmt.Sprintf("%s%s", prev.URL(), path))
	}

	if err := ValidateServerURL(serverURL); err != nil {
		return err
	}

	scheme, endpoint, path, err := ParseServerURL(serverURL)
	if err != nil {
		return err
	}

	cfg := Config{
		Scheme:     scheme,
		Endpoint:   endpoint,
		ServerPath: path,
	}

	client := GetAPIClient(cfg)
	version, err := getVersionMetadata(ctx, client)
	if err != nil {
		return fmt.Errorf("cannot get version metadata: %w", err)
	}

	cfg.AgentEndpoint = version.GetAgentEndpoint()
	cfg.UIEndpoint = version.GetUiEndpoint()

	serverType := version.GetType()
	if serverType == "oss" {
		err := Save(ctx, cfg)
		if err != nil {
			return fmt.Errorf("could not save configuration: %w", err)
		}

		c.ui.Success("Successfully configured Tracetest CLI")
		return nil
	}

	if prev.Jwt != "" {
		cfg.Jwt = prev.Jwt
		cfg.Token = prev.Token

		c.ShowOrganizationSelector(ctx, cfg)
		return nil
	}

	oauthServer := oauth.NewOAuthServer(cfg.URL(), cfg.UIEndpoint)
	err = oauthServer.WithOnSuccess(c.onOAuthSuccess(ctx, cfg)).
		WithOnFailure(c.onOAuthFailure).
		GetAuthJWT()

	return err
}

func (c Configurator) onOAuthSuccess(ctx context.Context, cfg Config) func(token, jwt string) {
	return func(token, jwt string) {
		cfg.Jwt = jwt
		cfg.Token = token

		c.ShowOrganizationSelector(ctx, cfg)
	}
}

func (c Configurator) onOAuthFailure(err error) {
	c.ui.Exit(err.Error())
}

func (c Configurator) ShowOrganizationSelector(ctx context.Context, cfg Config) {
	cfg, err := c.organizationSelector(ctx, cfg)
	if err != nil {
		c.ui.Exit(err.Error())
		return
	}

	cfg, err = c.environmentSelector(ctx, cfg)
	if err != nil {
		c.ui.Exit(err.Error())
		return
	}

	err = Save(ctx, cfg)
	if err != nil {
		c.ui.Exit(err.Error())
		return
	}

	c.onFinish(ctx, cfg)
}

func (c Configurator) ShowEnvironmentSelector(ctx context.Context, cfg Config) {
	cfg, err := c.environmentSelector(ctx, cfg)
	if err != nil {
		c.ui.Exit(err.Error())
		return
	}

	err = Save(ctx, cfg)
	if err != nil {
		c.ui.Exit(err.Error())
		return
	}

	c.onFinish(ctx, cfg)
}

func SetupHttpClient(cfg Config) *resourcemanager.HTTPClient {
	extraHeaders := http.Header{}
	extraHeaders.Set("x-client-id", analytics.ClientID())
	extraHeaders.Set("x-source", "cli")
	extraHeaders.Set("x-organization-id", cfg.OrganizationID)
	extraHeaders.Set("x-environment-id", cfg.EnvironmentID)
	extraHeaders.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Jwt))

	return resourcemanager.NewHTTPClient(fmt.Sprintf("%s%s", cfg.URL(), cfg.Path()), extraHeaders)
}
