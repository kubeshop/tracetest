package processor

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"go.uber.org/zap"
)

type applyFn func(ctx context.Context, resource fileutil.File) (fileutil.File, error)
type updateEnvFn func(ctx context.Context, envID string) error

type environment struct {
	logger      *zap.Logger
	applyFn     applyFn
	updateEnvFn updateEnvFn
}

func Environment(logger *zap.Logger, applyFn applyFn, updateEnvFn updateEnvFn) environment {
	return environment{logger: logger, applyFn: applyFn, updateEnvFn: updateEnvFn}
}

func (e environment) Postprocess(ctx context.Context, input fileutil.File) (fileutil.File, error) {
	var env openapi.EnvironmentResource
	err := yaml.Unmarshal(input.Contents(), &env)
	if err != nil {
		e.logger.Error("error parsing test suite", zap.String("content", string(input.Contents())), zap.Error(err))
		return input, fmt.Errorf("could not unmarshal test suite yaml: %w", err)
	}

	if env.GetSpec().Id != nil {
		err = e.updateEnvFn(ctx, *env.GetSpec().Id)
		if err != nil {
			return input, fmt.Errorf("could not update environment: %w", err)
		}
	}

	if env.GetSpec().Resources != nil {
		err = e.mapResources(ctx, input, *env.GetSpec().Resources)
		if err != nil {
			return input, fmt.Errorf("could not map environment resources: %w", err)
		}
	}

	marshalled, err := yaml.Marshal(env)
	if err != nil {
		return input, fmt.Errorf("could not marshal environment yaml: %w", err)
	}

	return fileutil.New(input.AbsPath(), marshalled), nil
}

func (e environment) mapResources(ctx context.Context, input fileutil.File, resources string) error {
	// resources
	if !fileutil.LooksLikeRelativeFilePath(resources) {
		return nil
	}

	files := fileutil.ReadDirFileNames(resources)

	for _, fileName := range files {
		if fileutil.IsDir(fileName) {
			err := e.mapResources(ctx, input, fileName)
			if err != nil {
				return fmt.Errorf("cannot map resources: %w", err)
			}

			continue
		}

		resource, err := fileutil.Read(fileName)
		if err != nil {
			return fmt.Errorf("cannot read resource file: %w", err)
		}

		_, err = e.applyFn(ctx, resource)
		if err != nil {
			return fmt.Errorf("cannot apply resource: %w", err)
		}
	}

	return nil
}
