package types

// Note: these types are very similar to the types on the server folder
// however they are defined here to avoid bias with the current implementation

// DataStore
type DataStore struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

type DataStoreResource struct {
	Type string    `json:"type"`
	Spec DataStore `json:"spec"`
}

// Environment

type EnvironmentKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Environment struct {
	ID     string                `json:"id"`
	Name   string                `json:"name"`
	Values []EnvironmentKeyValue `json:"values"`
}

type EnvironmentResource struct {
	Type string      `json:"type"`
	Spec Environment `json:"spec"`
}

// Config

type Config struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	AnalyticsEnabled bool `json:"analyticsEnabled"`
}

type ConfigResource struct {
	Type string `json:"type"`
	Spec Config `json:"spec"`
}

// PollingProfile

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
