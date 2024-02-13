package sensors

import (
	"fmt"

	"github.com/fluidtruck/deepcopy"
)

type Sensor interface {
	On(string, func(Event))
	Emit(string, interface{})
	Reset()
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
	lastEvent map[string]Event
}

func NewSensor() Sensor {
	return &sensor{
		listeners: make(map[string][]func(Event)),
		lastEvent: make(map[string]Event),
	}
}

func (r *sensor) Reset() {
	r.listeners = make(map[string][]func(Event))
	r.lastEvent = make(map[string]Event)
}

func (r *sensor) On(eventName string, cb func(Event)) {
	var slice []func(Event)
	if existingSlice, ok := r.listeners[eventName]; ok {
		slice = existingSlice
	} else {
		slice = make([]func(Event), 0)
	}
	r.listeners[eventName] = append(slice, cb)

	if event, ok := r.lastEvent[eventName]; ok {
		cb(event)
	}
}

func (r *sensor) Emit(eventName string, event interface{}) {
	listeners := r.listeners[eventName]
	e := Event{
		Name: eventName,
		data: event,
	}

	r.lastEvent[eventName] = e

	for _, listener := range listeners {
		listener(e)
	}
}
