package config

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/pkg/oauth"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	cliUI "github.com/kubeshop/tracetest/cli/ui"
	"go.uber.org/zap"
)

type onFinishFn func(context.Context, Config)

type Configurator struct {
	logger         *zap.Logger
	resources      *resourcemanager.Registry
	ui             cliUI.UI
	onFinish       onFinishFn
	errorHandlerFn errorHandlerFn
	flags          *agentConfig.Flags
	finalServerURL string
}

func NewConfigurator(resources *resourcemanager.Registry) Configurator {
	ui := cliUI.DefaultUI

	return Configurator{
		logger:    zap.NewNop(),
		resources: resources,
		ui:        ui,
		onFinish: func(_ context.Context, _ Config) {
			ui.Success("Successfully configured Tracetest CLI")
			ui.Finish()
		},
		errorHandlerFn: func(ctx context.Context, err error) {
			ui.Exit(err.Error())
		},
		flags: &agentConfig.Flags{},
	}
}

func (c Configurator) WithLogger(logger *zap.Logger) Configurator {
	c.logger = logger
	return c
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
	c.flags = &flags
	var serverURL string

	c.logger.Debug("Starting configurator", zap.Any("flags", flags), zap.Any("prev", prev), zap.String("serverURL", serverURL))

	if c.flags.AutomatedEnvironmentCanBeInferred() {
		c.logger.Debug("Automated environment detected, skipping prompts")
		// avoid prompts on automated or non-interactive environments
		serverURL = c.lastUsedURL(prev)
	} else {
		c.logger.Debug("Interactive environment detected, prompting for server URL")
		var err error
		serverURL, err = c.getServerURL(prev)
		if err != nil {
			c.logger.Error("Invalid server URL", zap.Error(err))
			return err
		}
	}
	c.finalServerURL = serverURL
	c.logger.Debug("Final server URL", zap.String("serverURL", serverURL))

	cfg, err := c.createConfig(serverURL)
	if err != nil {
		c.logger.Error("Could not create config", zap.Error(err))
		return err
	}

	cfg, err, isOSS := c.populateConfigWithVersionInfo(ctx, cfg)
	if err != nil {
		c.logger.Error("Could not populate config with version info", zap.Error(err))
		return err
	}

	if isOSS {
		c.logger.Debug("OSS server detected, skipping OAuth")
		// we don't need anything else for OSS
		return nil
	}

	if c.flags.CI {
		c.logger.Debug("CI environment detected, skipping OAuth")
		_, err = Save(ctx, cfg)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = c.handleOAuth(ctx, cfg, prev)
	if err != nil {
		c.logger.Error("Could not handle OAuth", zap.Error(err))
		return err
	}

	c.logger.Debug("Successfully configured OAuth")

	return nil
}

func (c Configurator) lastUsedURL(prev *Config) string {
	if c.flags.ServerURL != "" {
		return c.flags.ServerURL
	}

	possibleValues := []string{}
	if prev != nil {
		possibleValues = append(possibleValues, prev.UIEndpoint, prev.URL())
	}
	possibleValues = append(possibleValues, DefaultCloudEndpoint)

	return getFirstNonEmptyString(possibleValues)
}

func (c Configurator) getServerURL(prev *Config) (string, error) {
	serverURL := c.flags.ServerURL

	// if flag was passed, don't show prompt
	if c.flags.ServerURL == "" {
		serverURL = c.ui.TextInput("What tracetest server do you want to use?", c.lastUsedURL(prev))
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

func (c Configurator) populateConfigWithDevConfig(ctx context.Context, cfg *Config) {
	cfg.AgentEndpoint = os.Getenv("TRACETEST_DEV_AGENT_ENDPOINT")
	if cfg.AgentEndpoint == "" {
		cfg.AgentEndpoint = "localhost:8091"
	}

	cfg.UIEndpoint = os.Getenv("TRACETEST_DEV_UI_ENDPOINT")
	if cfg.UIEndpoint == "" {
		cfg.UIEndpoint = "http://localhost:3000"
	}

	cfg.Scheme = os.Getenv("TRACETEST_DEV_SCHEME")
	if cfg.Scheme == "" {
		cfg.Scheme = "http"
	}

	cfg.Endpoint = os.Getenv("TRACETEST_DEV_ENDPOINT")
	if cfg.Endpoint == "" {
		cfg.Endpoint = "localhost:11633"
	}

	cfg.ServerPath = os.Getenv("TRACETEST_DEV_SERVER_PATH")
}

func (c Configurator) populateConfigWithVersionInfo(ctx context.Context, cfg Config) (_ Config, _ error, isOSS bool) {
	useDevVersion := os.Getenv("TRACETEST_AGENT_DEV_CONFIG") == "true"
	if useDevVersion && Version == "dev" {
		c.populateConfigWithDevConfig(ctx, &cfg)

		c.ui.Success("Configured Tracetest CLI in development mode")

		return cfg, nil, false
	}

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

func (c Configurator) handleOAuth(ctx context.Context, cfg Config, prev *Config) (Config, error) {
	if prev != nil && cfg.UIEndpoint == prev.UIEndpoint {
		c.logger.Debug("Using previous UI endpoint", zap.String("uiEndpoint", cfg.UIEndpoint))
		if prev != nil && prev.Jwt != "" {
			c.logger.Debug("Using previous JWT")
			cfg.Jwt = prev.Jwt
			cfg.Token = prev.Token
		}
	}

	if c.flags.Token != "" {
		c.logger.Debug("Using token from flag")
		var err error
		cfg, err = c.exchangeToken(cfg, c.flags.Token)
		if err != nil {
			c.logger.Error("Could not exchange token", zap.Error(err))
			return Config{}, err
		}
	}

	if c.flags.AgentApiKey != "" {
		c.logger.Debug("Using agent API key from flag")
		cfg.AgentApiKey = c.flags.AgentApiKey
		c.showOrganizationSelector(ctx, prev, cfg)
		return cfg, nil
	}

	if cfg.Jwt != "" {
		c.logger.Debug("Using JWT from config")
		c.showOrganizationSelector(ctx, prev, cfg)
		return cfg, nil
	}

	c.logger.Debug("No JWT found, executing user login")

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
	c.logger.Debug("Exchanging token", zap.String("token", token))
	jwt, err := oauth.ExchangeToken(cfg.OAuthEndpoint(), token)
	if err != nil {
		c.logger.Error("Could not exchange token", zap.Error(err))
		return Config{}, err
	}

	cfg.Jwt = jwt
	cfg.Token = token

	claims, err := GetTokenClaims(jwt)
	if err != nil {
		c.logger.Error("Could not get token claims", zap.Error(err))
		return Config{}, err
	}

	organizationId := claims["organization_id"].(string)
	environmentId := claims["environment_id"].(string)

	if organizationId != "" {
		c.logger.Debug("Using organization ID from token", zap.String("organizationID", organizationId))
		c.flags.OrganizationID = organizationId
	}
	if environmentId != "" {
		c.logger.Debug("Using environment ID from token", zap.String("environmentID", environmentId))
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
		c.logger.Debug("OAuth success")
		cfg.Jwt = jwt
		cfg.Token = token

		c.showOrganizationSelector(ctx, prev, cfg)
	}
}

func (c Configurator) onOAuthFailure(err error) {
	c.errorHandlerFn(context.Background(), err)
}

func (c Configurator) showOrganizationSelector(ctx context.Context, prev *Config, cfg Config) {
	c.logger.Debug("Showing organization selector", zap.String("organizationID", cfg.OrganizationID), zap.String("environmentID", cfg.EnvironmentID))
	cfg.OrganizationID = c.flags.OrganizationID
	if cfg.OrganizationID == "" && c.flags.AgentApiKey == "" {
		c.logger.Debug("Organization ID not found, prompting for organization")
		orgID, err := c.organizationSelector(ctx, cfg, prev)
		if err != nil {
			c.logger.Error("Could not select organization", zap.Error(err))
			c.errorHandlerFn(ctx, err)
			return
		}

		cfg.OrganizationID = orgID
	}

	cfg.EnvironmentID = c.flags.EnvironmentID
	if cfg.EnvironmentID == "" && c.flags.AgentApiKey == "" {
		c.logger.Debug("Environment ID not found, prompting for environment")
		envID, err := c.environmentSelector(ctx, cfg, prev)
		if err != nil {
			c.errorHandlerFn(ctx, err)
			return
		}

		cfg.EnvironmentID = envID
	}

	ctx, err := Save(ctx, cfg)
	if err != nil {
		c.logger.Error("Could not save configuration", zap.Error(err))
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
