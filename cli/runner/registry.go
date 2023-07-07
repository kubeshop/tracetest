package runner

import (
	"fmt"
)

type Registry map[string]Runner

func NewRegistry() Registry {
	return Registry{}
}

func (r Registry) Register(runner Runner) {
	r[runner.Name()] = runner
}

var ErrNotFound = fmt.Errorf("runner not found")

func (r Registry) Get(name string) (Runner, error) {
	if runner, ok := r[name]; ok {
		return runner, nil
	}

	return nil, ErrNotFound
}
