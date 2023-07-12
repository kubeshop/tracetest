package types

// Note: these types are very similar to the types on the server folder
// however they are defined here to avoid bias with the current implementation

type AnalyzerResource struct {
	Type string   `json:"type"`
	Spec Analyzer `json:"spec"`
}

type Analyzer struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Enabled      bool             `json:"enabled"`
	MinimumScore int              `json:"minimumScore"`
	Plugins      []AnalyzerPlugin `json:"plugins"`
}

type AnalyzerPlugin struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`

	Rules []AnalyzerRule `json:"rules"`
}

type AnalyzerRule struct {
	ID               string   `json:"id"`
	Weight           int      `json:"weight"`
	ErrorLevel       string   `json:"errorLevel"`
	ErrorDescription string   `json:"errorDescription"`
	Description      string   `json:"description"`
	Tips             []string `json:"tips"`
	Name             string   `json:"name"`
	Enabled          bool     `json:"enabled"`
}
