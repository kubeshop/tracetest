package config

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/pkg/oauth"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	cliUI "github.com/kubeshop/tracetest/cli/ui"
)

type onFinishFn func(context.Context, Config)

type Configurator struct {
	resources      *resourcemanager.Registry
	ui             cliUI.UI
	onFinish       onFinishFn
	errorHandlerFn errorHandlerFn
	flags          agentConfig.Flags
	finalServerURL string
}

func NewConfigurator(resources *resourcemanager.Registry) Configurator {
	ui := cliUI.DefaultUI

	return Configurator{
		resources: resources,
		ui:        ui,
		onFinish: func(_ context.Context, _ Config) {
			ui.Success("Successfully configured Tracetest CLI")
			ui.Finish()
		},
		errorHandlerFn: func(ctx context.Context, err error) {
			ui.Exit(err.Error())
		},
		flags: agentConfig.Flags{},
	}
}

func (c Configurator) WithOnFinish(onFinish onFinishFn) Configurator {
	c.onFinish = onFinish
	return c
}

type errorHandlerFn func(ctx context.Context, err error)

func (c Configurator) WithErrorHandler(fn errorHandlerFn) Configurator {
	c.errorHandlerFn = fn
	return c
}

func (c Configurator) Start(ctx context.Context, prev *Config, flags agentConfig.Flags) error {
	c.flags = flags
	serverURL, err := c.getServerURL(prev, flags)
	c.finalServerURL = serverURL
	if err != nil {
		return err
	}

	cfg, err := c.createConfig(serverURL)
	if err != nil {
		return err
	}

	cfg, err, isOSS := c.populateConfigWithVersionInfo(ctx, cfg)
	if err != nil {
		return err
	}

	if isOSS {
		// we don't need anything else for OSS
		return nil
	}

	if flags.CI {
		_, err = Save(ctx, cfg)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = c.handleOAuth(ctx, cfg, prev, flags)
	if err != nil {
		return err
	}

	return nil
}

func (c Configurator) getServerURL(prev *Config, flags agentConfig.Flags) (string, error) {
	serverURL := flags.ServerURL

	// if flag was passed, don't show prompt
	if flags.ServerURL == "" {
		possibleValues := []string{}
		if prev != nil {
			possibleValues = append(possibleValues, prev.UIEndpoint, prev.URL())
		}
		possibleValues = append(possibleValues, DefaultCloudEndpoint)

		lastUsed := getFirstNonEmptyString(possibleValues)
		serverURL = c.ui.TextInput("What tracetest server do you want to use?", lastUsed)
	}

	if err := validateServerURL(serverURL); err != nil {
		return "", err
	}

	return serverURL, nil
}

func (c Configurator) createConfig(serverURL string) (Config, error) {
	scheme, endpoint, path, err := ParseServerURL(serverURL)
	if err != nil {
		return Config{}, err
	}

	if strings.Contains(serverURL, DefaultCloudDomain) {
		path = DefaultCloudPath
	} else if !strings.HasSuffix(path, "/api") {
		path = strings.TrimSuffix(path, "/") + "/api"
	}

	return Config{
		Scheme:     scheme,
		Endpoint:   endpoint,
		ServerPath: path,
	}, nil
}

type invalidServerErr struct {
	ui        cliUI.UI
	serverURL string
	parent    error
}

func (e invalidServerErr) Error() string {
	return fmt.Errorf("cannot reach %s: %w", e.serverURL, e.parent).Error()
}

func (e invalidServerErr) Render() {
	msg := fmt.Sprintf(`Cannot reach "%s". Please verify the url and enter it again.`, e.serverURL)
	e.ui.Error(msg)
}

func (c Configurator) populateConfigWithVersionInfo(ctx context.Context, cfg Config) (_ Config, _ error, isOSS bool) {
	client := GetAPIClient(cfg)
	version, err := getVersionMetadata(ctx, client)
	if err != nil {
		err = invalidServerErr{c.ui, c.finalServerURL, fmt.Errorf("cannot get version metadata: %w", err)}
		return Config{}, err, false
	}

	serverType := version.GetType()
	if serverType == "oss" {
		_, err = Save(ctx, cfg)
		if err != nil {
			return Config{}, fmt.Errorf("could not save configuration: %w", err), false
		}

		c.ui.Success("Successfully configured Tracetest CLI")
		return cfg, nil, true
	}

	cfg.AgentEndpoint = version.GetAgentEndpoint()
	cfg.UIEndpoint = version.GetUiEndpoint()
	cfg.Scheme, cfg.Endpoint, cfg.ServerPath, err = ParseServerURL(version.GetApiEndpoint())
	if err != nil {
		return Config{}, fmt.Errorf("cannot parse server url: %w", err), false
	}

	return cfg, nil, false
}

func (c Configurator) handleOAuth(ctx context.Context, cfg Config, prev *Config, flags agentConfig.Flags) (Config, error) {
	if prev != nil && cfg.UIEndpoint == prev.UIEndpoint {
		if prev != nil && prev.Jwt != "" {
			cfg.Jwt = prev.Jwt
			cfg.Token = prev.Token
		}
	}

	if flags.Token != "" {
		var err error
		cfg, err = c.exchangeToken(cfg, flags.Token)
		if err != nil {
			return Config{}, err
		}
	}

	if flags.AgentApiKey != "" {
		cfg.AgentApiKey = flags.AgentApiKey
		c.showOrganizationSelector(ctx, prev, cfg, flags)
		return cfg, nil
	}

	if cfg.Jwt != "" {
		c.showOrganizationSelector(ctx, prev, cfg, flags)
		return cfg, nil
	}

	return c.ExecuteUserLogin(ctx, cfg, prev)
}

func (c Configurator) ExecuteUserLogin(ctx context.Context, cfg Config, prev *Config) (Config, error) {
	oauthServer := oauth.NewOAuthServer(cfg.OAuthEndpoint(), cfg.UIEndpoint)
	err := oauthServer.WithOnSuccess(c.onOAuthSuccess(ctx, cfg, prev)).
		WithOnFailure(c.onOAuthFailure).
		GetAuthJWT()
	if err != nil {
		return Config{}, err
	}

	return cfg, err
}

func (c Configurator) exchangeToken(cfg Config, token string) (Config, error) {
	jwt, err := oauth.ExchangeToken(cfg.OAuthEndpoint(), token)
	if err != nil {
		return Config{}, err
	}

	cfg.Jwt = jwt
	cfg.Token = token

	claims, err := GetTokenClaims(jwt)
	if err != nil {
		return Config{}, err
	}

	organizationId := claims["organization_id"].(string)
	environmentId := claims["environment_id"].(string)

	if organizationId != "" {
		c.flags.OrganizationID = organizationId
	}
	if environmentId != "" {
		c.flags.EnvironmentID = environmentId
	}

	return cfg, nil
}

func getFirstNonEmptyString(values []string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}

	return ""
}

