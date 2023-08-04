package validation

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kubeshop/tracetest/server/expression/linting"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/transaction"
	"github.com/kubeshop/tracetest/server/variableset"
)

var ErrMissingVariables = errors.New("variables are missing")

type augmentedTestGetter interface {
	GetAugmented(context.Context, id.ID) (test.Test, error)
}

func ValidateMissingVariables(ctx context.Context, testRepo augmentedTestGetter, runRepo test.RunRepository, test test.Test, env variableset.VariableSet) (openapi.MissingVariablesError, error) {
	missingVariables := getMissingVariables(test, env)
	previousValues := map[string]variableset.VariableSetValue{}
	var err error
	if len(missingVariables) > 0 {
		previousValues, err = getPreviousEnvironmentValues(ctx, testRepo, runRepo, test)
		if err != nil {
			return openapi.MissingVariablesError{}, err
		}
	}
	return buildErrorObject(test, missingVariables, previousValues)
}

func getMissingVariables(test test.Test, environment variableset.VariableSet) []string {
	availableTestVariables := getAvailableVariables(test, environment)
	expectedVariables := linting.DetectMissingVariables(test, availableTestVariables)

	availableVariablesMap := make(map[string]bool, 0)
	for _, availableVariable := range availableTestVariables {
		availableVariablesMap[availableVariable] = true
	}

	missingVariables := []string{}

	for _, expectedVariable := range expectedVariables {
		if _, exists := availableVariablesMap[expectedVariable]; !exists {
			missingVariables = append(missingVariables, expectedVariable)
		}
	}

	return missingVariables
}

func getAvailableVariables(test test.Test, environment variableset.VariableSet) []string {
	availableVariables := make([]string, 0)
	for _, env := range environment.Values {
		availableVariables = append(availableVariables, env.Key)
	}

	for _, output := range test.Outputs {
		availableVariables = append(availableVariables, output.Name)
	}

	return availableVariables
}

func getPreviousEnvironmentValues(ctx context.Context, testRepo augmentedTestGetter, runRepo test.RunRepository, test test.Test) (map[string]variableset.VariableSetValue, error) {
	latestTestVersion, err := testRepo.GetAugmented(ctx, test.ID)
	if err != nil {
		return map[string]variableset.VariableSetValue{}, err
	}

	previousTestRun, err := runRepo.GetLatestRunByTestVersion(ctx, test.ID, *latestTestVersion.Version)
	if err != nil {
		// If error is not found, it means this is the first run. So just ignore this error
		// and provide empty values in the default values for the missing variables
		if err != sql.ErrNoRows {
			return map[string]variableset.VariableSetValue{}, err
		}
	} else {
		envMap := make(map[string]variableset.VariableSetValue, len(previousTestRun.VariableSet.Values))
		for _, envVar := range previousTestRun.VariableSet.Values {
			envMap[envVar.Key] = envVar
		}

		return envMap, nil
	}

	return map[string]variableset.VariableSetValue{}, nil
}

func ValidateMissingVariablesFromTransaction(ctx context.Context, testRepo augmentedTestGetter, runRepo test.RunRepository, transaction transaction.Transaction, env variableset.VariableSet) (openapi.MissingVariablesError, error) {
	missingVariables := make([]openapi.MissingVariable, 0)
	for _, step := range transaction.Steps {
		stepValidationResult, err := ValidateMissingVariables(ctx, testRepo, runRepo, step, env)
		if err != ErrMissingVariables {
			return openapi.MissingVariablesError{}, err
		}

		missingVariables = append(missingVariables, stepValidationResult.MissingVariables...)

		// update env with this test outputs
		outputs := make([]variableset.VariableSetValue, 0)
		for _, output := range step.Outputs {
			outputs = append(outputs, variableset.VariableSetValue{Key: output.Name})
		}

		env.Values = append(env.Values, outputs...)
	}

	if len(missingVariables) > 0 {
		return openapi.MissingVariablesError{MissingVariables: missingVariables}, ErrMissingVariables
	}

	return openapi.MissingVariablesError{}, nil
}

func buildErrorObject(test test.Test, missingVariables []string, previousValues map[string]variableset.VariableSetValue) (openapi.MissingVariablesError, error) {
	if len(missingVariables) > 0 {
		missingVariableObjects := make([]openapi.Variable, 0, len(missingVariables))
		for _, variable := range missingVariables {

			missingVariableObjects = append(missingVariableObjects, openapi.Variable{
				Key:          variable,
				DefaultValue: "",
			})
		}

		missingVariablesError := openapi.MissingVariablesError{
			MissingVariables: []openapi.MissingVariable{
				{TestId: string(test.ID), Variables: missingVariableObjects},
			},
		}

		for i, missingVariables := range missingVariablesError.MissingVariables {
			for j, missingVariable := range missingVariables.Variables {
				if envVar, found := previousValues[missingVariable.Key]; found {
					missingVariablesError.MissingVariables[i].Variables[j].DefaultValue = envVar.Value
				}
			}
		}

		return missingVariablesError, ErrMissingVariables
	}

	return openapi.MissingVariablesError{}, nil
}
