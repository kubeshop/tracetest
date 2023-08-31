package config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	cliUI "github.com/kubeshop/tracetest/cli/ui"
)

type entry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c Configurator) organizationSelector(ctx context.Context, cfg Config) (Config, error) {
	resource, err := c.resources.Get("organization")
	if err != nil {
		return cfg, err
	}

	elements, err := getElements(ctx, resource, cfg)
	if err != nil {
		return cfg, err
	}

	options := make([]cliUI.Option, len(elements))
	for i, org := range elements {
		options[i] = cliUI.Option{
			Text: org.Name,
			Fn: func(o entry) func(ui cliUI.UI) {
				return func(ui cliUI.UI) {
					cfg.OrganizationID = o.ID
				}
			}(org),
		}
	}

	option := c.ui.Select("What Organization do you want to use?", options, 0)
	option.Fn(c.ui)
	return cfg, nil
}

func (c Configurator) environmentSelector(ctx context.Context, cfg Config) (Config, error) {
	resource, err := c.resources.Get("env")
	if err != nil {
		return cfg, err
	}
	resource = resource.WithOptions(resourcemanager.WithPrefixGetter(func() string {
		return fmt.Sprintf("/organizations/%s/", cfg.OrganizationID)
	}))

	elements, err := getElements(ctx, resource, cfg)
	if err != nil {
		return cfg, err
	}

	options := make([]cliUI.Option, len(elements))
	for i, env := range elements {
		options[i] = cliUI.Option{
			Text: env.Name,
			Fn: func(e entry) func(ui cliUI.UI) {
				return func(ui cliUI.UI) {
					cfg.EnvironmentID = e.ID
				}
			}(env),
		}
	}

	option := c.ui.Select("What Environment do you want to use?", options, 0)
	option.Fn(c.ui)
	return cfg, err
}

type entryList struct {
	Elements []entry `json:"elements"`
}

func getElements(ctx context.Context, resource resourcemanager.Client, cfg Config) ([]entry, error) {
	resource = resource.WithHttpClient(setupHttpClient(cfg))

	var list entryList
	resultFormat, err := resourcemanager.Formats.GetWithFallback("json", "json")
	if err != nil {
		return []entry{}, err
	}

	envs, err := resource.List(ctx, resourcemanager.ListOption{}, resultFormat)
	if err != nil {
		return []entry{}, err
	}

	err = json.Unmarshal([]byte(envs), &list)
	if err != nil {
		return []entry{}, err
	}

	return list.Elements, nil
}
