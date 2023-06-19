package types

// Note: these types are very similar to the types on the server folder
// however they are defined here to avoid bias with the current implementation

type DataStore struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

type DataStoreResource struct {
	Type string    `json:"type"`
	Spec DataStore `json:"spec"`
}
