package runner

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type VariableSetFetcher interface {
	Fetch(context.Context, string) (string, error)
}

type internalVariableSetFetcher struct {
	logger *zap.Logger

	variableSetClient resourcemanager.Client
}

func GetVariableSetFetcher(logger *zap.Logger, variableSetClient resourcemanager.Client) VariableSetFetcher {
	return &internalVariableSetFetcher{
		logger:            logger,
		variableSetClient: variableSetClient,
	}
}

var _ VariableSetFetcher = &internalVariableSetFetcher{}

func (f *internalVariableSetFetcher) Fetch(ctx context.Context, varsID string) (string, error) {
	if varsID == "" {
		return "", nil // user have not defined variables, skipping it
	}

	if !fileutil.IsFilePath(varsID) {
		f.logger.Debug("varsID is not a file path", zap.String("vars", varsID))

		// validate that env exists
		_, err := f.variableSetClient.Get(ctx, varsID, resourcemanager.Formats.Get(resourcemanager.FormatYAML))
		if errors.Is(err, resourcemanager.ErrNotFound) {
			return "", fmt.Errorf("variable set '%s' not found", varsID)
		}
		if err != nil {
			return "", fmt.Errorf("cannot get variable set '%s': %w", varsID, err)
		}

		f.logger.Debug("envID is valid")

		return varsID, nil
	}

	file, err := fileutil.Read(varsID)
	if err != nil {
		return "", fmt.Errorf("cannot read environment set file %s: %w", varsID, err)
	}

	f.logger.Debug("envID is a file path", zap.String("filePath", varsID), zap.Any("file", f))
	updatedEnv, err := f.variableSetClient.Apply(ctx, file, yamlFormat)
	if err != nil {
		return "", fmt.Errorf("could not read environment set file: %w", err)
	}

	var vars openapi.VariableSetResource
	err = yaml.Unmarshal([]byte(updatedEnv), &vars)
	if err != nil {
		f.logger.Error("error parsing json", zap.String("content", updatedEnv), zap.Error(err))
		return "", fmt.Errorf("could not unmarshal variable set json: %w", err)
	}

	return vars.Spec.GetId(), nil
}
