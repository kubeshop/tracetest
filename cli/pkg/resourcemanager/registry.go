package resourcemanager

import "fmt"

type Registry struct {
	resources map[string]client
}

func NewRegistry() *Registry {
	return &Registry{
		resources: make(map[string]client),
	}
}

func (r *Registry) Register(c client) {
	r.resources[c.resourceName] = c
}

func (r *Registry) Get(resourceName string) (client, error) {
	c, ok := r.resources[resourceName]
	if !ok {
		return client{}, fmt.Errorf("resource %s not found", resourceName)
	}

	return c, nil
}
