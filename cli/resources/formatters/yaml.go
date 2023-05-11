package resources_formatters

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/kubeshop/tracetest/cli/file"
)

type Yaml struct {
	toStructFn     ToStruct
	toListStructFn ToListStruct
}

var _ FormatterInterface = Yaml{}

func NewYaml(resourceFormatter ResourceFormatter) Yaml {
	return Yaml{
		toStructFn:     resourceFormatter.ToStruct,
		toListStructFn: resourceFormatter.ToListStruct,
	}
}

func (Yaml) Type() string {
	return "yaml"
}

func (y Yaml) FormatList(file *file.File) (string, error) {
	data, err := y.toListStructFn(file)
	if err != nil {
		return "", fmt.Errorf("could not convert file to struct: %w", err)
	}

	result := ""
	for _, value := range data {
		bytes, err := yaml.Marshal(value)
		if err != nil {
			return "", fmt.Errorf("could not marshal output json: %w", err)
		}

		result += "---\n" + string(bytes) + "\n"
	}

	return result, nil
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

	return "---\n" + string(bytes), nil
}
