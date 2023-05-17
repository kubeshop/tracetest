package lintern_resource

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"golang.org/x/exp/slices"
)

type (
	Lintern struct {
		ID           id.ID           `json:"id"`
		Name         string          `json:"name"`
		Enabled      bool            `json:"enabled"`
		MinimumScore int             `json:"minimumScore"`
		Plugins      []LinternPlugin `json:"plugins"`
	}

	LinternPlugin struct {
		Name     string `json:"name"`
		Enabled  bool   `json:"enabled"`
		Required bool   `json:"required"`
	}
)

func (l Lintern) Validate() error {
	if l.Name == "" {
		return fmt.Errorf("lintern name cannot be empty")
	}

	for _, p := range l.Plugins {
		if p.Name == "" {
			return fmt.Errorf("plugin name cannot be empty")
		}
	}

	return nil
}

func (l Lintern) HasID() bool {
	return l.ID != ""
}

func (l Lintern) ValidateResult(result model.LinternResult) error {
	if l.MinimumScore != 0 && result.Score < l.MinimumScore {
		return fmt.Errorf("lintern score validation failed. Minimum %d, Actual: %d", l.MinimumScore, result.Score)
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
			return fmt.Errorf("lintern failed. Required plugin %s failed", plugin)
		}
	}

	return nil
}
