package functions

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/expression/types"
)

type Invoker func(args ...types.TypedValue) string

type Function struct {
	name       string
	invoker    Invoker
	parameters []Parameter
}

func (f Function) NumRequiredParams() int {
	count := 0
	for _, param := range f.parameters {
		if !param.optional {
			count++
		}
	}

	return count
}

func (f Function) Validate() error {
	foundOptionalParameterIndex := -1
	for i, param := range f.parameters {
		if param.optional {
			foundOptionalParameterIndex = i
		}

		if foundOptionalParameterIndex > -1 && !param.optional {
			// there's a required parameter after a optional parameter
			// this is invalid in most programming languages and makes it
			// extremely hard to resolve the function, so fail it.
			return fmt.Errorf("argument at index %d is required, but it's after optional argument at index %d", i, foundOptionalParameterIndex)
		}
	}

	return nil
}

type Parameter struct {
	pType    types.Type
	optional bool
}

func (f Function) Invoke(args ...types.TypedValue) (types.TypedValue, error) {
	numberRequiredParams := f.NumRequiredParams()
	if len(args) < numberRequiredParams {
		return types.TypedValue{}, fmt.Errorf("missing required parameters. Expected at least %d params, got %d", numberRequiredParams, len(args))
	}

	for i, arg := range args {
		param := f.parameters[i]
		if arg.Type != param.pType {
			return types.TypedValue{}, fmt.Errorf("invalid argument type on index %d: expected %s, got %s", i, arg.Type.String(), param.pType.String())
		}
	}

	returnedValue := f.invoker(args...)
	return types.GetTypedValue(returnedValue), nil
}

func DefaultRegistry() Registry {
	registry := newRegistry()

	registry.Add("uuid", generateUUID)
	registry.Add("firstName", generateFirstName)
	registry.Add("lastName", generateLastName)
	registry.Add("fullName", generateFullName)
	registry.Add("email", generateEmail)
	registry.Add("phone", generatePhoneNumber)
	registry.Add("date", generateDate, OptionalParam(types.TypeString))
	registry.Add("dateTime", generateDateTime, OptionalParam(types.TypeString))
	registry.Add("creditCard", generateCreditCard)
	registry.Add("creditCardCvv", generateCreditCardCVV)
	registry.Add("creditCardExpDate", generateCreditCardExpiration)
	registry.Add("randomInt", generateRandomInt, Param(types.TypeNumber), Param(types.TypeNumber))

	return registry
}

func Param(pType types.Type) Parameter {
	return Parameter{pType: pType, optional: false}
}

func OptionalParam(pType types.Type) Parameter {
	return Parameter{pType: pType, optional: true}
}
