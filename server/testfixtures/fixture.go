package testfixtures

import (
	"fmt"
	"sync"
)

type Fixture[T comparable] struct {
	mutex     sync.Mutex
	value     T
	generator Generator[T]
}

type Generator[T comparable] func(options FixtureOptions) (T, error)

var fixtures = make(map[string]interface{})

func emptyValue[T comparable]() T {
	var empty T
	return empty
}

type FixtureOptions struct {
	DisableCache bool
	Arguments    interface{}
}

type Option func(opt *FixtureOptions)

func RegisterFixture[T comparable](name string, generator Generator[T]) {
	if _, exists := fixtures[name]; exists {
		panic(fmt.Errorf("fixture %s already exists", name))
	}

	fixture := Fixture[T]{
		mutex:     sync.Mutex{},
		value:     emptyValue[T](),
		generator: generator,
	}

	fixtures[name] = &fixture
}

func GetFixtureValue[T comparable](name string, options ...Option) (T, error) {
	fixtureOptions := &FixtureOptions{}
	for _, option := range options {
		option(fixtureOptions)
	}

	obj := fixtures[name]
	fixture, ok := obj.(*Fixture[T])
	if !ok {
		return emptyValue[T](), fmt.Errorf("fixture \"%s\": conflict between configured and requested types", name)
	}

	fixture.mutex.Lock()
	defer fixture.mutex.Unlock()

	if !fixtureOptions.DisableCache && fixture.value != emptyValue[T]() {
		return fixture.value, nil
	}

	value, err := fixture.generator(*fixtureOptions)
	if err != nil {
		return emptyValue[T](), fmt.Errorf("fixture \"%s\": could not get value from generator: %w", name, err)
	}

	if !fixtureOptions.DisableCache {
		fixture.value = value
	}

	fixtures[name] = fixture
	return value, nil
}
