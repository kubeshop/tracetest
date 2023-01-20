package linting

import (
	"reflect"

	"github.com/kubeshop/tracetest/server/expression"
)

func DetectMissingVariables(target interface{}, availableVariables []string) []string {
	missingVariables := []string{}

	traverseObject(target, func(f reflect.StructField, value string) {
		tokens := getTokens(f, value)
		variables := make([]string, 0)
		for _, token := range tokens {
			if token.Type == expression.EnvironmentType {
				variables = append(variables, token.Identifier)
			}
		}

		missingVariables = append(missingVariables, getSetDifference(variables, availableVariables)...)
	})

	uniqueMissingVariables := make([]string, 0)
	existingMissingVariables := make(map[string]bool)

	for _, variable := range missingVariables {
		if _, exists := existingMissingVariables[variable]; !exists {
			uniqueMissingVariables = append(uniqueMissingVariables, variable)
			existingMissingVariables[variable] = true
		}
	}

	return uniqueMissingVariables
}

func getTokens(fieldInfo reflect.StructField, value string) []expression.ReflectionToken {
	if fieldInfo.Tag.Get("expr_enabled") == "true" {
		exprTokens, err := expression.GetTokensFromExpression(value)
		if err != nil {
			// probably not an expression, just skip it
			return []expression.ReflectionToken{}
		}

		return exprTokens
	}

	if fieldInfo.Tag.Get("stmt_enabled") == "true" {
		stmtTokens, err := expression.GetTokens(value)
		if err != nil {
			// probably not a statement, just skip it
			return []expression.ReflectionToken{}
		}

		return stmtTokens
	}

	return []expression.ReflectionToken{}
}

func traverseObject(target interface{}, f func(reflect.StructField, string)) {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)

	switch v.Kind() {
	case reflect.Pointer:
		traverseObject(v.Elem(), f)
	case reflect.Struct:
		traverseStruct(v, t, f)
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			traverseObject(item, f)
		}
	}
}

func traverseStruct(v reflect.Value, t reflect.Type, f func(reflect.StructField, string)) {
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.Kind() == reflect.Struct {
			traverseObject(value.Interface(), f)
			continue
		}

		if value.Kind() == reflect.Pointer {
			if !value.IsNil() && value.CanInterface() {
				traverseObject(value.Elem().Interface(), f)
			}
			continue
		}

		if value.Kind() == reflect.Slice {
			traverseArrayField(field, value.Interface(), f)
		}

		if !value.CanInterface() {
			continue
		}

		f(field, value.String())
	}
}

func traverseArrayField(field reflect.StructField, value interface{}, f func(reflect.StructField, string)) {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	switch t.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			if item.Kind() == reflect.String {
				f(field, item.String())
			} else {
				traverseArrayField(field, item.Interface(), f)
			}
		}
	default:
		traverseObject(value, f)
	}
}

func getSetDifference(a, b []string) []string {
	bMap := make(map[string]bool, len(b))
	for _, x := range b {
		bMap[x] = true
	}

	var diff []string
	for _, x := range a {
		if _, found := bMap[x]; !found {
			diff = append(diff, x)
		}
	}

	return diff
}
