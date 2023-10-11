package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type RunFn func(cmd *cobra.Command, args []string) (string, error)
type CobraRunFn func(cmd *cobra.Command, args []string)
type MiddlewareWrapper func(RunFn) RunFn

func WithResultHandler(runFn RunFn) CobraRunFn {
	return func(cmd *cobra.Command, args []string) {
		res, err := runFn(cmd, args)

		if err != nil {
			OnError(err)
			return
		}

		if res != "" {
			fmt.Println(res)
		}
	}
}

func OnError(err error) {
	fmt.Fprintf(os.Stderr, `
Version
%s

An error ocurred when executing the command

%s
`, versionText, err.Error())
	ExitCLI(1)
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
	optional     bool
}

func (p *resourceParameters) Validate(cmd *cobra.Command, args []string) []error {
	// if the resourceName is optional, skip validation.
	if p.optional {
		// we still need to bind it to the struct in case the user provided a value.
		// we need to check the args has at least one element to avoid a panic.
		if len(args) > 0 {
			p.ResourceName = args[0]
		}
		return nil
	}

	if len(args) == 0 || args[0] == "" {
		return []error{
			paramError{
				Parameter: "resource",
				Message:   fmt.Sprintf("resource name must be provided. Available resources: %s", resourceList()),
			},
		}
	}

	p.ResourceName = args[0]

	exists := resources.Exists(p.ResourceName)
	if !exists {
		suggestion := resources.Suggest(p.ResourceName)
		if suggestion != "" {
			return []error{
				paramError{
					Parameter: "resource",
					Message:   fmt.Sprintf("resource \"%s\" not found. Did you mean this?\n\t%s", p.ResourceName, suggestion),
				},
			}
		}

		return []error{
			paramError{
				Parameter: "resource",
				Message:   fmt.Sprintf("resource must be %s", resourceList()),
			},
		}
	}

	return nil
}

type paramError struct {
	Parameter string
	Message   string
}

func (pe paramError) Error() string {
	return fmt.Sprintf(`[%s] %s`, pe.Parameter, pe.Message)
}
