package model

import (
	"strings"
	"time"
)

type (
	Environment struct {
		ID          string
		Name        string
		Description string
		CreatedAt   time.Time
		Values      []EnvironmentValue
	}

	EnvironmentValue struct {
		Key   string
		Value string
	}
)

func (e Environment) HasID() bool {
	return e.ID != ""
}

func (e Environment) GetSlug() string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(e.Name), " ", "-"))
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
