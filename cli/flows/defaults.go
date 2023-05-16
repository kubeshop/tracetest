package flows

import (
	global_middlewares "github.com/kubeshop/tracetest/cli/global/middlewares"
	global_setup "github.com/kubeshop/tracetest/cli/global/setup"
	"github.com/spf13/cobra"
)

type args[P any] struct {
	Setup      *global_setup.Setup
	Parameters P
	Cmd        *cobra.Command
}

func NewArgs[P any](defaults defaults, parameters P) args[P] {
	return args[P]{
		Setup:      defaults.Setup,
		Parameters: parameters,
	}
}

type defaults struct {
	Setup   *global_setup.Setup
	PreRun  global_middlewares.CobraFn
	Run     global_middlewares.CobraFnWrapper
	PostRun global_middlewares.CobraFn
}

func NewDefaults(name string, setup *global_setup.Setup) defaults {
	run := global_middlewares.ComposeRun(global_middlewares.WithAnalytics(name, "cmd"), global_middlewares.WithResultHandler(setup))
	postRun := global_middlewares.ComposeNoopRun(global_middlewares.WithTeardownMiddleware(setup))

	return defaults{
		Setup:   setup,
		PreRun:  run(setup.PreRun),
		Run:     run,
		PostRun: postRun(),
	}
}
