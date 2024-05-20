package processor

import (
	"context"

	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
)

type Postprocessor interface {
	Postprocess(ctx context.Context, input fileutil.File) error
}
