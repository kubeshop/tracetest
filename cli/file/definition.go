package file

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
)

func LoadDefinition(file string) (definition.File, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return definition.File{}, fmt.Errorf("could not read definition file %s: %w", file, err)
	}

	f, err := definition.Decode(string(b))
	if err != nil {
		return definition.File{}, fmt.Errorf("could not parse definition file: %w", err)
	}

	return f, nil
}

func SaveDefinition(file, definition string) error {
	err := os.WriteFile(file, []byte(definition), 0644)
	if err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}

func SetTestID(file, id string) error {
	idStatementRegex := regexp.MustCompile("^id: [0-9a-zA-Z\\-_]+\n")
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

	return SaveDefinition(file, fileContent)
}
