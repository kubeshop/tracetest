package functions

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/expression/types"
)

type Invoker func(args ...types.TypedValue) string

type Function struct {
	name     string
	invoker  Invoker
	argTypes []types.Type
}

func (f Function) Invoke(args ...types.TypedValue) (types.TypedValue, error) {
	if len(args) != len(f.argTypes) {
		return types.TypedValue{}, fmt.Errorf("invalid number of arguments. Expected %d, got %d", len(f.argTypes), len(args))
	}

	for i, arg := range args {
		expectedArgType := f.argTypes[i]
		if arg.Type != expectedArgType {
			return types.TypedValue{}, fmt.Errorf("invalid argument type on index %d: expected %s, got %s", i, arg.Type.String(), expectedArgType.String())
		}
	}

	returnedValue := f.invoker(args...)
	return types.GetTypedValue(returnedValue), nil
}

func GetFunctionRegistry() Registry {
	emptyArgsConfig := []types.Type{}
	registry := newRegistry()

	registry.Add("uuid", generateUUID, emptyArgsConfig)
	registry.Add("firstName", generateFirstName, emptyArgsConfig)
	registry.Add("lastName", generateLastName, emptyArgsConfig)
	registry.Add("fullName", generateFullName, emptyArgsConfig)
	registry.Add("email", generateEmail, emptyArgsConfig)
	registry.Add("phone", generatePhoneNumber, emptyArgsConfig)
	registry.Add("creditCard", generateCreditCard, emptyArgsConfig)
	registry.Add("creditCardCvv", generateCreditCardCVV, emptyArgsConfig)
	registry.Add("creditCardExpDate", generateCreditCardExpiration, emptyArgsConfig)
	registry.Add("randomInt", generateRandomInt, []types.Type{types.TypeNumber, types.TypeNumber})

	return registry
}