func (c Configurator) onOAuthSuccess(ctx context.Context, cfg Config, prev *Config) func(token, jwt string) {
	return func(token, jwt string) {
		cfg.Jwt = jwt
		cfg.Token = token

		c.showOrganizationSelector(ctx, prev, cfg, c.flags)
	}
}

func (c Configurator) onOAuthFailure(err error) {
	c.errorHandlerFn(context.Background(), err)
}

func (c Configurator) showOrganizationSelector(ctx context.Context, prev *Config, cfg Config, flags agentConfig.Flags) {
	cfg.OrganizationID = flags.OrganizationID
	if cfg.OrganizationID == "" && flags.AgentApiKey == "" {
		orgID, err := c.organizationSelector(ctx, cfg, prev)
		if err != nil {
			c.errorHandlerFn(ctx, err)
			return
		}

		cfg.OrganizationID = orgID
	}

	cfg.EnvironmentID = flags.EnvironmentID
	if cfg.EnvironmentID == "" && flags.AgentApiKey == "" {
		envID, err := c.environmentSelector(ctx, cfg, prev)
		if err != nil {
			c.errorHandlerFn(ctx, err)
			return
		}

		cfg.EnvironmentID = envID
	}

	ctx, err := Save(ctx, cfg)
	if err != nil {
		c.errorHandlerFn(ctx, err)
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
