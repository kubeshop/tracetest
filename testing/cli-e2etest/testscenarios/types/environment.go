package types

// Note: these types are very similar to the types on the server folder
// however they are defined here to avoid bias with the current implementation

type VariableSetKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type VariableSet struct {
	ID     string                `json:"id"`
	Name   string                `json:"name"`
	Values []VariableSetKeyValue `json:"values"`
}

type VariableSetResource struct {
	Type string      `json:"type"`
	Spec VariableSet `json:"spec"`
}
