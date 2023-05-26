package cmd

import (
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/parameters"
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

func WithParamsHandler(params ...parameters.Params) MiddlewareWrapper {
	return func(runFn RunFn) RunFn {
		return func(cmd *cobra.Command, args []string) (string, error) {
			errors := parameters.ValidateParams(cmd, args, params...)

			if len(errors) > 0 {
				errorText := `The following errors occurred when validating the flags:`
				for _, err := range errors {
					errorText += fmt.Sprintf(`
[%s] %s`, err.Parameter, err.Message)
				}

				return "", fmt.Errorf(errorText)
			}

			return runFn(cmd, args)
		}
	}
}

func WithResourceMiddleware(runFn RunFn, params ...parameters.Params) CobraRunFn {
	params = append(params, resourceParams)
	return WithResultHandler(WithParamsHandler(params...)(runFn))
}
