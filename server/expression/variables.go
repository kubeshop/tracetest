package expression

import (
	"github.com/kubeshop/tracetest/server/model"
)

type TestVariables struct {
	TestId      string
	Environment []string
	Variables   []string
	Missing     []MissingVariables
}

type Variables struct {
	Environment model.Environment
	Executor    Executor
}

type MissingVariables struct {
	Key          string
	DefaultValue string
}

type VariablesMap map[string]bool

func (vm VariablesMap) Merge(other VariablesMap) VariablesMap {
	for key := range other {
		vm[key] = true
	}

	return vm
}

func (vm VariablesMap) MergeStringArray(strArray []string) VariablesMap {
	for _, value := range strArray {
		vm[value] = true
	}

	return vm
}

func (vm VariablesMap) ToArray() []string {
	array := []string{}

	for key := range vm {
		array = append(array, key)
	}

	return array
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

func (v Variables) GetEnvironmentVariables(test model.Test) VariablesMap {
	envVariables := VariablesMap{}

	for _, envVar := range v.Environment.Values {
		envVariables[envVar.Key] = true
	}

	test.Outputs.ForEach(func(key string, _ model.Output) error {
		envVariables[key] = true

		return nil
	})

	return envVariables
}

func (v Variables) GetSpecsVariables(test model.Test) (VariablesMap, error) {
	specVariables := VariablesMap{}
	err := test.Specs.ForEach(func(_ model.SpanQuery, namedAssertions model.NamedAssertions) error {
		for _, assertion := range namedAssertions.Assertions {
			variables, err := v.Executor.StatementTermsByType(string(assertion), EnvironmentType)
			if err != nil {
				return err
			}

			specVariables = specVariables.MergeStringArray(variables)
		}

		return nil
	})

	return specVariables, err
}

func (v Variables) GetOutputVariables(test model.Test) (VariablesMap, error) {
	specVariables := VariablesMap{}
	err := test.Outputs.ForEach(func(_ string, output model.Output) error {
		variables, err := v.Executor.StatementTermsByType(string(output.Value), EnvironmentType)
		if err != nil {
			return err
		}

		specVariables = specVariables.MergeStringArray(variables)

		return nil
	})

	return specVariables, err
}

func (v Variables) getMissingVariables(testVariables, environmentVariables VariablesMap, environment model.Environment) []MissingVariables {
	missingVariables := []MissingVariables{}

	for testVariableKey := range testVariables {
		exists := environmentVariables[testVariableKey]
		if !exists {
			missingVariables = append(missingVariables, MissingVariables{
				Key:          testVariableKey,
				DefaultValue: environment.Get(testVariableKey),
			})
		}
	}

	return missingVariables
}

func (v Variables) GetTestVariables(testId string, environmentVariables, testVariables VariablesMap, environment model.Environment) TestVariables {
	variablesResult := TestVariables{
		TestId:      testId,
		Environment: environmentVariables.ToArray(),
		Variables:   testVariables.ToArray(),
		Missing:     v.getMissingVariables(testVariables, environmentVariables, environment),
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
