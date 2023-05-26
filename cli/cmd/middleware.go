package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type RunFn func(cmd *cobra.Command, args []string) (string, error)
type CobraRunFn func(cmd *cobra.Command, args []string)

func WithResultHandler(runFn RunFn) CobraRunFn {
	return func(cmd *cobra.Command, args []string) {
		res, err := runFn(cmd, args)
		if err != nil {
			cliLogger.Error(fmt.Sprintf(`
Version
%s

An error ocurred when executing the command`, versionText), zap.Error(err))
			ExitCLI(1)
			return
		}

		if res != "" {
			fmt.Println(res)
		}
	}
}
