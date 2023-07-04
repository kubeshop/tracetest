package analyzer

type Result struct {
	SpanID string  `json:"span_id"`
	Passed bool    `json:"passed"`
	Errors []Error `json:"errors"`
}

type Error struct {
	Value       string   `json:"value"`
	Expected    string   `json:"expected"`
	Description string   `json:"description"`
	Suggestions []string `json:"suggestions"`
}

type EvalRuleResult struct {
	Passed  bool     `json:"passed"`
	Results []Result `json:"results"`
}

type RuleResult struct {
	// config
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	ErrorDescription string   `json:"errorDescription"`
	Tips             []string `json:"tips"`
	Weight           int      `json:"weight"`
	Level            string   `json:"level"`

	// results
	Passed  bool     `json:"passed"`
	Results []Result `json:"results"`
}

func NewRuleResult(config LinterRule, results EvalRuleResult) RuleResult {
	return RuleResult{
		// config
		Id:               config.Id,
		Name:             config.Name,
		Description:      config.Description,
		ErrorDescription: config.ErrorDescription,
		Tips:             config.Tips,
		Weight:           config.Weight,
		Level:            config.ErrorLevel,

		// results
		Passed:  results.Passed,
		Results: results.Results,
	}
}

type LinterResult struct {
	Plugins      []PluginResult `json:"plugins"`
	Score        int            `json:"score"`
	MinimumScore int            `json:"minimumScore"`
	Passed       bool           `json:"passed"`
}

func NewLinterResult(pluginResults []PluginResult, totalScore int, passed bool) LinterResult {
	// if no plugins are enabled, then the score should be 100
	score := 100
	if len(pluginResults) > 0 {
		score = totalScore / len(pluginResults)
	}

	return LinterResult{
		Plugins: pluginResults,
		Score:   score,
		Passed:  passed,
	}
}

type PluginResult struct {
	// metadata
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	//results
	Passed bool         `json:"passed"`
	Score  int          `json:"score"`
	Rules  []RuleResult `json:"rules"`
}

func (pr PluginResult) CalculateResults() PluginResult {
	totalScore := 0
	passScore := 0
	passed := true

	for _, result := range pr.Rules {
		// only error level rules are taken into account
		if result.Level == ErrorLevelError {
			totalScore += result.Weight

			if result.Passed {
				passScore += result.Weight
			} else {
				passed = false
			}
		}
	}

	pr.Score = 0
	if totalScore > 0 {
		pr.Score = (100 * passScore) / totalScore
	}
	pr.Passed = passed
	return pr
}
