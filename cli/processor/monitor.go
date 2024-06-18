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

type monitor struct {
	logger           *zap.Logger
	applyTestFn      applyResourceFunc
	applyTestSuiteFn applyResourceFunc
}

func Monitor(logger *zap.Logger, applyTestSuiteFn, applyTestFn applyResourceFunc) monitor {
	return monitor{
		logger:           cmdutil.GetLogger(),
		applyTestFn:      applyTestFn,
		applyTestSuiteFn: applyTestSuiteFn,
	}
}

func (m monitor) Preprocess(ctx context.Context, input fileutil.File) (fileutil.File, error) {
	var monitor openapi.MonitorResource
	err := yaml.Unmarshal(input.Contents(), &monitor)
	if err != nil {
		m.logger.Error("error parsing monitor", zap.String("content", string(input.Contents())), zap.Error(err))
		return input, fmt.Errorf("could not unmarshal monitor yaml: %w", err)
	}

	monitor, err = m.mapMonitorTests(ctx, input, monitor)
	if err != nil {
		return input, fmt.Errorf("could not map monitor tests: %w", err)
	}

	monitor, err = m.mapMonitorTestSuites(ctx, input, monitor)
	if err != nil {
		return input, fmt.Errorf("could not map monitor test suites: %w", err)
	}

	marshalled, err := yaml.Marshal(monitor)
	if err != nil {
		return input, fmt.Errorf("could not marshal monitor yaml: %w", err)
	}
	return fileutil.New(input.AbsPath(), marshalled), nil
}

func (m monitor) mapMonitorTests(ctx context.Context, input fileutil.File, monitor openapi.MonitorResource) (openapi.MonitorResource, error) {
	for i, test := range monitor.Spec.GetTests() {
		m.logger.Debug("mapping monitor test",
			zap.Int("index", i),
			zap.String("test", test),
		)
		if !fileutil.LooksLikeFilePath(test) {
			m.logger.Debug("does not look like a file path",
				zap.Int("index", i),
				zap.String("test", test),
			)
			continue
		}

		f, err := fileutil.Read(input.RelativeFile(test))
		if err != nil {
			return openapi.MonitorResource{}, fmt.Errorf("cannot read test file: %w", err)
		}

		testFile, err := m.applyTestFn(ctx, f)
		if err != nil {
			return openapi.MonitorResource{}, fmt.Errorf("cannot apply test '%s': %w", test, err)
		}

		var test openapi.TestResource
		err = yaml.Unmarshal(testFile.Contents(), &test)
		if err != nil {
			return openapi.MonitorResource{}, fmt.Errorf("cannot unmarshal updated test '%s': %w", test, err)
		}

		m.logger.Debug("mapped monitor test",
			zap.Int("index", i),
			zap.String("mapped test", test.Spec.GetId()),
		)

		monitor.Spec.Tests[i] = test.Spec.GetId()
	}

	return monitor, nil
}

func (m monitor) mapMonitorTestSuites(ctx context.Context, input fileutil.File, monitor openapi.MonitorResource) (openapi.MonitorResource, error) {
	for i, suite := range monitor.Spec.GetTestSuites() {
		m.logger.Debug("mapping monitor test suites",
			zap.Int("index", i),
			zap.String("suite", suite),
		)
		if !fileutil.LooksLikeFilePath(suite) {
			m.logger.Debug("does not look like a file path",
				zap.Int("index", i),
				zap.String("suite", suite),
			)
			continue
		}

		f, err := fileutil.Read(input.RelativeFile(suite))
		if err != nil {
			return openapi.MonitorResource{}, fmt.Errorf("cannot read suite file: %w", err)
		}

		suiteFile, err := m.applyTestSuiteFn(ctx, f)
		if err != nil {
			return openapi.MonitorResource{}, fmt.Errorf("cannot apply suite '%s': %w", suite, err)
		}

		var suite openapi.TestSuiteResource
		err = yaml.Unmarshal(suiteFile.Contents(), &suite)
		if err != nil {
			return openapi.MonitorResource{}, fmt.Errorf("cannot unmarshal updated suite '%s': %w", suite, err)
		}

		m.logger.Debug("mapped monitor suite",
			zap.Int("index", i),
			zap.String("mapped suite", suite.Spec.GetId()),
		)

		monitor.Spec.TestSuites[i] = suite.Spec.GetId()
	}

	return monitor, nil
}
