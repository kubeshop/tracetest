package formatters

import "golang.org/x/exp/slices"

type Output string

var (
	Outputs = []Output{
		Pretty,
		JSON,
		YAML,
	}
	Pretty Output = "pretty"
	JSON   Output = "json"
	YAML   Output = "yaml"
)

func OuputsStr() []string {
	out := make([]string, len(Outputs))
	for i, o := range Outputs {
		out[i] = string(o)
	}

	return out
}

func ValidOutput(o Output) bool {
	return slices.Contains(Outputs, o)
}
