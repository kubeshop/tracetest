package yaml

import (
	"fmt"

	dc "github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
)

type Transaction struct {
	ID          string            `mapstructure:"id"`
	Name        string            `mapstructure:"name"`
	Description string            `mapstructure:"description" yaml:",omitempty"`
	Env         map[string]string `mapstructure:"env" yaml:",omitempty"`
	Steps       []string          `mapstructure:"steps"`
}

func (t Transaction) Model() model.Transaction {
	mt := model.Transaction{}
	dc.DeepCopy(t, &mt)
	steps := make([]model.Test, 0, len(t.Steps))
	for _, stepID := range t.Steps {
		steps = append(steps, model.Test{
			ID: id.ID(stepID),
		})
	}
	mt.Steps = steps

	return mt
}

func (t Transaction) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("transaction name cannot be empty")
	}

	return nil
}
