package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/ui"

	"github.com/spf13/cobra"
)

type RunFn func(ctx context.Context, cmd *cobra.Command, args []string) (string, error)
type CobraRunFn func(cmd *cobra.Command, args []string)
type MiddlewareWrapper func(RunFn) RunFn

func rootCtx(cmd *cobra.Command) context.Context {
	// cobra does not correctly progpagate rootcmd context to sub commands,
	// so we need to manually traverse the command tree to find the root context
	if cmd == nil {
		return nil
	}

	var (
		ctx = cmd.Context()
		p   = cmd.Parent()
	)
	if cmd.Parent() == nil {
		return ctx
	}
	for {
		ctx = p.Context()
		p = p.Parent()
		if p == nil {
			break
		}
	}
	return ctx
}

func WithResultHandler(runFn RunFn) CobraRunFn {
	return func(cmd *cobra.Command, args []string) {
		// we need the root cmd context in case of an error caused rerun
		ctx := rootCtx(cmd)

		res, err := runFn(ctx, cmd, args)

		if err != nil {
			handleError(ctx, err)
			return
		}

		if res != "" {
			fmt.Println(res)
		}
	}
}

func handleError(ctx context.Context, err error) {
	reqErr := resourcemanager.RequestError{}
	if errors.As(err, &reqErr) && reqErr.IsAuthError {
		handleAuthError(ctx)
	} else {
		OnError(err)
	}
}

func handleAuthError(ctx context.Context) {
	ui.DefaultUI.Warning("Your authentication token has expired, please log in again.")
	configurator.
		WithOnFinish(func(ctx context.Context, _ config.Config) {
			retryCommand(ctx)
		}).
		ExecuteUserLogin(ctx, cliConfig, nil)
}

func retryCommand(ctx context.Context) {
	handleRootExecErr(rootCmd.ExecuteContext(ctx))
}

type errorMessager interface {
	Message() string
}

const defaultErrorFormat = `
Version
%s

An error ocurred when executing the command

%s
`

func OnError(err error) {
	if renderer := findErrorMessageRenderer(err); renderer != nil {
		ui.DefaultUI.Error(renderer.Message())
	} else {
		errorMessage := handleErrorMessage(err)
		fmt.Fprintf(os.Stderr, defaultErrorFormat, versionText, errorMessage)
	}
	ExitCLI(1)
}

func findErrorMessageRenderer(err error) errorMessager {
	for err != nil {
		if renderer, ok := err.(errorMessager); ok {
			return renderer
		}
		err = errors.Unwrap(err)
	}
	return nil
}

func handleErrorMessage(err error) string {
	var requestError resourcemanager.RequestError
	hasRequestError := errors.As(err, &requestError)

	if !hasRequestError || requestError.Code != 401 {
		return err.Error()
	}

	return fmt.Sprintf("user is not authenticated on %s", cliConfig.Endpoint)
}

func WithParamsHandler(validators ...cmdutil.Validator) MiddlewareWrapper {
	return func(runFn RunFn) RunFn {
		return func(ctx context.Context, cmd *cobra.Command, args []string) (string, error) {
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

			return runFn(ctx, cmd, args)
		}
	}
}

func WithResourceMiddleware(runFn RunFn, params ...cmdutil.Validator) CobraRunFn {
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
			cmdutil.ParamError{
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
				cmdutil.ParamError{
					Parameter: "resource",
					Message:   fmt.Sprintf("resource \"%s\" not found. Did you mean this?\n\t%s", p.ResourceName, suggestion),
				},
			}
		}

		return []error{
			cmdutil.ParamError{
				Parameter: "resource",
				Message:   fmt.Sprintf("resource must be %s", resourceList()),
			},
		}
	}

	return nil
}
