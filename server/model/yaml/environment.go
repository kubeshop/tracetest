package yaml

import (
	"fmt"

	dc "github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
)

type Environment struct {
	ID          string             `mapstructure:"id"`
	Name        string             `mapstructure:"name"`
	Description string             `mapstructure:"description" yaml:",omitempty"`
	Values      []EnvironmentValue `mapstructure:"values"`
}

type EnvironmentValue struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}

func (e Environment) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("environment name cannot be empty")
	}

	for _, v := range e.Values {
		if v.Key == "" {
			return fmt.Errorf("environment value name cannot be empty")
		}
	}

	return nil
}

func (e Environment) Model() model.Environment {
	me := model.Environment{}
	dc.DeepCopy(e, &me)

	return me
}
