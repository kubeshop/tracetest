package subscription

// Subscription represents a group of subscribers that are waiting for updates of a specific resource
type Subscription struct {
	resourceID  string
	subscribers []Subscriber
}
