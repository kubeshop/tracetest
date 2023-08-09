package preprocessor

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"go.uber.org/zap"
)

type applyTestFunc func(context.Context, fileutil.File) (fileutil.File, error)

type testSuite struct {
	logger      *zap.Logger
	applyTestFn applyTestFunc
}

func TestSuite(logger *zap.Logger, applyTestFn applyTestFunc) testSuite {
	return testSuite{
		logger:      logger,
		applyTestFn: applyTestFn,
	}
}

func (t testSuite) Preprocess(ctx context.Context, input fileutil.File) (fileutil.File, error) {
	var suite openapi.TestSuiteResource
	err := yaml.Unmarshal(input.Contents(), &suite)
	if err != nil {
		t.logger.Error("error parsing test suite", zap.String("content", string(input.Contents())), zap.Error(err))
		return input, fmt.Errorf("could not unmarshal test suite yaml: %w", err)
	}

	// TODO: This should be removed at some point as the Environment type is deprecated
	if *suite.Type == "Transaction" {
		textType := "TestSuite"
		suite.Type = &textType
	}

	suite, err = t.mapTestSuiteSteps(ctx, input, suite)
	if err != nil {
		return input, fmt.Errorf("could not map test suite steps: %w", err)
	}

	marshalled, err := yaml.Marshal(suite)
	if err != nil {
		return input, fmt.Errorf("could not marshal test yaml: %w", err)
	}
	return fileutil.New(input.AbsPath(), marshalled), nil
}

func (t testSuite) mapTestSuiteSteps(ctx context.Context, input fileutil.File, suite openapi.TestSuiteResource) (openapi.TestSuiteResource, error) {
	for i, step := range suite.Spec.GetSteps() {
		t.logger.Debug("mapping test suite step",
			zap.Int("index", i),
			zap.String("step", step),
		)
		if !fileutil.LooksLikeFilePath(step) {
			t.logger.Debug("does not look like a file path",
				zap.Int("index", i),
				zap.String("step", step),
			)
			continue
		}

		f, err := fileutil.Read(input.RelativeFile(step))
		if err != nil {
			return openapi.TestSuiteResource{}, fmt.Errorf("cannot read test file: %w", err)
		}

		testFile, err := t.applyTestFn(ctx, f)
		if err != nil {
			return openapi.TestSuiteResource{}, fmt.Errorf("cannot apply test '%s': %w", step, err)
		}

		var test openapi.TestResource
		err = yaml.Unmarshal(testFile.Contents(), &test)
		if err != nil {
			return openapi.TestSuiteResource{}, fmt.Errorf("cannot unmarshal updated test '%s': %w", step, err)
		}

		t.logger.Debug("mapped test suite step",
			zap.Int("index", i),
			zap.String("step", step),
			zap.String("mapped step", test.Spec.GetId()),
		)

		suite.Spec.Steps[i] = test.Spec.GetId()
	}

	return suite, nil
}
