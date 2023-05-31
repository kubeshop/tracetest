package linter_resource

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"golang.org/x/exp/slices"
)

type (
	Linter struct {
		ID           id.ID          `json:"id"`
		Name         string         `json:"name"`
		Enabled      bool           `json:"enabled"`
		MinimumScore int            `json:"minimumScore"`
		Plugins      []LinterPlugin `json:"plugins"`
	}

	LinterPlugin struct {
		Name     string `json:"name"`
		Enabled  bool   `json:"enabled"`
		Required bool   `json:"required"`
	}
)

func (l Linter) Validate() error {
	if l.Name == "" {
		return fmt.Errorf("linter name cannot be empty")
	}

	for _, p := range l.Plugins {
		if p.Name == "" {
			return fmt.Errorf("plugin name cannot be empty")
		}
	}

	return nil
}

func (l Linter) HasID() bool {
	return l.ID != ""
}

func (l Linter) GetID() id.ID {
	return l.ID
}

func (l Linter) ValidateResult(result model.LinterResult) error {
	if l.MinimumScore != 0 && result.Score < l.MinimumScore {
		return fmt.Errorf("linter score validation failed. Minimum %d, Actual: %d", l.MinimumScore, result.Score)
	}

	failedPlugins := make([]string, 0)
	for _, plugin := range result.Plugins {
		if !plugin.Passed {
			failedPlugins = append(failedPlugins, plugin.Name)
		}
	}

	if len(failedPlugins) == 0 {
		return nil
	}

	requiredPlugins := make([]string, 0)
	for _, plugin := range l.Plugins {
		if plugin.Required {
			requiredPlugins = append(requiredPlugins, plugin.Name)
		}
	}

	for _, plugin := range requiredPlugins {
		index := slices.IndexFunc(failedPlugins, func(failedPlugin string) bool { return failedPlugin == plugin })
		if index >= 0 {
			return fmt.Errorf("linter failed. Required plugin %s failed", plugin)
		}
	}

	return nil
}
