package linting

import (
	"fmt"
	"reflect"
	"time"

	"github.com/kubeshop/tracetest/server/expression"
)

func DetectMissingVariables(target interface{}, availableVariables []string) []string {
	missingVariables := []string{}

	traverse(target, func(f reflect.StructField, value string) {
		tokens := getTokens(f, value)
		variables := make([]string, 0)
		for _, token := range tokens {
			fmt.Println("@@@TOKEN", token, token.Identifier)
			if token.Type == expression.EnvironmentType {
				variables = append(variables, token.Identifier)
			}

			if token.Type == expression.VariableType {
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

func getTokens(fieldInfo reflect.StructField, value string) []expression.Token {
	if fieldInfo.Tag.Get("expr_enabled") == "true" {
		value = fmt.Sprintf("'%s'", value)
		exprTokens, err := expression.GetTokensFromExpression(value)
		if err != nil {
			// probably not an expression, just skip it
			return []expression.Token{}
		}

		return exprTokens
	}

	if fieldInfo.Tag.Get("stmt_enabled") == "true" {
		stmtTokens, err := expression.GetTokens(value)
		if err != nil {
			// probably not a statement, just skip it
			return []expression.Token{}
		}

		return stmtTokens
	}

	return []expression.Token{}
}

func traverse(target interface{}, f func(reflect.StructField, string)) {
	typ := reflect.TypeOf(target)
	value := reflect.ValueOf(target)

	traverseValue(value, typ, reflect.StructField{}, f)
}

func traverseValue(value reflect.Value, typ reflect.Type, parentField reflect.StructField, f func(reflect.StructField, string)) {
	switch value.Kind() {
	case reflect.Pointer:
		if value.IsZero() {
			return
		}

		realValue := value.Elem().Interface()
		typ := reflect.TypeOf(realValue)
		value := reflect.ValueOf(realValue)
		traverseValue(value, typ, parentField, f)
	case reflect.Struct:
		// Don't traverse time fields. Time fields are structs with internal information,
		// we don't have to go thought its internals to find missing variables. So,
		// just call the callback on them directly instead.
		if !typ.AssignableTo(reflect.TypeOf(time.Now())) {
			traverseStruct(value, typ, f)
		} else {
			f(parentField, value.String())
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			item := value.Index(i)
			traverseValue(item, item.Type(), parentField, f)
		}
	default:
		f(parentField, value.String())
	}
}

func traverseStruct(value reflect.Value, typ reflect.Type, f func(reflect.StructField, string)) {
	for i := 0; i < value.NumField(); i++ {
		fieldType := typ.Field(i)
		field := value.Field(i)
		typ := field.Type()
		traverseValue(field, typ, fieldType, f)
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
