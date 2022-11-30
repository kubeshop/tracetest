package fs

import (
	"context"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
)

var _ model.TestRepository = &fsDB{}

func (td *fsDB) TestIDExists(ctx context.Context, id id.ID) (bool, error) {
	panic("not implemented")
}

func (td *fsDB) CreateTest(ctx context.Context, test model.Test) (model.Test, error) {
	panic("not implemented")
}

func (td *fsDB) UpdateTest(ctx context.Context, test model.Test) (model.Test, error) {
	panic("not implemented")
}

func (td *fsDB) DeleteTest(ctx context.Context, test model.Test) error {
	panic("not implemented")
}

func (td *fsDB) GetTestVersion(ctx context.Context, id id.ID, version int) (model.Test, error) {
	panic("not implemented")
}

func (td *fsDB) GetLatestTestVersion(ctx context.Context, id id.ID) (model.Test, error) {
	panic("not implemented")
}

func (td *fsDB) GetTests(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.Test], error) {
	files, err := td.getFiles()
	if err != nil {
		return model.List[model.Test]{}, err
	}

	res := model.List[model.Test]{}
	for _, f := range files {
		if !f.isYaml() {
			continue
		}

		yf, err := f.read()
		if err != nil {
			return model.List[model.Test]{}, err
		}

		test, err := yf.Test()
		if err != nil {
			// not a test, ignore
			continue
		}

		res.Items = append(res.Items, test.Model())
	}

	return res, nil
}
