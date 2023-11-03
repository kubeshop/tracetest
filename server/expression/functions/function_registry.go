package functions

import (
	"fmt"
	"sync"
)

type Registry interface {
	Add(string, Invoker, ...Parameter)
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

func (r *registry) Add(name string, function Invoker, argsConfig ...Parameter) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	f := Function{
		name:       name,
		invoker:    function,
		parameters: argsConfig,
	}

	if err := f.Validate(); err != nil {
		// this is a development error. Fail it as fast as possible
		panic(err)
	}

	r.functions[name] = f
}

func (r *registry) Get(name string) (Function, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if function, found := r.functions[name]; found {
		return function, nil
	}

	return Function{}, fmt.Errorf(`could not find function "%s"`, name)
}
