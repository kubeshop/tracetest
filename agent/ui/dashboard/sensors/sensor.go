package sensors

import (
	"fmt"
	"sync"

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

	mutex sync.Mutex
}

func NewSensor() Sensor {
	return &sensor{
		listeners: make(map[string][]func(Event)),
		lastEvent: make(map[string]Event),
	}
}

func (r *sensor) Reset() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.listeners = make(map[string][]func(Event))
	r.lastEvent = make(map[string]Event)
}

func (r *sensor) On(eventName string, cb func(Event)) {
	r.mutex.Lock()

	var slice []func(Event)
	if existingSlice, ok := r.listeners[eventName]; ok {
		slice = existingSlice
	} else {
		slice = make([]func(Event), 0)
	}
	r.listeners[eventName] = append(slice, cb)
	r.mutex.Unlock()

	if event, ok := r.lastEvent[eventName]; ok {
		cb(event)
	}
}

func (r *sensor) Emit(eventName string, event interface{}) {
	r.mutex.Lock()
	r.mutex.Unlock()

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
