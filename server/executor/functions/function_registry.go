package functions

import (
	"fmt"
	"sync"
)

type Registry interface {
	Add(string, Invoker, ArgConfig)
	Get(string) (Function, error)
}

func newRegistry() Registry {
	return &registry{
		functions: make(map[string]Function, 0),
	}
}

type registry struct {
	functions map[string]Function
	mutex     sync.Mutex
}

func (r *registry) Add(name string, function Invoker, argsConfig ArgConfig) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.functions[name] = Function{
		name:      name,
		invoker:   function,
		argConfig: argsConfig,
	}
}

func (r *registry) Get(name string) (Function, error) {
	if function, found := r.functions[name]; found {
		return function, nil
	}

	return Function{}, fmt.Errorf(`could not find function "%s"`, name)
}
