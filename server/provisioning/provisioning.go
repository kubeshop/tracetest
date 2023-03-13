package provisioning

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/server/config/configresource"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

func New(db model.Repository, configDB *configresource.Repository) provisioner {
	return provisioner{db, configDB}
}

type provisioner struct {
	db       model.Repository
	configDB *configresource.Repository
}

var (
	ErrEnvEmpty = errors.New("cannot read provisioning from env variable TRACETEST_PROVISIONING: variable is empty")
)

func (p provisioner) FromEnv() error {
	envVar := os.Getenv("TRACETEST_PROVISIONING")
	if envVar == "" {
		return ErrEnvEmpty
	}

	data, err := base64.StdEncoding.DecodeString(envVar)
	if err != nil {
		return fmt.Errorf("cannot decode env variable TRACETEST_PROVISIONING: %w", err)
	}
	return p.do(data)
}

var ErrFileNotExists = errors.New("provisioning file does not exists")

func (p provisioner) FromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExists
		}
		return fmt.Errorf("cannot read provisioning file '%s'", path)
	}

	return p.do(data)
}

type provisionConfig struct {
	DataStore dataStore `mapstructure:"dataStore"`
	Config    config    `mapstructure:"config"`
}

func (p provisioner) do(data []byte) error {
	var rawYaml map[string]interface{}
	err := yaml.Unmarshal(data, &rawYaml)
	if err != nil {
		return fmt.Errorf("cannot unmarshal yaml: %w", err)
	}

	config := provisionConfig{}
	mapstructure.Decode(rawYaml, &config)
	if err != nil {
		return fmt.Errorf("cannot unmarshal yaml: %w", err)
	}

	// Provision data store
	dsModel := config.DataStore.model()
	dsModel.IsDefault = true

	_, err = p.db.CreateDataStore(context.TODO(), dsModel)
	if err != nil {
		return fmt.Errorf("cannot provision data store: %w", err)
	}

	// Provision config
	cfgModel := config.Config.model()
	_, err = p.configDB.Update(context.TODO(), cfgModel)
	if err != nil {
		return fmt.Errorf("cannot provision config: %w", err)
	}

	return nil
}
