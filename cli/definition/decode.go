package definition

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

func Decode(contents string) (File, error) {
	var f File

	err := yaml.Unmarshal([]byte(contents), &f)
	if err != nil {
		return File{}, fmt.Errorf("cannot decode file: %w", err)
	}

	switch f.Type {
	case FileTypeTest:
		var test Test
		err := mapstructure.Decode(f.Spec, &test)
		if err != nil {
			return File{}, fmt.Errorf("cannot decode test: %w", err)
		}
		f.Spec = test
	case FileTypeTransaction:
		var transaction Transaction
		err := mapstructure.Decode(f.Spec, &transaction)
		if err != nil {
			return File{}, fmt.Errorf("cannot decode transaction: %w", err)
		}
		f.Spec = transaction
	default:
		return File{}, fmt.Errorf("invalid file type %s", f.Type)
	}

	return f, nil
}

type FileType string

const (
	FileTypeTest        FileType = "Test"
	FileTypeTransaction FileType = "Transaction"
)

type File struct {
	Type FileType `yaml:"type"`
	Spec any      `yaml:"spec"`
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
