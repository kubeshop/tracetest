package provisioning

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

func New(db model.Repository) provisioner {
	return provisioner{db}
}

type provisioner struct {
	db model.Repository
}

func (p provisioner) FromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read provisioning file '%s'", path)
	}

	return p.do(data)
}

type provisionConfig struct {
	DataStore dataStore `mapstructure:"dataStore"`
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

	dsModel := config.DataStore.model()
	dsModel.IsDefault = true

	_, err = p.db.CreateDataStore(context.TODO(), dsModel)
	if err != nil {
		return fmt.Errorf("cannot provision data store: %w", err)
	}

	return nil
}
