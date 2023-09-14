package config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	cliUI "github.com/kubeshop/tracetest/cli/ui"
)

type Entry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c Configurator) organizationSelector(ctx context.Context, cfg Config) (string, error) {
	resource, err := c.resources.Get("organization")
	if err != nil {
		return "", err
	}

	elements, err := getElements(ctx, resource, cfg)
	if err != nil {
		return "", err
	}

	if len(elements) == 1 {
		c.ui.Println(fmt.Sprintf("Defaulting to only available Organization: %s", elements[0].Name))
		return elements[0].ID, nil
	}

	orgID := ""
	options := make([]cliUI.Option, len(elements))
	for i, org := range elements {
		options[i] = cliUI.Option{
			Text: org.Name,
			Fn: func(o Entry) func(ui cliUI.UI) {
				return func(ui cliUI.UI) {
					orgID = o.ID
				}
			}(org),
		}
	}

	option := c.ui.Select("What Organization do you want to use?", options, 0)
	option.Fn(c.ui)

	return orgID, nil
}

func (c Configurator) environmentSelector(ctx context.Context, cfg Config) (string, error) {
	resource, err := c.resources.Get("env")
	if err != nil {
		return "", err
	}
	resource = resource.WithOptions(resourcemanager.WithPrefixGetter(func() string {
		return fmt.Sprintf("/organizations/%s/", cfg.OrganizationID)
	}))

	elements, err := getElements(ctx, resource, cfg)
	if err != nil {
		return "", err
	}

	if len(elements) == 1 {
		c.ui.Println(fmt.Sprintf("Defaulting to only available Environment: %s", elements[0].Name))
		return elements[0].ID, nil
	}

	envID := ""
	options := make([]cliUI.Option, len(elements))
	for i, env := range elements {
		options[i] = cliUI.Option{
			Text: env.Name,
			Fn: func(e Entry) func(ui cliUI.UI) {
				return func(ui cliUI.UI) {
					envID = e.ID
				}
			}(env),
		}
	}

	option := c.ui.Select("What Environment do you want to use?", options, 0)
	option.Fn(c.ui)

	return envID, err
}

type entryList struct {
	Elements []Entry `json:"elements"`
}

func getElements(ctx context.Context, resource resourcemanager.Client, cfg Config) ([]Entry, error) {
	resource = resource.WithHttpClient(SetupHttpClient(cfg))

	var list entryList
	resultFormat, err := resourcemanager.Formats.GetWithFallback("json", "json")
	if err != nil {
		return []Entry{}, err
	}

	envs, err := resource.List(ctx, resourcemanager.ListOption{}, resultFormat)
	if err != nil {
		return []Entry{}, err
	}

	err = json.Unmarshal([]byte(envs), &list)
	if err != nil {
		return []Entry{}, err
	}

	return list.Elements, nil
}
