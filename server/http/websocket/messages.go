package websocket

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/http"
	"github.com/kubeshop/tracetest/server/model"
)

type Message struct {
	Type     string      `json:"type"`
	Resource string      `json:"resource"`
	Message  interface{} `json:"message"`
}

type Event struct {
	Type     string      `json:"type"`
	Resource string      `json:"resource"`
	Event    interface{} `json:"event"`
}

func SubscriptionSuccess(resource, subscriptionId string) Message {
	return Message{
		Type:     "success",
		Resource: resource,
		Message: struct {
			SubscriptionId string `json:"subscriptionId"`
		}{SubscriptionId: subscriptionId},
	}
}

func UnsubscribeSuccess() Message {
	return Message{
		Type:    "success",
		Message: "ok",
	}
}

func ErrorMessage(err error) Message {
	return Message{
		Type:    "error",
		Message: err.Error(),
	}
}

func ResourceUpdatedEvent(resource interface{}) Event {
	var mapped interface{}
	switch v := resource.(type) {
	case model.Run:
		mapped = (http.OpenAPIMapper{}).Run(&v)
	case *model.Run:
		mapped = (http.OpenAPIMapper{}).Run(v)
	default:
		fmt.Printf("type %T mapping not supported\n", v)
		mapped = v
	}

	return Event{
		Type:  "update",
		Event: mapped,
	}
}
