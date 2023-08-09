package yaml

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type FileType string

func (ft FileType) String() string {
	return string(ft)
}

const (
	FileTypeTest           FileType = "Test"
	FileTypeTestSuite      FileType = "TestSuite"
	FileTypeEnvironment    FileType = "Environment"
	FileTypeDataStore      FileType = "DataStore"
	FileTypeConfig         FileType = "Config"
	FileTypeDemo           FileType = "Demo"
	FileTypePollingProfile FileType = "PollingProfile"
	FileTypeAnalyzer       FileType = "Analyzer"
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
	case FileTypeTestSuite:
		transaction, err := f.TestSuite()
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

func (f File) TestSuite() (TestSuite, error) {
	if f.Type != FileTypeTestSuite {
		return TestSuite{}, fmt.Errorf("file is not a testsuite")
	}

	transaction, ok := f.Spec.(TestSuite)
	if !ok {
		return TestSuite{}, fmt.Errorf("file spec cannot be casted to a testsuite")
	}

	return transaction, nil
}
