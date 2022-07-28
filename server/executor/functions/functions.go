package functions

import "fmt"

type FunctionInvoker func(args ...FunctionArg) string
type FunctionArg struct {
	Value string
	Type  string
}

type FunctionArgConfig struct {
	NumberArgs int
	ArgsType   []string
}

type Function struct {
	name      string
	invoker   FunctionInvoker
	argConfig FunctionArgConfig
}

func (f Function) Invoke(args ...FunctionArg) (string, error) {
	if len(args) != f.argConfig.NumberArgs {
		return "", fmt.Errorf(`wrong number of arguments for "%s". Expected %d args, got %d`, f.name, f.argConfig.NumberArgs, len(args))
	}

	for i, argType := range f.argConfig.ArgsType {
		if args[i].Type != argType {
			return "", fmt.Errorf("wrong argument type: argument %d should be of type %s, but it is %s", i, argType, args[i].Type)
		}
	}

	return f.invoker(args...), nil
}

var emptyArgsConfig = FunctionArgConfig{
	NumberArgs: 0,
	ArgsType:   []string{},
}

func GetFunctionRegistry() FunctionRegistry {
	registry := newFunctionRegistry()

	registry.Add("uuid", generateUUID, emptyArgsConfig)
	registry.Add("firstName", generateFirstName, emptyArgsConfig)
	registry.Add("lastName", generateLastName, emptyArgsConfig)
	registry.Add("fullName", generateFullName, emptyArgsConfig)
	registry.Add("email", generateEmail, emptyArgsConfig)
	registry.Add("phone", generatePhoneNumber, emptyArgsConfig)
	registry.Add("creditCard", generateCreditCard, emptyArgsConfig)
	registry.Add("creditCardCvv", generateCreditCardCVV, emptyArgsConfig)
	registry.Add("randomInt", generateRandomInt, FunctionArgConfig{
		NumberArgs: 2,
		ArgsType:   []string{"number", "number"},
	})

	return registry
}
