package types

// Note: these types are very similar to the types on the server folder
// however they are defined here to avoid bias with the current implementation

type Config struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	AnalyticsEnabled bool `json:"analyticsEnabled"`
}

type ConfigResource struct {
	Type string `json:"type"`
	Spec Config `json:"spec"`
}
