package resourcemanager

import (
	"fmt"
	"sort"
)

type Registry struct {
	resources map[string]client
}

func NewRegistry() *Registry {
	return &Registry{
		resources: make(map[string]client),
	}
}

func (r *Registry) Register(c client) *Registry {
	r.resources[c.resourceName] = c
	return r
}

func (r *Registry) Get(resourceName string) (client, error) {
	c, ok := r.resources[resourceName]
	if !ok {
		return client{}, fmt.Errorf("resource %s not found", resourceName)
	}

	return c, nil
}

func (r *Registry) List() []string {
	var resources []string
	for k := range r.resources {
		resources = append(resources, k)
	}

	sort.Strings(resources)

	return resources
}
