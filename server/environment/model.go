package environment

import (
	"strings"

	"github.com/kubeshop/tracetest/server/pkg/id"
)

type (
	Environment struct {
		ID          id.ID              `json:"id"`
		Name        string             `json:"name"`
		Description string             `json:"description"`
		CreatedAt   string             `json:"createdAt"`
		Values      []EnvironmentValue `json:"values"`
	}

	EnvironmentValue struct {
		Key   string `json:"key"`
		Value string `json:"value"`
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
