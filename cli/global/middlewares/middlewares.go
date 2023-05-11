package global_middlewares

import "github.com/spf13/cobra"

type RunFn func(cmd *cobra.Command, args []string) (string, error)
type Middleware func(RunFn) RunFn
type CobraFn func(cmd *cobra.Command, args []string)
type CobraFnWrapper func(run RunFn) CobraFn
type NoopCobraFnWrapper func() CobraFn

func ComposeRun(middlewares ...Middleware) CobraFnWrapper {
	return func(run RunFn) CobraFn {
		return func(cmd *cobra.Command, args []string) {
			runMiddleware := ComposeMiddleware(middlewares...)(run)

			runMiddleware(cmd, args)
		}
	}
}

func ComposeNoopRun(middlewares ...Middleware) NoopCobraFnWrapper {
	return func() CobraFn {
		return ComposeRun(middlewares...)(NoopRun)
	}
}

func NoopRun(cmd *cobra.Command, args []string) (string, error) {
	return "", nil
}

func ComposeMiddleware(middlewares ...Middleware) Middleware {
	return func(runFn RunFn) RunFn {
		for _, middleware := range middlewares {
			runFn = middleware(runFn)
		}

		return runFn
	}
}
