package runner

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/ui"
)

func askForMissingVars(missingVars []varSet) []varSet {
	ui.DefaultUI.Warning("Some variables are required by one or more tests")
	ui.DefaultUI.Info("Fill the values for each variable:")

	filledVariables := make([]varSet, 0, len(missingVars))

	for _, missingVar := range missingVars {
		answer := missingVar
		answer.UserValue = ui.DefaultUI.TextInput(missingVar.Name, missingVar.DefaultValue)
		filledVariables = append(filledVariables, answer)
	}

	return filledVariables
}

type varSet struct {
	Name         string
	DefaultValue string
	UserValue    string
}

func (ev varSet) value() string {
	if ev.UserValue != "" {
		return ev.UserValue
	}

	return ev.DefaultValue
}

type varSets []varSet

func (evs varSets) toOpenapi() []openapi.VariableSetValue {
	vars := make([]openapi.VariableSetValue, len(evs))
	for i, ev := range evs {
		vars[i] = openapi.VariableSetValue{
			Key:   openapi.PtrString(ev.Name),
			Value: openapi.PtrString(ev.value()),
		}
	}

	return vars
}

func (evs varSets) unique() varSets {
	seen := make(map[string]bool)
	vars := make(varSets, 0, len(evs))
	for _, ev := range evs {
		if seen[ev.Name] {
			continue
		}

		seen[ev.Name] = true
		vars = append(vars, ev)
	}

	return vars
}

type missingVarsError varSets

func (e missingVarsError) Error() string {
	return fmt.Sprintf("missing variables: %v", []varSet(e))
}

func (e missingVarsError) Is(target error) bool {
	_, ok := target.(missingVarsError)
	return ok
}

func buildMissingVarsError(body []byte) error {
	var missingVarsErrResp openapi.MissingVariablesError
	err := jsonFormat.Unmarshal(body, &missingVarsErrResp)
	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %w", err)
	}

	missingVars := varSets{}

	for _, missingVarErr := range missingVarsErrResp.MissingVariables {
		for _, missingVar := range missingVarErr.Variables {
			missingVars = append(missingVars, varSet{
				Name:         missingVar.GetKey(),
				DefaultValue: missingVar.GetDefaultValue(),
			})
		}
	}

	return missingVarsError(missingVars.unique())
}
