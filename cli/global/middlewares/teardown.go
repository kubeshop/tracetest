package global_middlewares

import (
	"github.com/kubeshop/tracetest/cli/analytics"
	global_setup "github.com/kubeshop/tracetest/cli/global/setup"
	"github.com/spf13/cobra"
)

func WithTeardownMiddleware(setup *global_setup.Setup) Middleware {
	return func(runFn RunFn) RunFn {
		return func(cmd *cobra.Command, args []string) (string, error) {
			defer setup.Logger.Sync()
			defer analytics.Close()

			return runFn(cmd, args)
		}
	}
}
