package websocket

import "encoding/json"

type Message struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

func SuccessMessage(messageType string) Message {
	return Message{
		Type:    messageType,
		Message: "success",
	}
}

func Error(err error) Message {
	return Message{
		Type:    "error",
		Message: err.Error(),
	}
}

func UpdateMessage(object interface{}) Message {
	jsonContent, _ := json.Marshal(object)
	return Message{
		Type:    "update",
		Message: jsonContent,
	}
}
