package validation

import (
	"context"
	"errors"

	"github.com/kubeshop/tracetest/server/expression/linting"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/testdb"
)

var ErrMissingVariables = errors.New("variables are missing")

func ValidateMissingVariables(ctx context.Context, db model.Repository, test model.Test, environment model.Environment) (openapi.MissingVariablesError, error) {
	missingVariables := getMissingVariables(test, environment)
	return buildErrorObject(ctx, db, missingVariables, test)
}

func getMissingVariables(test model.Test, environment model.Environment) []string {
	availableTestVariables := getAvailableVariablesFromTest(test, environment)
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

// func ValidateMissingVariablesFromTransaction(ctx context.Context, db model.Repository, transaction model.Transaction, environment model.Environment) (openapi.MissingVariablesError, error) {
// 	missingVariables := getMissingVariablesFromTransaction(transaction, environment)
// 	return buildErrorObject(ctx, db, missingVariables)
// }

func buildErrorObject(ctx context.Context, db model.Repository, missingVariables []string, test model.Test) (openapi.MissingVariablesError, error) {
	if len(missingVariables) > 0 {
		missingVariablesError := openapi.MissingVariablesError{
			MissingVariables: make([]openapi.MissingVariable, 0, len(missingVariables)),
		}

		for _, variable := range missingVariables {
			missingVariablesError.MissingVariables = append(missingVariablesError.MissingVariables, openapi.MissingVariable{
				Key:          variable,
				DefaultValue: "",
			})
		}

		latestTestVersion, err := db.GetLatestTestVersion(ctx, test.ID)
		if err != nil {
			return openapi.MissingVariablesError{}, err
		}

		previousTestRun, err := db.GetLatestRunByTestVersion(ctx, test.ID, latestTestVersion.Version)
		if err != nil {
			// If error is not found, it means this is the first run. So just ignore this error
			// and provide empty values in the default values for the missing variables
			if err != testdb.ErrNotFound {
				return openapi.MissingVariablesError{}, err
			}
		} else {
			envMap := make(map[string]model.EnvironmentValue, len(previousTestRun.Environment.Values))
			for _, envVar := range previousTestRun.Environment.Values {
				envMap[envVar.Key] = envVar
			}

			for i, missingVariable := range missingVariablesError.MissingVariables {
				if envVar, found := envMap[missingVariable.Key]; found {
					missingVariablesError.MissingVariables[i].DefaultValue = envVar.Value
				}
			}
		}

		return missingVariablesError, ErrMissingVariables
	}

	return openapi.MissingVariablesError{}, nil
}

// func getMissingVariablesFromTransaction(transaction model.Transaction, environment model.Environment) []string {
// 	// TODO
// 	return []string{}
// }

func getAvailableVariablesFromTest(test model.Test, environment model.Environment) []string {
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
