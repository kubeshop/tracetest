package runner

import (
	"fmt"

	"go.uber.org/zap"
)

type Registry struct {
	runners map[string]Runner
	proxies map[string]string
	logger  *zap.Logger
}

func NewRegistry(logger *zap.Logger) Registry {
	return Registry{
		runners: map[string]Runner{},
		proxies: map[string]string{},
		logger:  logger,
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
	if ok {
		return runner, nil // found runner, return it to the user
	}

	// fallback, check if the runner has a proxy
	runnerName, ok := r.proxies[name]
	if !ok {
		return nil, ErrNotFound
	}

	r.logger.Warn(fmt.Sprintf("The resource `%s` is deprecated and will be removed in a future version. Please use `%s` instead.", name, runnerName))
	return r.Get(runnerName)
}

func (r Registry) List() []string {
	var list []string
	for name := range r.runners {
		list = append(list, name)
	}

	return list
}
