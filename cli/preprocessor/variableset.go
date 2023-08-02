package preprocessor

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"go.uber.org/zap"
)

type variableSet struct {
	logger *zap.Logger
}

type generic struct {
	Type string      `yaml:"type"`
	Spec interface{} `yaml:"spec"`
}

func VariableSet(logger *zap.Logger) variableSet {
	return variableSet{
		logger: logger,
	}
}

func (vs variableSet) Preprocess(ctx context.Context, input fileutil.File) (fileutil.File, error) {
	var resource generic
	err := yaml.Unmarshal(input.Contents(), &resource)
	if err != nil {
		vs.logger.Error("error parsing test", zap.String("content", string(input.Contents())), zap.Error(err))
		return input, fmt.Errorf("could not unmarshal test yaml: %w", err)
	}

	if resource.Type == "Environment" {
		resource.Type = "VariableSet"
	}

	marshalled, err := yaml.Marshal(resource)
	if err != nil {
		return input, fmt.Errorf("could not marshal test yaml: %w", err)
	}

	return fileutil.New(input.AbsPath(), marshalled), nil
}
