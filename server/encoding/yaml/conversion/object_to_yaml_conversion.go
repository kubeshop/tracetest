package conversion

import (
	"encoding/json"
	"reflect"

	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"gopkg.in/yaml.v2"
)

func GetYamlFileFromDefinition(def definition.Test) ([]byte, error) {
	defMap := make(map[string]interface{}, 0)
	jsonBytes, err := json.Marshal(def)
	if err != nil {
		return []byte{}, nil
	}

	err = json.Unmarshal(jsonBytes, &defMap)
	if err != nil {
		return []byte{}, nil
	}

	defMap = removeEmptyFields(defMap)

	bytes, err := yaml.Marshal(defMap)
	if err != nil {
		return []byte{}, nil
	}

	return bytes, nil
}

func removeEmptyFields(m map[string]interface{}) map[string]interface{} {
	for key, value := range m {
		tValue := reflect.ValueOf(value)
		if tValue.Kind() == reflect.Map {
			m[key] = removeEmptyFields(value.(map[string]interface{}))
			innerMap := m[key].(map[string]interface{})
			if len(innerMap) == 0 {
				delete(m, key)
			}
		}

		if tValue.IsZero() {
			delete(m, key)
		}
	}

	return m
}
