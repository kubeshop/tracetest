package expression

import "strings"

const (
	metaPrefix    = "tracetest.selected_spans."
	metaPrefixLen = len(metaPrefix)
)

type Attribute struct {
	name string
}

func NewAttribute(name string) Attribute {
	return Attribute{name: name}
}

func (a *Attribute) Capture(in []string) error {
	input := in[0]
	input = strings.TrimPrefix(input, "attr:")

	a.name = input
	return nil
}

func (a *Attribute) Name() string {
	if a.IsMeta() {
		return a.name[metaPrefixLen:]
	}

	return a.name
}

func (a *Attribute) IsMeta() bool {
	return len(a.name) > metaPrefixLen && a.name[0:metaPrefixLen] == metaPrefix
}
