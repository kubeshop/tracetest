package analyzer

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
)

const (
	ResourceName       = "Analyzer"
	ResourceNamePlural = "Analyzers"
)

var Operations = []resourcemanager.Operation{
	resourcemanager.OperationGet,
	resourcemanager.OperationList,
	resourcemanager.OperationUpdate,
}

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
