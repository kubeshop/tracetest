package global_decorators

import (
	"context"

	global_types "github.com/kubeshop/tracetest/cli/global/types"
	misc_actions "github.com/kubeshop/tracetest/cli/misc/actions"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

type version struct {
	Config
	VersionText *string
}

type Version interface {
	Config
	GetVersionText() string
}

var _ Version = &version{}

func WithVersion(command global_types.Command) global_types.Command {
	config, err := command.(Config)
	if !err {
		panic("command must implement Config interface")
	}

	version := version{
		Config: config,
	}

	cmd := version.Get()
	cmd.PreRun = version.preRun(cmd.PreRun)
	version.Set(cmd)

	return version
}

func (d version) GetVersionText() string {
	return *d.VersionText
}

func (d version) preRun(next CobraFn) CobraFn {
	return func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		config := *d.GetConfig()
		options := []misc_actions.ActionArgsOption{
			misc_actions.ActionWithClient(utils.GetAPIClient(*d.GetConfig())),
			misc_actions.ActionWithConfig(config),
			misc_actions.ActionWithLogger(d.GetLogger()),
		}

		action := misc_actions.NewGetServerVersionAction(options...)
		version := action.GetVersionText(ctx)

		d.VersionText = &version

		next(cmd, args)
	}
}
