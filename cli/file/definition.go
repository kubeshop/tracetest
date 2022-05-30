package file

import (
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/definition"
	"gopkg.in/yaml.v2"
)

func LoadDefinition(file string) (definition.Test, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return definition.Test{}, fmt.Errorf("could not read file %s: %w", file, err)
	}

	test := definition.Test{}
	err = yaml.Unmarshal(fileBytes, &test)
	if err != nil {
		return definition.Test{}, fmt.Errorf("could not parse definition file: %w", err)
	}

	return test, nil
}
