package types

// Note: these types are very similar to the types on the server folder
// however they are defined here to avoid bias with the current implementation

type Transaction struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
}

type TransactionResource struct {
	Type string      `json:"type"`
	Spec Transaction `json:"spec"`
}
