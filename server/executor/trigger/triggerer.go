package trigger

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/trace"
)

type Triggerer interface {
	Trigger(context.Context, model.Test, trace.TraceID, trace.SpanID) (Response, error)
	Type() string
}

type Response struct {
	SpanAttributes map[string]string
	Response       any
}

func NewRegsitry(tracer trace.Tracer) *Registry {
	return &Registry{
		tracer: tracer,
		reg:    map[string]Triggerer{},
	}
}

type Registry struct {
	sync.Mutex
	tracer trace.Tracer
	reg    map[string]Triggerer
}

func (r *Registry) Add(t Triggerer) {
	r.Lock()
	defer r.Unlock()

	r.reg[t.Type()] = t
}

var ErrTriggererTypeNotRegistered = errors.New("triggerer type not found")

func (r *Registry) Get(triggererType string) (Triggerer, error) {
	r.Lock()
	defer r.Unlock()

	t, found := r.reg[triggererType]
	if !found {
		return nil, fmt.Errorf(`cannot get trigger type "%s": %w`, triggererType, ErrTriggererTypeNotRegistered)
	}

	return Instrument(r.tracer, t), nil
}
