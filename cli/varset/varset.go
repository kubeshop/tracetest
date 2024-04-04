package varset

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/ui"
)

var jsonFormat = resourcemanager.Formats.Get(resourcemanager.FormatJSON)

func AskForMissingVars(missingVars []VarSet) VarSets {
	ui.DefaultUI.Warning("Some variables are required by one or more tests")
	ui.DefaultUI.Info("Fill the values for each variable:")

	filledVariables := make([]VarSet, 0, len(missingVars))

	for _, missingVar := range missingVars {
		answer := missingVar
		answer.UserValue = ui.DefaultUI.TextInput(missingVar.Name, missingVar.DefaultValue)
		filledVariables = append(filledVariables, answer)
	}

	return filledVariables
}

type VarSet struct {
	Name         string
	DefaultValue string
	UserValue    string
}

func (ev VarSet) value() string {
	if ev.UserValue != "" {
		return ev.UserValue
	}

	return ev.DefaultValue
}

type VarSets []VarSet

func (evs VarSets) ToOpenapi() []openapi.VariableSetValue {
	vars := make([]openapi.VariableSetValue, len(evs))
	for i, ev := range evs {
		vars[i] = openapi.VariableSetValue{
			Key:   openapi.PtrString(ev.Name),
			Value: openapi.PtrString(ev.value()),
		}
	}

	return vars
}

func (evs VarSets) Unique() VarSets {
	seen := make(map[string]bool)
	vars := make(VarSets, 0, len(evs))
	for _, ev := range evs {
		if seen[ev.Name] {
			continue
		}

		seen[ev.Name] = true
		vars = append(vars, ev)
	}

	return vars
}

type MissingVarsError VarSets

func (e MissingVarsError) Error() string {
	return fmt.Sprintf("missing variables: %v", []VarSet(e))
}

func (e MissingVarsError) Is(target error) bool {
	_, ok := target.(MissingVarsError)
	return ok
}

func BuildMissingVarsError(body []byte) error {
	var missingVarsErrResp openapi.MissingVariablesError
	err := jsonFormat.Unmarshal(body, &missingVarsErrResp)
	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %w", err)
	}

	missingVars := VarSets{}

	for _, missingVarErr := range missingVarsErrResp.MissingVariables {
		for _, missingVar := range missingVarErr.Variables {
			missingVars = append(missingVars, VarSet{
				Name:         missingVar.GetKey(),
				DefaultValue: missingVar.GetDefaultValue(),
			})
		}
	}

	return MissingVarsError(missingVars.Unique())
}
