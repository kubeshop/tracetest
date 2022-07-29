package replacer

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Variables map[string]string

type Injector struct {
	provider VariableProvider
}

func NewInjector(options ...InjectorOption) Injector {
	injector := Injector{
		provider: expressionValueProvider{},
	}

	for _, option := range options {
		option(&injector)
	}

	return injector
}

func (i Injector) Inject(target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("could not inject value because target is not a pointer")
	}

	targetStruct := targetValue.Elem()

	return i.inject(targetStruct)
}

func (i Injector) inject(target reflect.Value) error {
	switch target.Kind() {
	case reflect.Struct:
		return i.injectValuesIntoStruct(target)
	case reflect.String:
		return i.injectValueIntoField(target)
	case reflect.Slice:
		for index := 0; index < target.Len(); index++ {
			item := target.Index(index)
			err := i.inject(item)
			if err != nil {
				return err
			}
		}
	case reflect.Ptr:
		value := target.Elem()
		return i.inject(value)
	}

	return nil
}

func (i Injector) injectValuesIntoStruct(target reflect.Value) error {
	for index := 0; index < target.NumField(); index++ {
		field := target.Field(index)
		err := i.inject(field)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i Injector) injectValueIntoField(field reflect.Value) error {
	if !field.CanSet() {
		// Unexported field. Should be skipped
		return nil
	}

	// We only support variables replacements in strings right now
	strValue := field.String()
	newValue, err := i.replaceExpression(strValue)
	if err != nil {
		return err
	}

	field.SetString(newValue)
	return nil
}

func (i Injector) replaceExpression(value string) (string, error) {
	expressionRegex, err := regexp.Compile(`\{\{([^\}]*)\}\}`)
	if err != nil {
		return "", fmt.Errorf("could not compile env variable regex: %w", err)
	}

	allExpressions := expressionRegex.FindAllString(value, -1)

	for _, expression := range allExpressions {
		expressionText := expression[2 : len(expression)-2] // removes '{{' and '}}'
		expressionText = strings.TrimSpace(expressionText)
		expressionValue, err := i.provider.GetVariable(expressionText)
		if err != nil {
			return "", err
		}

		value = strings.Replace(value, expression, expressionValue, -1)
	}

	return value, nil
}
