package variableset

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/pkg/id"
)

type (
	VariableSet struct {
		ID          id.ID              `json:"id"`
		Name        string             `json:"name"`
		Description string             `json:"description"`
		CreatedAt   string             `json:"createdAt"`
		Values      []VariableSetValue `json:"values"`
	}

	VariableSetValue struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
)

func (e VariableSet) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("variable set name cannot be empty")
	}

	for _, v := range e.Values {
		if v.Key == "" {
			return fmt.Errorf("variable set value name cannot be empty")
		}
	}

	return nil
}

func (e VariableSet) HasID() bool {
	return e.ID != ""
}

func (e VariableSet) GetID() id.ID {
	return e.ID
}

func (e VariableSet) Slug() id.ID {
	return id.ID(strings.ToLower(strings.ReplaceAll(strings.TrimSpace(e.Name), " ", "-")))
}

func (e VariableSet) Get(key string) string {
	for _, v := range e.Values {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}

func (e VariableSet) Merge(env VariableSet) VariableSet {
	values := make(map[string]string)
	for _, variable := range e.Values {
		values[variable.Key] = variable.Value
	}

	for _, variable := range env.Values {
		values[variable.Key] = variable.Value
	}

	newValues := make([]VariableSetValue, 0, len(values))
	for key, value := range values {
		newValues = append(newValues, VariableSetValue{
			Key:   key,
			Value: value,
		})
	}

	e.Values = newValues
	return e
}
