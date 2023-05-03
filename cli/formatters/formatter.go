package formatters

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/file"
)

type ToStruct func(*file.File) (interface{}, error)
type ToListStruct func(*file.File) (interface{}, error)

type ToTable func(*file.File) (*simpletable.Header, *simpletable.Body, error)
type ToListTable ToTable

type ResourceFormatter interface {
	ToTable(*file.File) (*simpletable.Header, *simpletable.Body, error)
	ToListTable(*file.File) (*simpletable.Header, *simpletable.Body, error)
	ToStruct(*file.File) (interface{}, error)
	ToListStruct(*file.File) (interface{}, error)
}

type FormatterInterface interface {
	Format(*file.File) (string, error)
	Type() string
}

type Formatter struct {
	formatType string
	registry   map[string]FormatterInterface
}

func NewFormatter(formatType string, formatters ...FormatterInterface) Formatter {
	registry := make(map[string]FormatterInterface, len(formatters))

	for _, option := range formatters {
		registry[option.Type()] = option
	}

	return Formatter{formatType, registry}
}

func (f Formatter) Format(file *file.File) (string, error) {
	formatter, ok := f.registry[f.formatType]
	if !ok {
		return "", fmt.Errorf("formatter %s not found", f.formatType)
	}

	return formatter.Format(file)
}

func BuildFormatter(formatType string, defaultType Output, toTable ToTable, toStruct ToStruct) Formatter {
	jsonFormatter := NewJson(toStruct)
	yamlFormatter := NewYaml(toStruct)
	tableFormatter := NewTable(toTable)

	if defaultType == "" {
		defaultType = YAML
	}

	if formatType == "" {
		formatType = string(defaultType)
	}

	return NewFormatter(formatType, jsonFormatter, yamlFormatter, tableFormatter)
}
