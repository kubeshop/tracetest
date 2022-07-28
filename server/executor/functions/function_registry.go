package functions

import "fmt"

type FunctionRegistry interface {
	Add(string, FunctionInvoker, FunctionArgConfig)
	Get(string) (Function, error)
}

func newFunctionRegistry() FunctionRegistry {
	return &functionRegistry{
		functions: make(map[string]Function, 0),
	}
}

type functionRegistry struct {
	functions map[string]Function
}

func (r *functionRegistry) Add(name string, function FunctionInvoker, argsConfig FunctionArgConfig) {
	r.functions[name] = Function{
		name:      name,
		invoker:   function,
		argConfig: argsConfig,
	}
}

func (r *functionRegistry) Get(name string) (Function, error) {
	if function, found := r.functions[name]; found {
		return function, nil
	}

	return Function{}, fmt.Errorf(`could not find function "%s"`, name)
}
