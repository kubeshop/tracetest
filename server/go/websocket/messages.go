package websocket

type Message struct {
	Type    string
	Message interface{}
}

func Error(err error) Message {
	return Message{
		Type:    "error",
		Message: err.Error(),
	}
}
