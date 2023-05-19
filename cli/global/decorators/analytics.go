package global_decorators

import (
	global_analytics "github.com/kubeshop/tracetest/cli/analytics"
	global_types "github.com/kubeshop/tracetest/cli/global/types"
	"github.com/spf13/cobra"
)

type analytics struct {
	global_types.Command
	Name     string
	Category string
}

type Analytics interface{}

func WithAnalytics(name, category string) func(command global_types.Command) global_types.Command {
	return func(command global_types.Command) global_types.Command {
		analytics := analytics{
			Name:     name,
			Category: category,
		}

		cmd := command.Get()
		cmd.PreRun = analytics.preRun(cmd.PreRun)
		cmd.PostRun = analytics.postRun(cmd.PostRun)

		analytics.Set(cmd)
		return analytics
	}
}

func (d analytics) preRun(next CobraFn) CobraFn {
	return func(cmd *cobra.Command, args []string) {
		global_analytics.Track(d.Name, d.Category, map[string]string{})

		next(cmd, args)
	}
}

func (d analytics) postRun(next CobraFn) CobraFn {
	return func(cmd *cobra.Command, args []string) {
		defer global_analytics.Close()

		next(cmd, args)
	}
}
