package resourcemanager

import (
	"errors"
	"fmt"
	"sort"

	"github.com/agnivade/levenshtein"
)

type Registry struct {
	resources map[string]Client
}

func NewRegistry() *Registry {
	return &Registry{
		resources: make(map[string]Client),
	}
}

func (r *Registry) Register(c Client) *Registry {
	r.resources[c.resourceName] = c
	return r
}

var ErrResourceNotFound = errors.New("resource not found")

func (r *Registry) Get(resourceName string) (Client, error) {
	c, ok := r.resources[resourceName]
	if !ok {
		return Client{}, ErrResourceNotFound
	}

	if c.options.proxyResource != "" {
		c.logger.Warn(fmt.Sprintf("The resource `%s` is deprecated and will be removed in a future version. Please use `%s` instead.", c.resourceName, c.options.proxyResource))
		return r.Get(c.options.proxyResource)
	}

	return c, nil
}

func (r *Registry) Exists(resourceName string) bool {
	c, ok := r.resources[resourceName]
	if !ok {
		return false
	}

	if c.options.proxyResource != "" {
		return r.Exists(c.options.proxyResource)
	}

	return true
}

func (r *Registry) List() []string {
	var resources []string
	for k, c := range r.resources {
		if c.options.proxyResource == "" {
			resources = append(resources, k)
		}
	}

	sort.Strings(resources)

	return resources
}

const minDistanceForSuggestion = 2

func (r *Registry) Suggest(input string) string {
	for resource := range r.resources {
		distance := levenshtein.ComputeDistance(input, resource)
		if distance <= minDistanceForSuggestion {
			return resource
		}

	}

	return ""
}
