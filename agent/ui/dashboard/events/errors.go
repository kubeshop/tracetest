package events

type MessageType string

var (
	Warning MessageType = "warning"
	Error   MessageType = "error"
)
