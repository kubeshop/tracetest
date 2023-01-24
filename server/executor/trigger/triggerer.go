package trigger

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/trace"
)

type TriggerOptions struct {
	TraceID  trace.TraceID
	Executor expression.Executor
}

type Triggerer interface {
	Trigger(context.Context, model.Test, *TriggerOptions) (Response, error)
	Type() model.TriggerType
	Resolve(context.Context, model.Test, *TriggerOptions) (model.Test, error)
}

type Response struct {
	SpanAttributes map[string]string
	Result         model.TriggerResult
	TraceID        trace.TraceID
	SpanID         trace.SpanID
}

func NewRegsitry(tracer, triggerSpanTracer trace.Tracer) *Registry {
	return &Registry{
		tracer:            tracer,
		triggerSpanTracer: triggerSpanTracer,
		reg:               map[model.TriggerType]Triggerer{},
	}
}

type Registry struct {
	sync.Mutex
	tracer            trace.Tracer
	triggerSpanTracer trace.Tracer
	reg               map[model.TriggerType]Triggerer
}

func (r *Registry) Add(t Triggerer) {
	r.Lock()
	defer r.Unlock()

	r.reg[t.Type()] = t
}

var ErrTriggererTypeNotRegistered = errors.New("triggerer type not found")

func (r *Registry) Get(triggererType model.TriggerType) (Triggerer, error) {
	r.Lock()
	defer r.Unlock()

	t, found := r.reg[triggererType]
	if !found {
		return nil, fmt.Errorf(`cannot get trigger type "%s": %w`, triggererType, ErrTriggererTypeNotRegistered)
	}

	return Instrument(r.tracer, r.triggerSpanTracer, t), nil
}

func WrapInQuotes(input string, quoteChar string) string {
	return fmt.Sprintf("%s%s%s", quoteChar, input, quoteChar)
}
