package config

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
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
	flags     ConfigFlags
}

func NewConfigurator(resources *resourcemanager.Registry) Configurator {
	ui := cliUI.DefaultUI
	onFinish := func(_ context.Context, _ Config) {
		ui.Success("Successfully configured Tracetest CLI")
		ui.Finish()
	}
	flags := ConfigFlags{}

	return Configurator{resources, ui, onFinish, flags}
}

func (c Configurator) WithOnFinish(onFinish onFinishFn) Configurator {
	c.onFinish = onFinish
	return c
}

func (c Configurator) Start(ctx context.Context, prev Config, flags ConfigFlags) error {
	c.flags = flags
	var serverURL string
	if flags.Endpoint != "" {
		serverURL = flags.Endpoint
	} else if prev.UIEndpoint != "" {
		serverURL = prev.UIEndpoint
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

	if strings.Contains(serverURL, DefaultCloudDomain) {
		path = &DefaultCloudPath
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

	serverType := version.GetType()
	if serverType == "oss" {
		err := Save(cfg)
		if err != nil {
			return fmt.Errorf("could not save configuration: %w", err)
		}

		c.ui.Success("Successfully configured Tracetest CLI")
		return nil
	}

	cfg.AgentEndpoint = version.GetAgentEndpoint()
	cfg.UIEndpoint = version.GetUiEndpoint()
	cfg.Scheme, cfg.Endpoint, cfg.ServerPath, err = ParseServerURL(version.GetApiEndpoint())
	if err != nil {
		return fmt.Errorf("cannot parse server url: %w", err)
	}

	if flags.CI {
		return Save(cfg)
	}

	oauthEndpoint := fmt.Sprintf("%s%s", cfg.URL(), cfg.Path())

	if prev.Jwt != "" {
		cfg.Jwt = prev.Jwt
		cfg.Token = prev.Token
	}

	if flags.Token != "" {
		jwt, err := oauth.ExchangeToken(oauthEndpoint, flags.Token)
		if err != nil {
			return err
		}

		cfg.Jwt = jwt
		cfg.Token = flags.Token

		claims, err := GetTokenClaims(jwt)
		if err != nil {
			return err
		}

		flags.OrganizationID = claims["organization_id"].(string)
		flags.EnvironmentID = claims["environment_id"].(string)
	}

	if flags.AgentApiKey != "" {
		cfg.AgentApiKey = flags.AgentApiKey
		c.ShowOrganizationSelector(ctx, cfg, flags)
		return nil
	}

	if cfg.Jwt != "" {
		c.ShowOrganizationSelector(ctx, cfg, flags)
		return nil
	}

	oauthServer := oauth.NewOAuthServer(oauthEndpoint, cfg.UIEndpoint)
	return oauthServer.WithOnSuccess(c.onOAuthSuccess(ctx, cfg)).
		WithOnFailure(c.onOAuthFailure).
		GetAuthJWT()
}

func (c Configurator) onOAuthSuccess(ctx context.Context, cfg Config) func(token, jwt string) {
	return func(token, jwt string) {
		cfg.Jwt = jwt
		cfg.Token = token

		c.ShowOrganizationSelector(ctx, cfg, c.flags)
	}
}

func (c Configurator) onOAuthFailure(err error) {
	c.ui.Exit(err.Error())
}

func (c Configurator) ShowOrganizationSelector(ctx context.Context, cfg Config, flags ConfigFlags) {
	cfg.OrganizationID = flags.OrganizationID
	if cfg.OrganizationID == "" && flags.AgentApiKey == "" {
		orgID, err := c.organizationSelector(ctx, cfg)
		if err != nil {
			c.ui.Exit(err.Error())
			return
		}

		cfg.OrganizationID = orgID
	}

	cfg.EnvironmentID = flags.EnvironmentID
	if cfg.EnvironmentID == "" && flags.AgentApiKey == "" {
		envID, err := c.environmentSelector(ctx, cfg)
		if err != nil {
			c.ui.Exit(err.Error())
			return
		}

		cfg.EnvironmentID = envID
	}

	err := Save(cfg)
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

func GetTokenClaims(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
