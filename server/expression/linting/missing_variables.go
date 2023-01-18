package linting

import (
	"reflect"

	"github.com/kubeshop/tracetest/server/expression"
)

func DetectMissingVariables(target interface{}, availableVariables []string) []string {
	missingVariables := []string{}

	traverseObject(target, func(f reflect.StructField, v reflect.Value) {
		tokens := getTokens(f, v)
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

func getTokens(f reflect.StructField, v reflect.Value) []expression.ReflectionToken {
	value := v.String()
	if f.Tag.Get("expr_enabled") == "true" {
		exprTokens, err := expression.GetTokensFromExpression(value)
		if err != nil {
			// probably not an expression, just skip it
			return []expression.ReflectionToken{}
		}

		return exprTokens
	}

	if f.Tag.Get("stmt_enabled") == "true" {
		stmtTokens, err := expression.GetTokens(value)
		if err != nil {
			// probably not a statement, just skip it
			return []expression.ReflectionToken{}
		}

		return stmtTokens
	}

	return []expression.ReflectionToken{}
}

func traverseObject(target interface{}, f func(reflect.StructField, reflect.Value)) {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if field.Type.Kind() == reflect.Struct {
			traverseObject(value.Interface(), f)
			continue
		}

		if field.Type.Kind() == reflect.Slice {
			for i := 0; i < value.Len(); i++ {
				item := value.Index(i)
				traverseObject(item.Interface(), f)
			}
		}

		if !value.CanInterface() {
			continue
		}

		f(field, value)
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
