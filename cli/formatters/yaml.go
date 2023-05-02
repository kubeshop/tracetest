package formatters

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/kubeshop/tracetest/cli/file"
)

type Yaml struct {
	toStructFn ToStruct
}

var _ FormatterInterface = Yaml{}

func NewYaml(toStructFn ToStruct) Yaml {
	return Yaml{toStructFn}
}

func (Yaml) Type() string {
	return "yaml"
}

func (y Yaml) Format(file *file.File) (string, error) {
	data, err := y.toStructFn(file)
	if err != nil {
		return "", fmt.Errorf("could not convert file to struct: %w", err)
	}

	bytes, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("could not marshal output json: %w", err)
	}

	return string(bytes), nil
}
