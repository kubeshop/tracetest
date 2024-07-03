package processor

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"go.uber.org/zap"
)

type test struct {
	logger                  *zap.Logger
	applyPollingProfileFunc applyResourceFunc
}

func Test(logger *zap.Logger, applyPollingProfileFunc applyResourceFunc) test {
	return test{
		logger:                  cmdutil.GetLogger(),
		applyPollingProfileFunc: applyPollingProfileFunc,
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

	test, err = t.consolidateGRPCFile(input, test)
	if err != nil {
		return input, fmt.Errorf("could not consolidate grpc file: %w", err)
	}

	test, err = t.consolidateScriptFile(input, test)
	if err != nil {
		return input, fmt.Errorf("could not consolidate playwrightengine script file: %w", err)
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

func (t test) consolidateGRPCFile(input fileutil.File, test openapi.TestResource) (openapi.TestResource, error) {
	if test.Spec.Trigger.GetType() != "grpc" {
		t.logger.Debug("test does not use grpc", zap.String("triggerType", test.Spec.Trigger.GetType()))
		return test, nil
	}

	definedPBFile := test.Spec.Trigger.Grpc.GetProtobufFile()
	if !t.isValidGrpcFilePath(definedPBFile, input.AbsDir()) {
		t.logger.Debug("protobuf file is not a file path", zap.String("protobufFile", definedPBFile))
		return test, nil
	}

	pbFilePath := input.RelativeFile(definedPBFile)
	t.logger.Debug("protobuf file", zap.String("path", pbFilePath))

	pbFile, err := fileutil.Read(pbFilePath)
	if err != nil {
		return test, fmt.Errorf(`cannot read protobuf file: %w`, err)
	}
	t.logger.Debug("protobuf file contents", zap.String("contents", string(pbFile.Contents())))

	test.Spec.Trigger.Grpc.SetProtobufFile(string(pbFile.Contents()))

	return test, nil
}

func (t test) consolidateScriptFile(input fileutil.File, test openapi.TestResource) (openapi.TestResource, error) {
	if test.Spec.Trigger.GetType() != "playwrightengine" {
		t.logger.Debug("test does not use playwrightengine", zap.String("triggerType", test.Spec.Trigger.GetType()))
		return test, nil
	}

	definedScriptFile := test.Spec.Trigger.PlaywrightEngine.GetScript()
	if !t.isValidGrpcFilePath(definedScriptFile, input.AbsDir()) {
		t.logger.Debug("script file is not a file path", zap.String("protobufFile", definedScriptFile))
		return test, nil
	}

	scriptFilePath := input.RelativeFile(definedScriptFile)
	t.logger.Debug("script file", zap.String("path", scriptFilePath))

	scriptFile, err := fileutil.Read(scriptFilePath)
	if err != nil {
		return test, fmt.Errorf(`cannot read script file: %w`, err)
	}
	t.logger.Debug("script file contents", zap.String("contents", string(scriptFile.Contents())))

	test.Spec.Trigger.PlaywrightEngine.SetScript(string(scriptFile.Contents()))

	return test, nil
}

func (t test) isValidGrpcFilePath(grpcFilePath, testFile string) bool {
	if fileutil.LooksLikeRelativeFilePath(grpcFilePath) {
		// if looks like a relative file path, test if it exists
		return fileutil.IsFilePathToRelativeDir(grpcFilePath, testFile)
	}

	// it could be an absolute file path, test it
	return fileutil.IsFilePath(grpcFilePath)
}
