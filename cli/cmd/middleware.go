package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/ui"

	"github.com/spf13/cobra"
)

type RunFn func(ctx context.Context, cmd *cobra.Command, args []string) (string, error)
type CobraRunFn func(cmd *cobra.Command, args []string)
type MiddlewareWrapper func(RunFn) RunFn

// type keyCmd struct{}

// var cmdKey keyCmd

// func ContextWithCmd(cmd *cobra.Command, args []string) context.Context {
// 	return context.WithValue(
// 		cmd.Context(),
// 		cmdKey,
// 		buildCmdLine(cmd, args),
// 	)
// }

// func ContextGetCmd(ctx context.Context) string {
// 	v := ctx.Value(cmdKey)
// 	if v == nil {
// 		return ""
// 	}
// 	return v.(string)
// }

func rootCtx(cmd *cobra.Command) context.Context {
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
		ExecuteUserLogin(ctx, cliConfig)
}

func retryCommand(ctx context.Context) {
	handleRootExecErr(rootCmd.ExecuteContext(ctx))
}

// func buildCmdLine(cmd *cobra.Command, args []string) string {
// 	cmdLine := append([]string{cmd.CommandPath()}, args...)
// 	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
// 		if flag.Changed {
// 			cmdLine = append(cmdLine, fmt.Sprintf("--%s=%s", flag.Name, flag.Value.String()))
// 		}
// 	})
// 	return strings.Join(cmdLine, " ")
// }

// func exitWithCmdStatus(err error) {
// 	if exitError, ok := err.(*exec.ExitError); ok {
// 		ws := exitError.Sys().(syscall.WaitStatus)
// 		os.Exit(ws.ExitStatus())
// 	}
// }

type errorMessageRenderer interface {
	Render()
}

const defaultErrorFormat = `
Version
%s

An error ocurred when executing the command

%s
`

func OnError(err error) {
	errorMessage := handleErrorMessage(err)

	if renderer, ok := err.(errorMessageRenderer); ok {
		renderer.Render()
	} else {
		fmt.Fprintf(os.Stderr, defaultErrorFormat, versionText, errorMessage)
	}
	ExitCLI(1)
}

func handleErrorMessage(err error) string {
	var requestError resourcemanager.RequestError
	hasRequestError := errors.As(err, &requestError)

	if !hasRequestError || requestError.Code != 401 {
		return err.Error()
	}

	return fmt.Sprintf("user is not authenticated on %s", cliConfig.Endpoint)
}

func WithParamsHandler(validators ...Validator) MiddlewareWrapper {
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
