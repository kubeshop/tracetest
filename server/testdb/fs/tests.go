package fs

import (
	"context"
	"os"
	"path"
	"time"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml/yamlconvert"
	"github.com/kubeshop/tracetest/server/testdb"
)

var _ model.TestRepository = &fsDB{}

func (td *fsDB) TestIDExists(ctx context.Context, id id.ID) (bool, error) {
	panic("TestIDExists not implemented")
}

func (td *fsDB) CreateTest(ctx context.Context, test model.Test) (model.Test, error) {
	if !test.HasID() {
		test.ID = testdb.IDGen.ID()
	}

	test.Version = 1
	test.CreatedAt = time.Now()

	b, err := yamlconvert.Test(test).Encode()
	if err != nil {
		return model.Test{}, err
	}

	err = os.WriteFile(td.testFName(test), b, 0644)
	if err != nil {
		return model.Test{}, err
	}

	return test, nil
}

func (td *fsDB) testFName(test model.Test) string {
	return path.Join(td.root, model.Slug(test.Name)+".yaml")
}

func (td *fsDB) UpdateTest(ctx context.Context, test model.Test) (model.Test, error) {
	testFile, err := td.getTestFile(test.ID)
	if err != nil {
		return model.Test{}, err
	}

	err = testFile.write(yamlconvert.Test(test))
	if err != nil {
		return model.Test{}, err
	}

	return test, nil
}

func (td *fsDB) DeleteTest(ctx context.Context, test model.Test) error {
	panic("DeleteTest not implemented")
}

func (td *fsDB) GetTestVersion(ctx context.Context, id id.ID, version int) (model.Test, error) {
	// ignore version for fs db
	return td.GetLatestTestVersion(ctx, id)
}

func (td *fsDB) GetLatestTestVersion(ctx context.Context, id id.ID) (model.Test, error) {
	tests, err := td.getAllTests()
	if err != nil {
		return model.Test{}, err
	}

	for _, test := range tests {
		if test.ID == id {
			return test, nil
		}
	}

	return model.Test{}, nil

}

func (td *fsDB) GetTests(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.Test], error) {
	tests, err := td.getAllTests()
	if err != nil {
		return model.List[model.Test]{}, err
	}
	// TODO: paginate
	res := model.List[model.Test]{
		Items:      tests,
		TotalCount: len(tests),
	}
	return res, nil
}

func (td *fsDB) getTestFile(id id.ID) (file, error) {
	files, err := td.getAllTestFiles()
	if err != nil {
		return file{}, err
	}

	for _, f := range files {
		yf, _ := f.read()
		test, _ := yf.Test()
		if test.Model().ID == id {
			return f, nil
		}
	}

	return file{}, testdb.ErrNotFound

}

func (td *fsDB) getAllTests() ([]model.Test, error) {
	files, err := td.getAllTestFiles()
	if err != nil {
		return nil, err
	}

	res := []model.Test{}
	for _, f := range files {
		yf, _ := f.read()
		test, _ := yf.Test()
		res = append(res, test.Model())
	}

	return res, nil
}

func (td *fsDB) getAllTestFiles() ([]file, error) {
	files, err := td.getFiles()
	if err != nil {
		return nil, err
	}

	res := []file{}
	for _, f := range files {
		if !f.isYaml() {
			continue
		}

		yf, err := f.read()
		if err != nil {
			return nil, err
		}

		_, err = yf.Test()
		if err != nil {
			// not a test, ignore
			continue
		}
		res = append(res, f)
	}

	return res, nil
}
