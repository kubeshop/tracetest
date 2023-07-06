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

type transaction struct {
	logger      *zap.Logger
	applyTestFn applyTestFunc
}

func Transaction(logger *zap.Logger, applyTestFn applyTestFunc) transaction {
	return transaction{
		logger:      logger,
		applyTestFn: applyTestFn,
	}
}

func (t transaction) Preprocess(ctx context.Context, input fileutil.File) (fileutil.File, error) {
	var tran openapi.TransactionResource
	err := yaml.Unmarshal(input.Contents(), &tran)
	if err != nil {
		t.logger.Error("error parsing transaction", zap.String("content", string(input.Contents())), zap.Error(err))
		return input, fmt.Errorf("could not unmarshal transaction yaml: %w", err)
	}

	tran, err = t.mapTransactionSteps(ctx, input, tran)
	if err != nil {
		return input, fmt.Errorf("could not map transaction steps: %w", err)
	}

	marshalled, err := yaml.Marshal(tran)
	if err != nil {
		return input, fmt.Errorf("could not marshal test yaml: %w", err)
	}
	return fileutil.New(input.AbsPath(), marshalled), nil
}

func (t transaction) mapTransactionSteps(ctx context.Context, input fileutil.File, tran openapi.TransactionResource) (openapi.TransactionResource, error) {
	for i, step := range tran.Spec.GetSteps() {
		t.logger.Debug("mapping transaction step",
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
			return openapi.TransactionResource{}, fmt.Errorf("cannot read test file: %w", err)
		}

		testFile, err := t.applyTestFn(ctx, f)
		if err != nil {
			return openapi.TransactionResource{}, fmt.Errorf("cannot apply test '%s': %w", step, err)
		}

		var test openapi.TestResource
		err = yaml.Unmarshal(testFile.Contents(), &test)
		if err != nil {
			return openapi.TransactionResource{}, fmt.Errorf("cannot unmarshal updated test '%s': %w", step, err)
		}

		t.logger.Debug("mapped transaction step",
			zap.Int("index", i),
			zap.String("step", step),
			zap.String("mapped step", test.Spec.GetId()),
		)

		tran.Spec.Steps[i] = test.Spec.GetId()
	}

	return tran, nil
}
