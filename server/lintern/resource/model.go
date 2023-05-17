package lintern_resource

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/pkg/id"
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
