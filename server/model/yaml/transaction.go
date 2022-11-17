package yaml

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
)

type Transaction struct {
	ID          string            `mapstructure:"id"`
	Name        string            `mapstructure:"name"`
	Description string            `mapstructure:"description" yaml:",omitempty"`
	Env         map[string]string `mapstructure:"env" yaml:",omitempty"`
	Steps       []string          `mapstructure:"steps"`
}

func (t Transaction) Model() model.Transaction {
	steps := make([]model.Test, 0, len(t.Steps))
	for _, stepID := range t.Steps {
		steps = append(steps, model.Test{
			ID: id.ID(stepID),
		})
	}

	mt := model.Transaction{
		ID:          id.ID(t.ID),
		Name:        t.Name,
		Description: t.Description,
		Steps:       steps,
	}

	return mt
}

func (t Transaction) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("transaction name cannot be empty")
	}

	return nil
}
