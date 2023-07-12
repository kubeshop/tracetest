package variable

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
		provider: EnvironmentVariableProvider{},
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
	case reflect.Pointer:
		i.inject(target.Elem())
	case reflect.Struct:
		return i.injectValuesIntoStruct(target)
	case reflect.Interface:
		newValue, err := i.replaceValueInInterface(target)
		if err != nil {
			return err
		}

		if target.CanSet() {
			target.Set(newValue)
		}
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

func (i Injector) replaceValueInInterface(target reflect.Value) (reflect.Value, error) {
	target = target.Elem()
	newValue := reflect.New(target.Type()).Elem()
	for index := 0; index < target.NumField(); index++ {
		newValueField := newValue.Field(index)
		field := target.Field(index)
		if !field.CanInterface() {
			continue
		}

		newValueField.Set(field)
		err := i.inject(newValueField)
		if err != nil {
			return newValue, err
		}
	}

	return newValue, nil
}

func (i Injector) injectValueIntoField(field reflect.Value) error {
	if !field.CanInterface() {
		// cannot replace unexported fields, so skip them.
		return nil
	}

	if !field.CanAddr() {
		return nil
	}

	// We only support variables replacements in strings right now
	strValue := field.String()
	newValue, err := i.ReplaceInString(strValue)
	if err != nil {
		return err
	}

	field.SetString(newValue)
	return nil
}

func (i Injector) ReplaceInString(value string) (string, error) {
	envVarRegex, err := regexp.Compile(`\$\{\w+\}`)
	if err != nil {
		return "", fmt.Errorf("could not compile env variable regex: %w", err)
	}

	allEnvVariables := envVarRegex.FindAllString(value, -1)

	for _, envVariableExpression := range allEnvVariables {
		envVarName := envVariableExpression[2 : len(envVariableExpression)-1] // removes '${' and '}'
		envVarValue, err := i.provider.GetVariable(envVarName)
		if err != nil {
			return "", err
		}

		value = strings.Replace(value, envVariableExpression, envVarValue, -1)
	}

	return value, nil
}
