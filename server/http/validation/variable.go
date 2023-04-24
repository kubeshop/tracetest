package validation

import (
	"context"
	"errors"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/expression/linting"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/tests"
)

var ErrMissingVariables = errors.New("variables are missing")

func ValidateMissingVariables(ctx context.Context, db model.Repository, test model.Test, env environment.Environment) (openapi.MissingVariablesError, error) {
	missingVariables := getMissingVariables(test, env)
	previousValues := map[string]environment.EnvironmentValue{}
	var err error
	if len(missingVariables) > 0 {
		previousValues, err = getPreviousEnvironmentValues(ctx, db, test)
		if err != nil {
			return openapi.MissingVariablesError{}, err
		}
	}
	return buildErrorObject(test, missingVariables, previousValues)
}

func getMissingVariables(test model.Test, environment environment.Environment) []string {
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

func getAvailableVariables(test model.Test, environment environment.Environment) []string {
	availableVariables := make([]string, 0)
	for _, env := range environment.Values {
		availableVariables = append(availableVariables, env.Key)
	}

	test.Outputs.ForEach(func(key string, _ model.Output) error {
		availableVariables = append(availableVariables, key)
		return nil
	})

	return availableVariables
}

func getPreviousEnvironmentValues(ctx context.Context, db model.Repository, test model.Test) (map[string]environment.EnvironmentValue, error) {
	latestTestVersion, err := db.GetLatestTestVersion(ctx, test.ID)
	if err != nil {
		return map[string]environment.EnvironmentValue{}, err
	}

	previousTestRun, err := db.GetLatestRunByTestVersion(ctx, test.ID, latestTestVersion.Version)
	if err != nil {
		// If error is not found, it means this is the first run. So just ignore this error
		// and provide empty values in the default values for the missing variables
		if err != testdb.ErrNotFound {
			return map[string]environment.EnvironmentValue{}, err
		}
	} else {
		envMap := make(map[string]environment.EnvironmentValue, len(previousTestRun.Environment.Values))
		for _, envVar := range previousTestRun.Environment.Values {
			envMap[envVar.Key] = envVar
		}

		return envMap, nil
	}

	return map[string]environment.EnvironmentValue{}, nil
}

func ValidateMissingVariablesFromTransaction(ctx context.Context, db model.Repository, transaction tests.Transaction, environment environment.Environment) (openapi.MissingVariablesError, error) {	missingVariables := make([]openapi.MissingVariable, 0)
	for _, step := range transaction.Steps {
		stepValidationResult, err := ValidateMissingVariables(ctx, db, step, env)
		if err != ErrMissingVariables {
			return openapi.MissingVariablesError{}, err
		}

		missingVariables = append(missingVariables, stepValidationResult.MissingVariables...)

		// update env with this test outputs
		outputs := make([]environment.EnvironmentValue, 0)
		step.Outputs.ForEach(func(key string, val model.Output) error {
			outputs = append(outputs, environment.EnvironmentValue{Key: key})
			return nil
		})

		env.Values = append(env.Values, outputs...)
	}

	if len(missingVariables) > 0 {
		return openapi.MissingVariablesError{MissingVariables: missingVariables}, ErrMissingVariables
	}

	return openapi.MissingVariablesError{}, nil
}

func buildErrorObject(test model.Test, missingVariables []string, previousValues map[string]environment.EnvironmentValue) (openapi.MissingVariablesError, error) {
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
