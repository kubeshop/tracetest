package websocket

type Message struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

type Event struct {
	Type  string      `json:"type"`
	Event interface{} `json:"event"`
}

func SubscriptionSuccess(subscriptionId string) Message {
	return Message{
		Type: "success",
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
	return Event{
		Type:  "update",
		Event: resource,
	}
}
