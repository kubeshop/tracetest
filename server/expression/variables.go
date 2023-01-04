package expression

import (
	"github.com/kubeshop/tracetest/server/model"
)

type TestVariables struct {
	TestId      string
	Environment []string
	Variables   []string
	Missing     []string
}

type Variables struct {
	Environment model.Environment
	Executor    Executor
}

func NewVariables(environment model.Environment, executor Executor) Variables {
	return Variables{
		Environment: environment,
		Executor:    executor,
	}
}

func (v Variables) UpdateEnvironmentVariables(test model.Test) Variables {
	newEnv := v.mergeOutputsIntoEnv(v.Environment, test.Outputs)

	return Variables{
		Environment: newEnv,
		Executor:    v.Executor,
	}
}

func (v Variables) GetEnvironmentVariables(test model.Test) []string {
	envVariables := []string{}

	for _, envVar := range v.Environment.Values {
		envVariables = append(envVariables, envVar.Key)
	}

	test.Outputs.ForEach(func(key string, _ model.Output) error {
		envVariables = append(envVariables, key)

		return nil
	})

	return envVariables
}

func (v Variables) GetSpecsVariables(test model.Test) ([]string, error) {
	specVariables := []string{}
	err := test.Specs.ForEach(func(_ model.SpanQuery, namedAssertions model.NamedAssertions) error {
		for _, assertion := range namedAssertions.Assertions {
			variables, err := v.Executor.StatementTermsByType(string(assertion), EnvironmentType)
			if err != nil {
				return err
			}

			specVariables = append(specVariables, variables...)
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return specVariables, nil
}

func (v Variables) GetOutputVariables(test model.Test) ([]string, error) {
	specVariables := []string{}
	err := test.Outputs.ForEach(func(_ string, output model.Output) error {
		variables, err := v.Executor.StatementTermsByType(string(output.Value), EnvironmentType)
		if err != nil {
			return err
		}

		specVariables = append(specVariables, variables...)

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return specVariables, nil
}

func contains(slice []string, value string) bool {
	for _, a := range slice {
		if a == value {
			return true
		}
	}
	return false
}

func (v Variables) getMissingVariables(testVariables, environmentVariables []string) []string {
	missingVariables := []string{}

	for _, envVar := range testVariables {
		if !contains(environmentVariables, envVar) {
			missingVariables = append(missingVariables, envVar)
		}
	}

	return missingVariables
}

func (v Variables) GetTestVariables(testId string, environmentVariables, testVariables []string) TestVariables {
	variablesResult := TestVariables{
		TestId:      testId,
		Environment: environmentVariables,
		Variables:   testVariables,
		Missing:     v.getMissingVariables(testVariables, environmentVariables),
	}

	return variablesResult
}

func (v Variables) mergeOutputsIntoEnv(env model.Environment, outputs model.OrderedMap[string, model.Output]) model.Environment {
	newEnv := make([]model.EnvironmentValue, 0, outputs.Len())
	outputs.ForEach(func(key string, _ model.Output) error {
		newEnv = append(newEnv, model.EnvironmentValue{
			Key:   key,
			Value: "",
		})

		return nil
	})

	return env.Merge(model.Environment{
		Values: newEnv,
	})
}
