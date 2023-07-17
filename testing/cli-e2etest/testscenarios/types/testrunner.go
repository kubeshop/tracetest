package types

type TestRunnerResource struct {
	Type string     `json:"type"`
	Spec TestRunner `json:"spec"`
}

type TestRunner struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	RequiredGates []string `json:"requiredGates"`
}
