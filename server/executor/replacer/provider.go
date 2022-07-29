package replacer

import (
	"fmt"
	"regexp"

	"github.com/kubeshop/tracetest/server/executor/functions"
)

type VariableProvider interface {
	GetVariable(string) (string, error)
}

type expressionValueProvider struct{}

var _ VariableProvider = &expressionValueProvider{}

func (p expressionValueProvider) GetVariable(expression string) (string, error) {
	functionRegex := regexp.MustCompile(`[a-zA-Z_][a-zA-Z_0-9]*\(.*\)`)

	if functionRegex.Match([]byte(expression)) {
		// It's a function
		return p.getFunctionValue(expression)
	}

	// TODO: add support for variables later
	return "", fmt.Errorf(`expression "%s" has no implementation`, expression)
}

func (p expressionValueProvider) getFunctionValue(expression string) (string, error) {
	function, args, err := functions.ParseFunction(expression)
	if err != nil {
		return "", err
	}

	return function.Invoke(args...)
}
