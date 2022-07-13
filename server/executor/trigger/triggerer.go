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
	Type() model.TriggerType
}

type Response struct {
	SpanAttributes map[string]string
	Result         model.TriggerResult
}

func NewRegsitry(tracer trace.Tracer) *Registry {
	return &Registry{
		tracer: tracer,
		reg:    map[model.TriggerType]Triggerer{},
	}
}

type Registry struct {
	sync.Mutex
	tracer trace.Tracer
	reg    map[model.TriggerType]Triggerer
}

func (r *Registry) Add(t Triggerer) {
	r.Lock()
	defer r.Unlock()

	fmt.Println("*****", t.Type())
	r.reg[t.Type()] = t
}

var ErrTriggererTypeNotRegistered = errors.New("triggerer type not found")

func (r *Registry) Get(triggererType model.TriggerType) (Triggerer, error) {
	r.Lock()
	defer r.Unlock()

	t, found := r.reg[triggererType]
	fmt.Println("*****", triggererType, found)
	if !found {
		return nil, fmt.Errorf(`cannot get trigger type "%s": %w`, triggererType, ErrTriggererTypeNotRegistered)
	}

	return Instrument(r.tracer, t), nil
}
