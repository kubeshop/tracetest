package expression

import "strings"

type Duration string

func (a *Duration) Capture(in []string) error {
	input := in[0]
	input = strings.ReplaceAll(input, " ", "")

	*a = Duration(input)
	return nil
}
