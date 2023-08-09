package types

// Note: these types are very similar to the types on the server folder
// however they are defined here to avoid bias with the current implementation

type TestSuite struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
}

type TestSuiteResource struct {
	Type string    `json:"type"`
	Spec TestSuite `json:"spec"`
}

type AugmentedTestSuiteLastRun struct {
	Passes int `json:"passes"`
	Fails  int `json:"fails"`
}

type AugmentedTestSuiteSummary struct {
	Runs    int                       `json:"runs"`
	LastRun AugmentedTestSuiteLastRun `json:"lastRun"`
}

type AugmentedTestSuite struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Steps       []string                  `json:"steps"`
	Summary     AugmentedTestSuiteSummary `json:"summary"`
}

type AugmentedTestSuiteResource struct {
	Type string             `json:"type"`
	Spec AugmentedTestSuite `json:"spec"`
}
