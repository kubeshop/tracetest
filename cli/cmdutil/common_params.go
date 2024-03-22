package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Validator interface {
	Validate(cmd *cobra.Command, args []string) []error
}

type ParamError struct {
	Parameter string
	Message   string
}

func (pe ParamError) Error() string {
	return fmt.Sprintf(`[%s] %s`, pe.Parameter, pe.Message)
}
