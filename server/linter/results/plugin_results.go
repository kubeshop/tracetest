package results

type PluginResult struct {
	// metadata
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`

	//results
	Passed bool         `json:"passed"`
	Score  int          `json:"score"`
	Rules  []RuleResult `json:"rules"`
}

func (pr PluginResult) CalculateResults() PluginResult {
	failedScore := 0
	passed := true

	for _, result := range pr.Rules {
		if !result.Passed {
			passed = false
			failedScore += result.Weight
		}
	}

	pr.Score = 100 - failedScore
	pr.Passed = passed
	return pr
}
