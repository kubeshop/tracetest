package cmd

import "fmt"

type ParamError struct {
	Parameter string
	Message   string
}

func (pe ParamError) Error() string {
	return fmt.Sprintf(`[%s] %s`, pe.Parameter, pe.Message)
}
