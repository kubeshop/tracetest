package parameters

import "github.com/spf13/cobra"

type ParamError struct {
	Parameter string
	Message   string
}

type Params interface {
	Validate(cmd *cobra.Command, args []string) []ParamError
}

func ValidateParams(cmd *cobra.Command, args []string, params ...Params) []ParamError {
	errors := make([]ParamError, 0)

	for _, param := range params {
		errors = append(errors, param.Validate(cmd, args)...)
	}

	return errors
}
