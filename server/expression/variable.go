package expression

import "strings"

type Variable struct {
	name string
}

func NewVariable(name string) Variable {
	return Variable{name}
}

func (v *Variable) Capture(in []string) error {
	input := in[0]
	input = strings.TrimPrefix(input, "var:")

	v.name = input
	return nil
}

func (v *Variable) Name() string {
	return v.name
}
