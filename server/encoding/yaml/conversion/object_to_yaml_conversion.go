package conversion

import (
	"reflect"

	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"gopkg.in/yaml.v2"
)

func GetYamlFileFromDefinition(def definition.Test) ([]byte, error) {
	yamlBytes, err := yaml.Marshal(def)
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
