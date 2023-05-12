package model

import (
	"strings"
	"time"
)

type (
	Environment2 struct {
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

func (e Environment2) HasID() bool {
	return e.ID != ""
}

func (e Environment2) Slug() string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(e.Name), " ", "-"))
}

func (e Environment2) Get(key string) string {
	for _, v := range e.Values {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}

func (e Environment2) Merge(env Environment2) Environment2 {
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
