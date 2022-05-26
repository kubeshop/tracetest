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

type Generator[T comparable] func(args ...interface{}) (T, error)

var fixtures = make(map[string]interface{})

func emptyValue[T comparable]() T {
	var empty T
	return empty
}

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

func GetFixtureValue[T comparable](name string, args ...interface{}) (T, error) {
	obj := fixtures[name]
	fixture, ok := obj.(*Fixture[T])
	if !ok {
		return emptyValue[T](), fmt.Errorf("fixture \"%s\": conflict between configured and requested types", name)
	}

	fixture.mutex.Lock()
	defer fixture.mutex.Unlock()

	if fixture.value != emptyValue[T]() {
		return fixture.value, nil
	}

	value, err := fixture.generator(args)
	if err != nil {
		return emptyValue[T](), fmt.Errorf("fixture \"%s\": could not get value from generator: %w", name, err)
	}

	fixture.value = value
	fixtures[name] = fixture
	return value, nil
}
