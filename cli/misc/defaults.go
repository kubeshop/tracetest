package misc

import (
	global_middlewares "github.com/kubeshop/tracetest/cli/global/middlewares"
	global_setup "github.com/kubeshop/tracetest/cli/global/setup"
	misc_setup "github.com/kubeshop/tracetest/cli/misc/setup"
	"github.com/spf13/cobra"
)

type args[P any] struct {
	Setup      *misc_setup.Setup
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
	Setup   *misc_setup.Setup
	PreRun  global_middlewares.CobraFn
	Run     global_middlewares.CobraFnWrapper
	PostRun global_middlewares.CobraFn
}

func NewDefaults(name string, setup *global_setup.Setup) defaults {
	newSetup := misc_setup.NewSetup(setup, misc_setup.WithApiClient())
	run := global_middlewares.ComposeRun(global_middlewares.WithAnalytics(name, "cmd"), global_middlewares.WithResultHandler(setup))
	postRun := global_middlewares.ComposeNoopRun(global_middlewares.WithTeardownMiddleware(setup))

	return defaults{
		Setup:   newSetup,
		PreRun:  run(newSetup.PreRun),
		Run:     run,
		PostRun: postRun(),
	}
}
