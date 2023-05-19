package global

import (
	global_decorators "github.com/kubeshop/tracetest/cli/global/decorators"
	global_types "github.com/kubeshop/tracetest/cli/global/types"
	"github.com/spf13/cobra"
)

var (
	GroupConfig = &cobra.Group{
		ID:    "configuration",
		Title: "Configuration"
	}

	GroupResources = &cobra.Group{
		ID:    "resources",
		Title: "Resources",
	}

	GroupTests = &cobra.Group{
		ID:    "tests",
		Title: "Tests",
	}

	GroupMisc = &cobra.Group{
		ID:    "misc",
		Title: "Misc",
	}
)

type Root global_types.Command

func NewRoot() Root {
	cmd := &cobra.Command{
		Use:     "tracetest",
		Short:   "CLI to configure, install and execute tests on a Tracetest server",
		Long:    `CLI to configure, install and execute tests on a Tracetest server`,
		PreRun:  global_decorators.NoopRun,
		PostRun: global_decorators.NoopRun,
	}

	cmd.SetCompletionCommandGroupID(GroupConfig.ID)
	cmd.SetHelpCommandGroupID(GroupMisc.ID)

	return global_types.NewCommand(cmd)
}
