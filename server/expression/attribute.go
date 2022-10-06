package expression

import "strings"

type Attribute string

func (a *Attribute) Capture(in []string) error {
	input := in[0]
	input = strings.TrimPrefix(input, "attr:")

	*a = Attribute(input)
	return nil
}
