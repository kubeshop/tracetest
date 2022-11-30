package fs

import (
	"context"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/trace"
)

var _ model.RunRepository = &fsDB{}

func (td *fsDB) CreateRun(ctx context.Context, test model.Test, run model.Run) (model.Run, error) {
	panic("not implemented")
}

func (td *fsDB) UpdateRun(ctx context.Context, r model.Run) error {
	panic("not implemented")
}

func (td *fsDB) DeleteRun(ctx context.Context, r model.Run) error {
	panic("not implemented")
}

func (td *fsDB) GetRun(ctx context.Context, testID id.ID, runID int) (model.Run, error) {
	panic("not implemented")
}

func (td *fsDB) GetTestRuns(ctx context.Context, test model.Test, take, skip int32) (model.List[model.Run], error) {
	panic("not implemented")
}

func (td *fsDB) GetRunByTraceID(ctx context.Context, traceID trace.TraceID) (model.Run, error) {
	panic("not implemented")
}
