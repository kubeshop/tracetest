package runner

import (
	"fmt"
)

type Registry struct {
	runners map[string]Runner
	proxies map[string]string
}

func NewRegistry() Registry {
	return Registry{
		runners: map[string]Runner{},
		proxies: map[string]string{},
	}
}

func (r Registry) Register(runner Runner) Registry {
	r.runners[runner.Name()] = runner
	return r
}

func (r Registry) RegisterProxy(proxyName, runnerName string) Registry {
	r.proxies[proxyName] = runnerName
	return r
}

var ErrNotFound = fmt.Errorf("runner not found")

func (r Registry) Get(name string) (Runner, error) {
	runner, ok := r.runners[name]
	if !ok {
		if runnerName, ok := r.proxies[name]; ok {
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
	for name := range r.runners {
		list = append(list, name)
	}

	return list
}
