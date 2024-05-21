package processor

import (
	"context"

	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
)

type Preprocessor interface {
	Preprocess(ctx context.Context, input fileutil.File) (fileutil.File, error)
}
