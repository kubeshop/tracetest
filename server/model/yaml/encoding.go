package yaml

import (
	"bytes"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

func Encode(in File) ([]byte, error) {
	b := new(bytes.Buffer)
	enc := yaml.NewEncoder(b)
	defer enc.Close()
	enc.SetIndent(2)
	err := enc.Encode(in)
	if err != nil {
		return nil, fmt.Errorf("cannot encode File: %w", err)
	}

	return b.Bytes(), nil

}

func Decode(contents []byte) (File, error) {
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
