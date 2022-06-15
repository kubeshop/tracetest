package subscription

type Message struct {
	ResourceID string
	Type       string
	Content    interface{}
}
