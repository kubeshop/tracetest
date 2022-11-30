package fs

import (
	"context"

	"github.com/kubeshop/tracetest/server/model"
)

var _ model.EnvironmentRepository = &fsDB{}

func (td *fsDB) CreateEnvironment(ctx context.Context, environment model.Environment) (model.Environment, error) {
	panic("not implemented")
}

func (td *fsDB) UpdateEnvironment(ctx context.Context, environment model.Environment) (model.Environment, error) {
	panic("not implemented")
}

func (td *fsDB) DeleteEnvironment(ctx context.Context, environment model.Environment) error {
	panic("not implemented")
}

func (td *fsDB) GetEnvironments(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.Environment], error) {
	panic("not implemented")
}

func (td *fsDB) GetEnvironment(ctx context.Context, id string) (model.Environment, error) {
	panic("not implemented")
}

func (td *fsDB) EnvironmentIDExists(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}
