package fs

import (
	"context"
	"time"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/trace"
)

var _ model.RunRepository = &fsDB{}

func (td *fsDB) nextRunID(runs map[int]model.Run) int {
	max := 0
	for key := range runs {
		if key > max {
			max = key
		}
	}

	return max + 1
}

func (td *fsDB) CreateRun(ctx context.Context, test model.Test, run model.Run) (model.Run, error) {
	db := td.runs(test.ID.String())
	runs, err := db.read()
	if err != nil {
		return model.Run{}, err
	}

	run.ID = td.nextRunID(runs)
	run.TestID = test.ID
	run.State = model.RunStateCreated
	run.TestVersion = test.Version
	if run.CreatedAt.IsZero() {
		run.CreatedAt = time.Now()
	}

	err = db.write(run.ID, run)
	if err != nil {
		return model.Run{}, err
	}

	return run, nil
}

func (td *fsDB) UpdateRun(ctx context.Context, run model.Run) error {
	db := td.runs(run.TestID.String())
	return db.write(run.ID, run)
}

func (td *fsDB) DeleteRun(ctx context.Context, r model.Run) error {
	panic("DeleteRun not implemented")
}

func (td *fsDB) GetRun(ctx context.Context, testID id.ID, runID int) (model.Run, error) {
	return td.runs(testID.String()).get(runID)
}

func (td *fsDB) GetTestRuns(ctx context.Context, test model.Test, take, skip int32) (model.List[model.Run], error) {
	runs, err := td.runs(test.ID.String()).read()
	if err != nil {
		return model.List[model.Run]{}, err
	}

	res := model.List[model.Run]{
		Items:      make([]model.Run, 0, take),
		TotalCount: len(runs),
	}

	// TODO: paginate
	for _, run := range runs {
		res.Items = append(res.Items, run)
	}
	return res, nil
}

func (td *fsDB) GetRunByTraceID(ctx context.Context, traceID trace.TraceID) (model.Run, error) {
	panic("GetRunByTraceID not implemented")
}
