package runner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
)

type environment struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	AgentApiKey    string `json:"agentApiKey"`
	OrganizationID string `json:"organizationID"`
}

func (s *Runner) getEnvironment(ctx context.Context, cfg config.Config) (environment, error) {
	resource, err := s.resources.Get("env")
	if err != nil {
		return environment{}, err
	}

	resource = resource.
		WithHttpClient(config.SetupHttpClient(cfg)).
		WithOptions(resourcemanager.WithPrefixGetter(func() string {
			return fmt.Sprintf("/organizations/%s/", cfg.OrganizationID)
		}))

	resultFormat, err := resourcemanager.Formats.GetWithFallback("json", "json")
	if err != nil {
		return environment{}, err
	}

	raw, err := resource.Get(ctx, cfg.EnvironmentID, resultFormat)
	if err != nil {
		return environment{}, err
	}

	var env environment
	err = json.Unmarshal([]byte(raw), &env)
	if err != nil {
		return environment{}, err
	}

	return env, nil
}
