package yaml

import (
	"fmt"
	"reflect"

	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

func Encode(in File) ([]byte, error) {
	return getYaml(in)
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
	case FileTypeDataStore:
		var dataStore openapi.DataStore
		err := mapstructure.Decode(f.Spec, &dataStore)
		if err != nil {
			return File{}, fmt.Errorf("cannot decode datastore: %w", err)
		}
		f.Spec = dataStore
	default:
		return File{}, fmt.Errorf("invalid file type %s", f.Type)
	}

	return f, nil
}

func getYaml(in File) ([]byte, error) {
	yamlBytes, err := yaml.Marshal(in)
	if err != nil {
		return []byte{}, nil
	}

	mapSlice := yaml.MapSlice{}
	err = yaml.Unmarshal(yamlBytes, &mapSlice)
	if err != nil {
		return []byte{}, nil
	}

	mapSlice = removeEmptyFields(mapSlice)

	bytes, err := yaml.Marshal(mapSlice)
	if err != nil {
		return []byte{}, nil
	}

	return bytes, nil
}

func removeEmptyFields(slice yaml.MapSlice) yaml.MapSlice {
	newMap := make(yaml.MapSlice, 0)
	for _, entry := range slice {
		key, value := entry.Key, entry.Value
		tValue := reflect.ValueOf(value)

		if tValue.Kind() == reflect.Slice {
			newValue, length := removeEmptyFieldsFromSlice(value)
			if length > 0 {
				newMap = append(newMap, yaml.MapItem{Key: key, Value: newValue})
			}
		} else if !tValue.IsZero() {
			newMap = append(newMap, entry)
		}
	}

	return newMap
}

func removeEmptyFieldsFromSlice(slice interface{}) (interface{}, int) {
	// The slice can be a MapSlice or a slice of objects which
	// can be either another slice of interface{} or a MapSlice.
	// So we have to check for both cases before doing anything
	if mapSlice, ok := slice.(yaml.MapSlice); ok {
		newSlice := removeEmptyFields(mapSlice)
		return newSlice, len(newSlice)
	}

	// If this was not a MapSlice, this is can be a slice of interface{}
	if genericSlice, ok := slice.([]interface{}); ok {
		if len(genericSlice) == 0 {
			return yaml.MapSlice{}, 0
		}

		firstItem := genericSlice[0]
		if reflect.ValueOf(firstItem).Kind() != reflect.Slice {
			return genericSlice, len(genericSlice)
		}

		newSlice := make([]interface{}, 0)
		for _, item := range genericSlice {
			if reflect.ValueOf(item).Kind() == reflect.Slice {
				newItem, length := removeEmptyFieldsFromSlice(item)
				if length > 0 {
					newSlice = append(newSlice, newItem)
				}
			}
		}

		return newSlice, len(newSlice)
	}

	return yaml.MapSlice{}, 0
}
