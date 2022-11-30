package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/denisbrodbeck/machineid"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"golang.org/x/exp/slices"
)

type fsDB struct {
	root string
}

func New(path string) (model.Repository, error) {
	ps := &fsDB{
		root: path,
	}

	err := ps.ready()
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func (td *fsDB) ready() error {
	err := os.Mkdir(td.root, 0755)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil
		}
		return err
	}

	return nil
}

func (td *fsDB) ServerID() (id string, isNew bool, err error) {
	id, err = td.getConfigValue("serverID")
	if err != nil {
		err = fmt.Errorf("could not get machineID: %w", err)
		return
	}

	if id != "" {
		isNew = false
		return
	}

	// no id, let's creat it
	isNew = true
	id, err = machineid.ProtectedID("tracetest")
	if err != nil {
		err = fmt.Errorf("could not get machineID: %w", err)
		return
	}
	id = id[:10] // limit lenght to avoid issues with GA

	// id, err =
	if err != nil {
		err = fmt.Errorf("could not get machineID: %w", err)
		return
	}
	err = td.writeConfigValue("serverID", id)
	if err != nil {
		err = fmt.Errorf("could not save serverID into DB: %w", err)
		return
	}

	return
}

func (td *fsDB) getConfigValue(key string) (string, error) {
	config, err := td.getDB()
	if err != nil {
		return "", err
	}

	return config[key], nil
}

func (td *fsDB) writeConfigValue(key, value string) error {
	config, err := td.getDB()
	if err != nil {
		return err
	}

	config[key] = value

	return td.persistDB(config)
}

func (td *fsDB) getDB() (map[string]string, error) {
	db, err := os.ReadFile(path.Join(td.root, ".config.json"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// allow working with not existing dbs
			return map[string]string{}, nil
		}

		return nil, err
	}

	var config map[string]string
	err = json.Unmarshal(db, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (td *fsDB) persistDB(config map[string]string) error {
	db, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(td.root, ".config.json"), db, 0644)
}

type file struct {
	path string
	info os.FileInfo
}

func (f file) read() (yaml.File, error) {
	b, err := os.ReadFile(f.path)
	if err != nil {
		return yaml.File{}, fmt.Errorf("cannot read file %s: %w", f.path, err)
	}

	yf, err := yaml.Decode(b)
	if err != nil {
		return yaml.File{}, fmt.Errorf("cannot decode file %s: %w", f.path, err)
	}

	return yf, nil
}

func (f file) isYaml() bool {
	return slices.Contains([]string{".yaml", ".yml"}, filepath.Ext(f.path))
}

func (td *fsDB) getFiles() ([]file, error) {
	return td.readDir(td.root)
}

func (td *fsDB) readDir(path string) ([]file, error) {
	var files []file
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			files = append(files, file{path, info})
			return nil
		})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (td *fsDB) Drop() error {
	panic("not implemented")
}
