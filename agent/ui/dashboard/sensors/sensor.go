package sensors

import (
	"fmt"

	"github.com/fluidtruck/deepcopy"
)

type Sensor interface {
	On(string, func(Event))
	Emit(string, interface{})
}

type Event struct {
	Name string
	data interface{}
}

func (e *Event) Unmarshal(target interface{}) error {
	err := deepcopy.DeepCopy(e.data, target)
	if err != nil {
		return fmt.Errorf("could not unmarshal event into target: %w", err)
	}

	return nil
}

type sensor struct {
	listeners map[string][]func(Event)
}

func NewSensor() Sensor {
	return &sensor{
		listeners: make(map[string][]func(Event)),
	}
}

func (r *sensor) On(eventName string, cb func(Event)) {
	var slice []func(Event)
	if existingSlice, ok := r.listeners[eventName]; ok {
		slice = existingSlice
	} else {
		slice = make([]func(Event), 0)
		slice = append(slice, cb)
	}
	r.listeners[eventName] = append(slice, cb)
}

func (r *sensor) Emit(eventName string, event interface{}) {
	listeners := r.listeners[eventName]
	e := Event{
		Name: eventName,
		data: event,
	}

	for _, listener := range listeners {
		listener(e)
	}
}
