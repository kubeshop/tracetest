package yaml

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/config/configresource"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"gopkg.in/yaml.v3"
)

type FileType string

func (ft FileType) String() string {
	return string(ft)
}

const (
	FileTypeTest           FileType = "Test"
	FileTypeTransaction    FileType = "Transaction"
	FileTypeEnvironment    FileType = "Environment"
	FileTypeDataStore      FileType = "DataStore"
	FileTypeConfig         FileType = "Config"
	FileTypeDemo           FileType = "Demo"
	FileTypePollingProfile FileType = "PollingProfile"
)

type File struct {
	Type FileType `yaml:"type"`
	Spec any      `yaml:"spec"`
}

func (f File) Encode() ([]byte, error) {
	return yaml.Marshal(f)
}

func (f File) Validate() error {
	switch f.Type {
	case FileTypeTest:
		test, err := f.Test()
		if err != nil {
			return err
		}
		return test.Validate()
	case FileTypeTransaction:
		transaction, err := f.Transaction()
		if err != nil {
			return err
		}
		return transaction.Validate()
	}
	return fmt.Errorf("invalid file type %s", f.Type)
}

func (f File) Test() (Test, error) {
	if f.Type != FileTypeTest {
		return Test{}, fmt.Errorf("file is not a test")
	}

	test, ok := f.Spec.(Test)
	if !ok {
		return Test{}, fmt.Errorf("file spec cannot be casted to a test")
	}

	return test, nil
}

func (f File) Config() (configresource.Config, error) {
	if f.Type != FileTypeConfig {
		return configresource.Config{}, fmt.Errorf("file is not a test")
	}

	config, ok := f.Spec.(configresource.Config)
	if !ok {
		return configresource.Config{}, fmt.Errorf("file spec cannot be casted to a test")
	}

	return config, nil
}

func (f File) PollingProfile() (pollingprofile.PollingProfile, error) {
	if f.Type != FileTypePollingProfile {
		return pollingprofile.PollingProfile{}, fmt.Errorf("file is not a test")
	}

	profile, ok := f.Spec.(pollingprofile.PollingProfile)
	if !ok {
		return pollingprofile.PollingProfile{}, fmt.Errorf("file spec cannot be casted to a test")
	}

	return profile, nil
}

func (f File) Transaction() (Transaction, error) {
	if f.Type != FileTypeTransaction {
		return Transaction{}, fmt.Errorf("file is not a transaction")
	}

	transaction, ok := f.Spec.(Transaction)
	if !ok {
		return Transaction{}, fmt.Errorf("file spec cannot be casted to a transaction")
	}

	return transaction, nil
}

func (f File) Environment() (Environment, error) {
	if f.Type != FileTypeEnvironment {
		return Environment{}, fmt.Errorf("file is not an environment")
	}

	environment, ok := f.Spec.(Environment)
	if !ok {
		return Environment{}, fmt.Errorf("file spec cannot be casted to an environment")
	}

	return environment, nil
}
