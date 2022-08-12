package functions

import "fmt"

type Invoker func(args ...Arg) string
type Arg struct {
	Value string
	Type  string
}

type ArgTypes []string

type Function struct {
	name     string
	invoker  Invoker
	argTypes ArgTypes
}

func (f Function) Invoke(args ...Arg) (string, error) {
	if len(args) != len(f.argTypes) {
		return "", fmt.Errorf(`wrong number of arguments for "%s". Expected %d args, got %d`, f.name, len(f.argTypes), len(args))
	}

	for i, argType := range f.argTypes {
		if args[i].Type != argType {
			return "", fmt.Errorf("wrong argument type: argument %d should be of type %s, but it is %s", i, argType, args[i].Type)
		}
	}

	return f.invoker(args...), nil
}

var emptyArgsConfig = ArgTypes{}

func GetFunctionRegistry() Registry {
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
	registry.Add("randomInt", generateRandomInt, ArgTypes{"number", "number"})

	return registry
}
