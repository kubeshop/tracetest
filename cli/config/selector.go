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

func (c Configurator) organizationSelector(ctx context.Context, cfg Config, prev *Config) (string, error) {
	resource, err := c.resources.Get("organization")
	if err != nil {
		return "", err
	}

	elements, err := getElements(ctx, resource, cfg)
	if err != nil {
		return "", err
	}

	if len(elements) == 1 {
		c.ui.Println(fmt.Sprintf(`
Defaulting to only available Organization: %s`, elements[0].Name))
		return elements[0].ID, nil
	}

	orgID := ""
	options := make([]cliUI.Option, len(elements))

	defaultIndex := 0
	for i, org := range elements {
		options[i] = cliUI.Option{
			Text: fmt.Sprintf("%s (%s)", org.Name, org.ID),
			Fn: func(o Entry) func(ui cliUI.UI) {
				return func(ui cliUI.UI) {
					orgID = o.ID
				}
			}(org),
		}

		// if we have a previous organization, set it as default
		if prev != nil && prev.OrganizationID == org.ID {
			defaultIndex = i
		}
	}

	option := c.ui.Select(`What Organization do you want to use?`, options, defaultIndex)
	option.Fn(c.ui)

	return orgID, nil
}

func (c Configurator) environmentSelector(ctx context.Context, cfg Config, prev *Config) (string, error) {
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
	defaultIndex := 0
	for i, env := range elements {
		options[i] = cliUI.Option{
			Text: fmt.Sprintf("%s (%s)", env.Name, env.ID),
			Fn: func(e Entry) func(ui cliUI.UI) {
				return func(ui cliUI.UI) {
					envID = e.ID
				}
			}(env),
		}
		// if we have a previous env, set it as default
		if prev != nil && prev.EnvironmentID == env.ID {
			defaultIndex = i
		}
	}

	option := c.ui.Select("What Environment do you want to use?", options, defaultIndex)
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
