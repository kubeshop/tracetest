package fs

import (
	"context"
	"path"

	"github.com/joho/godotenv"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testdb"
)

var _ model.EnvironmentRepository = &fsDB{}

func (td *fsDB) envFName(env model.Environment) string {
	return path.Join(td.root, env.ID+".env")
}

func (td *fsDB) CreateEnvironment(ctx context.Context, env model.Environment) (model.Environment, error) {
	env.ID = env.Slug()

	err := godotenv.Write(
		encodeVars(env),
		td.envFName(env),
	)
	if err != nil {
		return model.Environment{}, err
	}

	return env, nil
}

func encodeVars(env model.Environment) map[string]string {
	ret := map[string]string{}
	for _, val := range env.Values {
		ret[val.Key] = val.Value
	}
	return ret
}

func (td *fsDB) UpdateEnvironment(ctx context.Context, env model.Environment) (model.Environment, error) {
	err := godotenv.Write(
		encodeVars(env),
		td.envFName(env),
	)
	if err != nil {
		return model.Environment{}, err
	}

	return env, nil
}

func (td *fsDB) DeleteEnvironment(ctx context.Context, env model.Environment) error {
	panic("DeleteEnvironment not implemented")
}

func (td *fsDB) GetEnvironments(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.Environment], error) {
	envs, err := td.getAllEnvs()
	if err != nil {
		return model.List[model.Environment]{}, err
	}

	// TODO: paginate
	res := model.List[model.Environment]{
		Items:      envs,
		TotalCount: len(envs),
	}

	return res, nil
}

func (td *fsDB) GetEnvironment(ctx context.Context, id string) (model.Environment, error) {
	envs, err := td.getAllEnvs()
	if err != nil {
		return model.Environment{}, err
	}
	for _, env := range envs {
		if env.ID == id {
			return env, nil
		}
	}
	return model.Environment{}, testdb.ErrNotFound
}

func (td *fsDB) EnvironmentIDExists(ctx context.Context, id string) (bool, error) {
	_, err := td.GetEnvironment(ctx, id)
	return err == nil, err
}

func (td *fsDB) getAllEnvs() ([]model.Environment, error) {
	files, err := td.getAllEnvFiles()
	if err != nil {
		return nil, err
	}

	res := []model.Environment{}
	for _, f := range files {
		envVars, err := f.readEnv()
		if err != nil {
			return nil, err
		}

		vals := make([]model.EnvironmentValue, 0, len(envVars))
		for k, v := range envVars {
			vals = append(vals, model.EnvironmentValue{k, v})
		}

		res = append(res, model.Environment{
			ID:     model.Slug(f.path),
			Name:   f.info.Name(),
			Values: vals,
		})
	}

	return res, nil
}

func (td *fsDB) getAllEnvFiles() ([]file, error) {
	files, err := td.getFiles()
	if err != nil {
		return nil, err
	}

	res := []file{}
	for _, f := range files {
		if !f.isEnv() {
			continue
		}
		res = append(res, f)
	}

	return res, nil
}
