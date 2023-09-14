package trigger

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"go.opentelemetry.io/otel/trace"
)

type Options struct {
	TraceID trace.TraceID
	SpanID  trace.SpanID
	TestID  id.ID
}

type Triggerer interface {
	Trigger(context.Context, trigger.Trigger, *Options) (Response, error)
	Type() trigger.TriggerType
}

type Response struct {
	SpanAttributes map[string]string
	Result         trigger.TriggerResult
	TraceID        trace.TraceID
	SpanID         trace.SpanID
}

func NewRegistry(tracer trace.Tracer) *Registry {
	return &Registry{
		tracer: tracer,
		reg:    map[trigger.TriggerType]Triggerer{},
	}
}

type Registry struct {
	sync.Mutex
	tracer trace.Tracer
	reg    map[trigger.TriggerType]Triggerer
}

func (r *Registry) Add(t Triggerer) {
	r.Lock()
	defer r.Unlock()

	r.reg[t.Type()] = t
}

var ErrTriggererTypeNotRegistered = errors.New("triggerer type not found")

func (r *Registry) Get(triggererType trigger.TriggerType) (Triggerer, error) {
	r.Lock()
	defer r.Unlock()

	t, found := r.reg[triggererType]
	if !found {
		return nil, fmt.Errorf(`cannot get trigger type "%s": %w`, triggererType, ErrTriggererTypeNotRegistered)
	}

	return Instrument(r.tracer, t), nil
}
