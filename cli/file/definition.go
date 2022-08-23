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

	err = os.WriteFile(file, yamlContent, 0644)
	if err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}

func SetTestID(file string, id string) error {
	idStatementRegex := regexp.MustCompile("^id: [0-9a-zA-Z\\-]+\n")
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("could not read test definition file %s: %w", file, err)
	}

	fileContent := string(fileBytes)
	idStatement := idStatementRegex.FindString(fileContent)
	if idStatement != "" {
		fileContent = strings.Replace(fileContent, idStatement, fmt.Sprintf("id: %s\n", id), 1)
	} else {
		fileContent = fmt.Sprintf("id: %s\n%s", id, fileContent)
	}

	err = os.WriteFile(file, []byte(fileContent), 0644)
	if err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}
