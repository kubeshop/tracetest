package file

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kubeshop/tracetest/cli/definition"
	"gopkg.in/yaml.v2"
)

func LoadDefinition(file string) (definition.Test, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return definition.Test{}, fmt.Errorf("could not read test definition file %s: %w", file, err)
	}

	fileBytes, err = injectEnvVariables(fileBytes)
	if err != nil {
		return definition.Test{}, fmt.Errorf("could not inject env variables into definition file: %w", err)
	}

	test := definition.Test{}
	err = yaml.Unmarshal(fileBytes, &test)
	if err != nil {
		return definition.Test{}, fmt.Errorf("could not parse test definition file: %w", err)
	}

	return test, nil
}

func SaveDefinition(file string, definition definition.Test) error {
	yamlContent, err := yaml.Marshal(definition)
	if err != nil {
		return fmt.Errorf("could not marshal definition into YAML: %w", err)
	}

	err = os.WriteFile(file, yamlContent, 0755)
	if err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}

func injectEnvVariables(fileBytes []byte) ([]byte, error) {
	envVarRegex, err := regexp.Compile(`\$\{\w+\}`)
	if err != nil {
		return []byte{}, fmt.Errorf("could not compile env variable regex: %w", err)
	}

	fileString := string(fileBytes)
	allEnvVariables := envVarRegex.FindAllString(fileString, -1)

	for _, envVariableExpression := range allEnvVariables {
		envVarName := envVariableExpression[2 : len(envVariableExpression)-1] // removes '${' and '}'
		envVarValue := os.Getenv(envVarName)

		fileString = strings.Replace(fileString, envVariableExpression, envVarValue, -1)
	}

	return []byte(fileString), nil
}
