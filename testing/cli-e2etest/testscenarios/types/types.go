package types

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

type Environment struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Values map[string]string `json:"values"`
}

type EnvironmentResource struct {
	Type string      `json:"type"`
	Spec Environment `json:"spec"`
}
