package runner

import (
	"fmt"
)

type Registry struct {
	registry map[string]Runner
	proxy    map[string]string
}

func NewRegistry() Registry {
	return Registry{
		registry: map[string]Runner{},
		proxy:    map[string]string{},
	}
}

func (r Registry) Register(runner Runner) Registry {
	r.registry[runner.Name()] = runner
	return r
}

func (r Registry) RegisterProxy(proxyName, runnerName string) Registry {
	r.proxy[proxyName] = runnerName
	return r
}

var ErrNotFound = fmt.Errorf("runner not found")

func (r Registry) Get(name string) (Runner, error) {
	runner, ok := r.registry[name]
	if !ok {
		if runnerName, ok := r.proxy[name]; ok {
			if !ok {
				return nil, ErrNotFound
			}

			return r.Get(runnerName)
		}
	}

	return runner, nil
}

func (r Registry) List() []string {
	var list []string
	for name := range r.registry {
		list = append(list, name)
	}

	return list
}
