package runner

import (
	"fmt"
)

type Registry map[string]Runner

func NewRegistry() Registry {
	return Registry{}
}

func (r Registry) Register(runner Runner) Registry {
	r[runner.Name()] = runner
	return r
}

var ErrNotFound = fmt.Errorf("runner not found")

func (r Registry) Get(name string) (Runner, error) {
	if runner, ok := r[name]; ok {
		return runner, nil
	}

	return nil, ErrNotFound
}

func (r Registry) List() []string {
	var list []string
	for name := range r {
		list = append(list, name)
	}

	return list
}
