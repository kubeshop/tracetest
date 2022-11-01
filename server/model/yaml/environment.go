package yaml

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
)

type Environment struct {
	ID          string            `mapstructure:"id"`
	Name        string            `mapstructure:"name"`
	Description string            `mapstructure:"description" yaml:",omitempty"`
	Values      EnvironmentValues `mapstructure:"values"`
}

type EnvironmentValue struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}

type EnvironmentValues []EnvironmentValue

func (evs EnvironmentValues) Model() []model.EnvironmentValue {
	mevs := make([]model.EnvironmentValue, 0, len(evs))
	for _, ev := range evs {
		mevs = append(mevs, model.EnvironmentValue{
			Key:   ev.Key,
			Value: ev.Value,
		})
	}

	return mevs
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
	me := model.Environment{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Values:      e.Values.Model(),
	}

	return me
}
