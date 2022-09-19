package formatters

import "golang.org/x/exp/slices"

type Output string

var (
	CurrentOutput = DefaultOutput

	Outputs = []Output{
		Pretty,
		JSON,
	}

	DefaultOutput = Pretty

	Pretty Output = "pretty"
	JSON   Output = "json"
)

func SetOutput(o Output) {
	CurrentOutput = o
}

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
