package types

// Note: these types are very similar to the types on the server folder
// however they are defined here to avoid bias with the current implementation

type PollingProfilePeriodicStrategy struct {
	Timeout              string `json:"timeout"`
	RetryDelay           string `json:"retryDelay"`
	SelectorMatchRetries string `json:"selectorMatchRetries"`
}

type PollingProfile struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default"`

	Strategy string                         `json:"strategy"`
	Periodic PollingProfilePeriodicStrategy `json:"periodic"`
}

type PollingProfileResource struct {
	Type string         `json:"type"`
	Spec PollingProfile `json:"spec"`
}
