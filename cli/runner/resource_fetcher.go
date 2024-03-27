package runner

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/variable"
	"go.uber.org/zap"
)

type ResourceFetcher interface {
	FetchWithDefinitionFile(context.Context, string) (any, error)
	FetchWithID(context.Context, string) (any, error)
}

type fetcher struct {
	logger *zap.Logger

	resourceRunner Runner
}

func GetResourceFetcher(logger *zap.Logger, resourceRunner Runner) ResourceFetcher {
	return &fetcher{
		logger:         logger,
		resourceRunner: resourceRunner,
	}
}

var _ ResourceFetcher = &fetcher{}

func (f *fetcher) FetchWithDefinitionFile(ctx context.Context, definitionFile string) (any, error) {
	file, err := fileutil.Read(definitionFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read definition file %s: %w", definitionFile, err)
	}
	df := file
	f.logger.Debug("Definition file read", zap.String("absolutePath", df.AbsPath()))

	df, err = f.injectLocalEnvVars(df)
	if err != nil {
		return nil, fmt.Errorf("cannot inject local env vars: %w", err)
	}

	resource, err := f.resourceRunner.Apply(ctx, df)
	if err != nil {
		return nil, fmt.Errorf("cannot apply definition file: %w", err)
	}
	f.logger.Debug("Definition file applied", zap.String("updated", string(df.Contents())))

	return resource, nil
}

func (f *fetcher) FetchWithID(ctx context.Context, resourceID string) (any, error) {
	f.logger.Debug("Definition file not provided, fetching resource by ID", zap.String("ID", resourceID))

	resource, err := f.resourceRunner.GetByID(ctx, resourceID)
	if err != nil {
		return nil, fmt.Errorf("cannot get resource by ID: %w", err)
	}
	f.logger.Debug("Resource fetched by ID", zap.String("ID", resourceID), zap.Any("resource", resource))

	return resource, nil
}

func (f *fetcher) injectLocalEnvVars(df fileutil.File) (fileutil.File, error) {
	variableInjector := variable.NewInjector(variable.WithVariableProvider(
		variable.EnvironmentVariableProvider{},
	))

	injected, err := variableInjector.ReplaceInString(string(df.Contents()))
	if err != nil {
		return df, fmt.Errorf("cannot inject local variable set: %w", err)
	}

	df = fileutil.New(df.AbsPath(), []byte(injected))

	return df, nil
}
