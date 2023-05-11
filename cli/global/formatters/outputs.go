package global_formatters

import (
	"strings"

	"golang.org/x/exp/slices"
)

type Output string

var (
	OutputFormats       = OuputsStr()
	OutputFormatsString = strings.Join(OutputFormats, "|")

	CurrentOutput = DefaultOutput

	Outputs = []Output{
		Pretty,
		JSON,
		YAML,
		Empty,
	}

	DefaultOutput = Pretty

	Empty  Output = ""
	Pretty Output = "pretty"
	JSON   Output = "json"
	YAML   Output = "yaml"
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
