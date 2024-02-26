package comparator

import "fmt"

type Comparator interface {
	Compare(expected, actual string) error
	fmt.Stringer
}

var (
	ErrNoMatch   = fmt.Errorf("no match")
	ErrNotFound  = fmt.Errorf("not found")
	ErrWrongType = fmt.Errorf("wrong type")
)

type Registry interface {
	Get(string) (Comparator, error)
}

func DefaultRegistry() Registry {
	r, _ := NewRegistry(Basic...)
	return r
}

func NewRegistry(comps ...Comparator) (Registry, error) {
	reg := make(map[string]Comparator)
	for _, c := range comps {
		if _, ok := reg[c.String()]; ok {
			return nil, fmt.Errorf(`comparator "%s" already registered`, c.String())
		}
		reg[c.String()] = c
	}
	return registry{reg}, nil
}

type registry struct {
	registry map[string]Comparator
}

func (r registry) Get(s string) (Comparator, error) {
	if c, ok := r.registry[s]; ok {
		return c, nil
	}

	return nil, fmt.Errorf("cannot get comparator %s: %w", s, ErrNotFound)
}
