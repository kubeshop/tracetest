package global_middlewares

import (
	"fmt"
	"os"

	global_setup "github.com/kubeshop/tracetest/cli/global/setup"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func WithResultHandler(setup *global_setup.Setup) Middleware {
	return func(runFn RunFn) RunFn {
		return func(cmd *cobra.Command, args []string) (string, error) {
			res, err := runFn(cmd, args)
			if err != nil {
				setup.Logger.Error(fmt.Sprintf(`
Version
%s

An error ocurred when executing the command`, *setup.VersionText), zap.Error(err))
				os.Exit(1)
				return res, nil
			}

			if res != "" {
				fmt.Println(res)
			}

			return res, nil
		}
	}
}
