package validation

import (
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

func ValidateMissingValidations(test model.Test, environment model.Environment) (openapi.MissingVariablesError, error) {
	availableTestVariables := getAvailableVariablesFromTest(test, environment)
	expectedVariables := getRequiredVariablesFromTest(test)
}

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

func getRequiredVariablesFromTest(test model.Test) []string {

}
