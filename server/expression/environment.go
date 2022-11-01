package expression

import "strings"

type Environment struct {
	name string
}

func NewEnvironment(name string) Variable {
	return Variable{name}
}

func (e *Environment) Capture(in []string) error {
	input := in[0]
	input = strings.TrimPrefix(input, "env:")

	e.name = input
	return nil
}

func (e *Environment) Name() string {
	return e.name
}
