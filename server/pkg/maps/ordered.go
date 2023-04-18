package maps

import (
	"encoding/json"
	"errors"
)

type Ordered[K comparable, V any] struct {
	list        []V
	keyPosition map[K]int
	positionKey map[int]K
}

func (om *Ordered[K, V]) replace(om2 *Ordered[K, V]) {
	*om = *om2
}

type jsonOrderedMapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

func (om Ordered[K, V]) MarshalJSON() ([]byte, error) {
	j := []jsonOrderedMapEntry[K, V]{}
	om.ForEach(func(key K, asserts V) error {
		j = append(j, jsonOrderedMapEntry[K, V]{key, asserts})
		return nil
	})

	return json.Marshal(j)
}

func (om *Ordered[K, V]) UnmarshalJSON(data []byte) error {
	aux := []jsonOrderedMapEntry[K, V]{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	newMap := Ordered[K, V]{}
	var err error
	for _, s := range aux {
		newMap, err = newMap.Add(s.Key, s.Value)
		if err != nil {
			return err
		}
	}

	om.replace(&newMap)

	return nil
}

func (om Ordered[K, V]) MustAdd(key K, asserts V) Ordered[K, V] {
	def, err := om.Add(key, asserts)
	if err != nil {
		panic(err)
	}
	return def
}

func (om Ordered[K, V]) Add(key K, asserts V) (Ordered[K, V], error) {
	if om.keyPosition == nil {
		om.keyPosition = make(map[K]int)
	}
	if om.positionKey == nil {
		om.positionKey = make(map[int]K)
	}

	if _, exists := om.keyPosition[key]; exists {
		return Ordered[K, V]{}, errors.New("selector already exists")
	}

	om.list = append(om.list, asserts)
	ix := len(om.list) - 1
	om.keyPosition[key] = ix
	om.positionKey[ix] = key

	return om, nil
}

func (om Ordered[K, V]) Len() int {
	return len(om.list)
}

func (om Ordered[K, V]) Get(key K) V {
	ix, exists := om.keyPosition[key]
	if !exists {
		var result V
		return result
	}

	return om.list[ix]
}

func (om *Ordered[K, V]) ForEach(fn func(key K, val V) error) error {
	for ix, asserts := range om.list {
		K := om.positionKey[ix]
		err := fn(K, asserts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (om Ordered[K, V]) Unordered() map[K]V {
	m := map[K]V{}
	om.ForEach(func(key K, val V) error {
		m[key] = val
		return nil
	})

	return m
}
