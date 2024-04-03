package resourcemanager

import (
	"fmt"
	"reflect"
)

func GetResourceType(input any) (string, error) {
	value := reflect.ValueOf(input)
	if value.Kind() != reflect.Struct {
		return "", fmt.Errorf("input must be a struct")
	}

	for i := 0; i < value.NumField(); i++ {
		valueField := value.Field(i)
		typeField := value.Type().Field(i)

		if typeField.Name == "Type" {
			f := valueField.Interface()
			val := reflect.ValueOf(f)
			val = reflect.Indirect(val)

			return val.String(), nil
		}
	}

	return "", fmt.Errorf(`struct has no "Type" field`)
}
