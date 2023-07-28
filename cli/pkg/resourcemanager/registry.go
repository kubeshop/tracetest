package resourcemanager

import (
	"errors"
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
