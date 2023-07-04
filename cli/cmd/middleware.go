package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/spf13/cobra"
)

type RunFn func(cmd *cobra.Command, args []string) (string, error)
type CobraRunFn func(cmd *cobra.Command, args []string)
type MiddlewareWrapper func(RunFn) RunFn

func WithResultHandler(runFn RunFn) CobraRunFn {
	return func(cmd *cobra.Command, args []string) {
		res, err := runFn(cmd, args)

		if err != nil {
			fmt.Fprintf(os.Stderr, `
Version
%s

An error ocurred when executing the command

%s
`, versionText, err.Error())
			ExitCLI(1)
			return
		}

		if res != "" {
			fmt.Println(res)
		}
	}
}

func WithParamsHandler(validators ...Validator) MiddlewareWrapper {
	return func(runFn RunFn) RunFn {
		return func(cmd *cobra.Command, args []string) (string, error) {
			errors := make([]error, 0)

			for _, validator := range validators {
				errors = append(errors, validator.Validate(cmd, args)...)
			}

			if len(errors) > 0 {
				errorText := "The following errors occurred when validating the flags:\n"
				for _, err := range errors {
					errorText += err.Error() + "\n"
				}

				return "", fmt.Errorf(errorText)
			}

			return runFn(cmd, args)
		}
	}
}

type Validator interface {
	Validate(cmd *cobra.Command, args []string) []error
}

func WithResourceMiddleware(runFn RunFn, params ...Validator) CobraRunFn {
	params = append(params, resourceParams)
	return WithResultHandler(WithParamsHandler(params...)(runFn))
}

type resourceParameters struct {
	ResourceName string
}

func (p *resourceParameters) Validate(cmd *cobra.Command, args []string) []error {
	if len(args) == 0 || args[0] == "" {
		return []error{
			ParamError{
				Parameter: "resource",
				Message:   "resource name must be provided",
			},
		}
	}

	p.ResourceName = args[0]

	_, err := resources.Get(p.ResourceName)
	if errors.Is(err, resourcemanager.ErrResourceNotFound) {
		return []error{
			ParamError{
				Parameter: "resource",
				Message:   fmt.Sprintf("resource must be %s", resourceList()),
			},
		}
	}

	return nil
}
