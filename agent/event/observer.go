package event

import (
	"github.com/kubeshop/tracetest/agent/proto"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type Observer interface {
	StartTriggerExecution(*proto.TriggerRequest)
	EndTriggerExecution(*proto.TriggerRequest, error)

	StartTracePoll(*proto.PollingRequest)
	EndTracePoll(*proto.PollingRequest, error)

	StartSpanReceive([]*v1.Span)
	EndSpanReceive([]*v1.Span, error)

	StartDataStoreConnection(*proto.DataStoreConnectionTestRequest)
	EndDataStoreConnection(*proto.DataStoreConnectionTestRequest, error)

	Error(error)
}

type wrapperObserver struct {
	wrappedObserver Observer
}

func NewNopObserver() Observer {
	return &wrapperObserver{
		wrappedObserver: nil,
	}
}

func WrapObserver(observer Observer) Observer {
	return &wrapperObserver{
		wrappedObserver: observer,
	}
}

var _ Observer = &wrapperObserver{}

func (o *wrapperObserver) StartDataStoreConnection(request *proto.DataStoreConnectionTestRequest) {
	if o.wrappedObserver == nil {
		return
	}

	o.wrappedObserver.StartDataStoreConnection(request)
}

func (o *wrapperObserver) EndDataStoreConnection(request *proto.DataStoreConnectionTestRequest, err error) {
	if o.wrappedObserver == nil {
		return
	}

	o.wrappedObserver.EndDataStoreConnection(request, err)
}

func (o *wrapperObserver) StartSpanReceive(spans []*v1.Span) {
	if o.wrappedObserver == nil {
		return
	}

	o.wrappedObserver.StartSpanReceive(spans)
}

func (o *wrapperObserver) EndSpanReceive(spans []*v1.Span, err error) {
	if o.wrappedObserver == nil {
		return
	}

	o.wrappedObserver.EndSpanReceive(spans, err)
}

func (o *wrapperObserver) Error(err error) {
	if o.wrappedObserver == nil {
		return
	}

	o.wrappedObserver.Error(err)
}

func (o *wrapperObserver) StartTracePoll(request *proto.PollingRequest) {
	if o.wrappedObserver == nil {
		return
	}

	o.wrappedObserver.StartTracePoll(request)
}

func (o *wrapperObserver) EndTracePoll(request *proto.PollingRequest, err error) {
	if o.wrappedObserver == nil {
		return
	}

	o.wrappedObserver.EndTracePoll(request, err)
}

func (o *wrapperObserver) StartTriggerExecution(request *proto.TriggerRequest) {
	if o.wrappedObserver == nil {
		return
	}

	o.wrappedObserver.StartTriggerExecution(request)
}

func (o *wrapperObserver) EndTriggerExecution(request *proto.TriggerRequest, err error) {
	if o.wrappedObserver == nil {
		return
	}

	o.wrappedObserver.EndTriggerExecution(request, err)
}
