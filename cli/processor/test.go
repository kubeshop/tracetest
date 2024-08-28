package processor

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/processor/trigger_preprocessor"
	"go.uber.org/zap"
)

type test struct {
	logger                   *zap.Logger
	applyPollingProfileFunc  applyResourceFunc
	triggerProcessorRegistry trigger_preprocessor.Registry
}

func Test(logger *zap.Logger, triggerProcessorRegistry trigger_preprocessor.Registry, applyPollingProfileFunc applyResourceFunc) test {
	return test{
		logger:                   cmdutil.GetLogger(),
		applyPollingProfileFunc:  applyPollingProfileFunc,
		triggerProcessorRegistry: triggerProcessorRegistry,
	}
}

func (t test) Preprocess(ctx context.Context, input fileutil.File) (fileutil.File, error) {
	var test openapi.TestResource
	err := yaml.Unmarshal(input.Contents(), &test)
	if err != nil {
		t.logger.Error("error parsing test", zap.String("content", string(input.Contents())), zap.Error(err))
		return input, fmt.Errorf("could not unmarshal test yaml: %w", err)
	}

	test, err = t.mapPollingProfiles(ctx, input, test)
	if err != nil {
		t.logger.Error("error mapping polling profiles from test", zap.String("content", string(input.Contents())), zap.Error(err))
		return input, fmt.Errorf("could not map polling profiles referenced in test yaml: %w", err)
	}

	test, err = t.triggerProcessorRegistry.Preprocess(input, test)
	if err != nil {
		return input, fmt.Errorf("could not preprocess trigger: %w", err)
	}

	marshalled, err := yaml.Marshal(test)
	if err != nil {
		return input, fmt.Errorf("could not marshal test yaml: %w", err)
	}

	return fileutil.New(input.AbsPath(), marshalled), nil
}

func (t test) mapPollingProfiles(ctx context.Context, input fileutil.File, test openapi.TestResource) (openapi.TestResource, error) {
	if test.Spec.PollingProfile == nil {
		return test, nil
	}

	pollingProfilePath := test.Spec.PollingProfile
	if !fileutil.LooksLikeFilePath(*pollingProfilePath) {
		t.logger.Debug("does not look like a file path",
			zap.String("path", *pollingProfilePath),
		)

		return test, nil
	}

	f, err := fileutil.Read(input.RelativeFile(*pollingProfilePath))
	if err != nil {
		return openapi.TestResource{}, fmt.Errorf("cannot read polling profile file: %w", err)
	}

	pollingProfileFile, err := t.applyPollingProfileFunc(ctx, f)
	if err != nil {
		return openapi.TestResource{}, fmt.Errorf("cannot apply polling profile '%s': %w", *pollingProfilePath, err)
	}

	var pollingProfile openapi.PollingProfile
	err = yaml.Unmarshal(pollingProfileFile.Contents(), &pollingProfile)
	if err != nil {
		return openapi.TestResource{}, fmt.Errorf("cannot unmarshal updated pollingProfile '%s': %w", *pollingProfilePath, err)
	}

	test.Spec.PollingProfile = &pollingProfile.Spec.Id

	return test, nil
}
