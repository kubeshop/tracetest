package types

// Note: these types are very similar to the types on the server folder
// however they are defined here to avoid bias with the current implementation

type AnalyzerResource struct {
	Type string   `json:"type"`
	Spec Analyzer `json:"spec"`
}

type Analyzer struct {
	Id           string           `json:"id"`
	Name         string           `json:"name"`
	Enabled      bool             `json:"enabled"`
	MinimumScore int              `json:"minimumScore"`
	Plugins      []AnalyzerPlugin `json:"plugins"`
}

type AnalyzerPlugin struct {
	Name     string `json:"name"`
	Enabled  bool   `json:"enabled"`
	Required bool   `json:"required"`
}
