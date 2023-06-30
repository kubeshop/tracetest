package results

type LinterResult struct {
	Plugins      []PluginResult `json:"plugins"`
	Score        int            `json:"score"`
	MinimumScore int            `json:"minimumScore"`
	Passed       bool           `json:"passed"`
}

func NewLinterResult(pluginResults []PluginResult, totalScore int, passed bool) LinterResult {
	return LinterResult{
		Plugins: pluginResults,
		Score:   totalScore / len(pluginResults),
		Passed:  passed,
	}
}
