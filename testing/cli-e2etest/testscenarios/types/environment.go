package types

// Note: these types are very similar to the types on the server folder
// however they are defined here to avoid bias with the current implementation

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
