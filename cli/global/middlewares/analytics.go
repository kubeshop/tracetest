package global_middlewares

import (
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
)

func WithAnalytics(name, category string) Middleware {
	return func(runFn RunFn) RunFn {
		return func(cmd *cobra.Command, args []string) (string, error) {
			analytics.Track(name, category, map[string]string{})

			return runFn(cmd, args)
		}
	}
}
