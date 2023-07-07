package runner

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/ui"
)

func askForMissingVars(missingVars []envVar) []envVar {
	ui.DefaultUI.Warning("Some variables are required by one or more tests")
	ui.DefaultUI.Info("Fill the values for each variable:")

	filledVariables := make([]envVar, 0, len(missingVars))

	for _, missingVar := range missingVars {
		answer := missingVar
		answer.UserValue = ui.DefaultUI.TextInput(missingVar.Name, missingVar.DefaultValue)
		filledVariables = append(filledVariables, answer)
	}

	return filledVariables
}

type envVar struct {
	Name         string
	DefaultValue string
	UserValue    string
}

func (ev envVar) value() string {
	if ev.UserValue != "" {
		return ev.UserValue
	}

	return ev.DefaultValue
}

type envVars []envVar

func (evs envVars) toOpenapi() []openapi.EnvironmentValue {
	vars := make([]openapi.EnvironmentValue, len(evs))
	for i, ev := range evs {
		vars[i] = openapi.EnvironmentValue{
			Key:   openapi.PtrString(ev.Name),
			Value: openapi.PtrString(ev.value()),
		}
	}

	return vars
}

func (evs envVars) unique() envVars {
	seen := make(map[string]bool)
	vars := make(envVars, 0, len(evs))
	for _, ev := range evs {
		if seen[ev.Name] {
			continue
		}

		seen[ev.Name] = true
		vars = append(vars, ev)
	}

	return vars
}

type missingEnvVarsError envVars

func (e missingEnvVarsError) Error() string {
	return fmt.Sprintf("missing env vars: %v", []envVar(e))
}
