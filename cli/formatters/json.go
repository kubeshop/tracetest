package formatters

import (
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/cli/file"
)

type Json struct {
	toStructFn ToStruct
}

var _ FormatterInterface = Json{}

func NewJson(toStructFn ToStruct) Json {
	return Json{toStructFn}
}

func (j Json) Type() string {
	return "json"
}

func (j Json) Format(file *file.File) (string, error) {
	data, err := j.toStructFn(file)
	if err != nil {
		return "", fmt.Errorf("could not convert file to struct: %w", err)
	}

	bytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", fmt.Errorf("could not marshal output json: %w", err)
	}

	return string(bytes), nil
}
