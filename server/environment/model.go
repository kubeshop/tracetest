package environment

import (
	"strings"

	"github.com/kubeshop/tracetest/server/id"
)

type (
	Environment struct {
		ID          id.ID              `mapstructure:"id" json:"id"`
		Name        string             `mapstructure:"name" json:"name"`
		Description string             `mapstructure:"description" json:"description"`
		CreatedAt   string             `mapstructure:"createdAt" json:"createdAt"`
		Values      []EnvironmentValue `mapstructure:"values" json:"values"`
	}

	EnvironmentValue struct {
		Key   string `mapstructure:"key" json:"key"`
		Value string `mapstructure:"value" json:"value"`
	}
)

func (e Environment) Validate() error {
	return nil
}

func (e Environment) HasID() bool {
	return e.ID != ""
}

func (e Environment) Slug() id.ID {
	return id.ID(strings.ToLower(strings.ReplaceAll(strings.TrimSpace(e.Name), " ", "-")))
}

func (e Environment) Get(key string) string {
	for _, v := range e.Values {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}

func (e Environment) Merge(env Environment) Environment {
	values := make(map[string]string)
	for _, variable := range e.Values {
		values[variable.Key] = variable.Value
	}

	for _, variable := range env.Values {
		values[variable.Key] = variable.Value
	}

	newValues := make([]EnvironmentValue, 0, len(values))
	for key, value := range values {
		newValues = append(newValues, EnvironmentValue{
			Key:   key,
			Value: value,
		})
	}

	e.Values = newValues
	return e
}
