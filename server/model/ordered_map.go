package model

import (
	"encoding/json"
	"errors"
)

type OrderedMap[K comparable, V any] struct {
	list        []V
	keyPosition map[K]int
	positionKey map[int]K
}

func (om *OrderedMap[K, V]) replace(om2 *OrderedMap[K, V]) {
	*om = *om2
}

type jsonOrderedMapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

func (om OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	j := []jsonOrderedMapEntry[K, V]{}
	om.Map(func(key K, asserts V) {
		j = append(j, jsonOrderedMapEntry[K, V]{key, asserts})
	})

	return json.Marshal(j)
}

func (om *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	aux := []jsonOrderedMapEntry[K, V]{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	newMap := OrderedMap[K, V]{}
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

func (om OrderedMap[K, V]) MustAdd(key K, asserts V) OrderedMap[K, V] {
	def, err := om.Add(key, asserts)
	if err != nil {
		panic(err)
	}
	return def
}

func (om OrderedMap[K, V]) Add(key K, asserts V) (OrderedMap[K, V], error) {
	if om.keyPosition == nil {
		om.keyPosition = make(map[K]int)
	}
	if om.positionKey == nil {
		om.positionKey = make(map[int]K)
	}

	if _, exists := om.keyPosition[key]; exists {
		return OrderedMap[K, V]{}, errors.New("selector already exists")
	}

	om.list = append(om.list, asserts)
	ix := len(om.list) - 1
	om.keyPosition[key] = ix
	om.positionKey[ix] = key

	return om, nil
}

func (om OrderedMap[K, V]) Len() int {
	return len(om.list)
}

func (om OrderedMap[K, V]) Get(key K) V {
	ix, exists := om.keyPosition[key]
	if !exists {
		var result V
		return result
	}

	return om.list[ix]
}

func (om *OrderedMap[K, V]) Map(fn func(key K, val V)) {
	for ix, asserts := range om.list {
		K := om.positionKey[ix]
		fn(K, asserts)
	}
}
