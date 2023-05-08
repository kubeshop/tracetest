package cmd

import (
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var dashboardCmd = &cobra.Command{
	GroupID: cmdGroupMisc.ID,
	Use:     "dashboard",
	Short:   "Opens the Tracetest Dashboard URL",
	Long:    "Opens the Tracetest Dashboard URL",
	PreRun:  setupCommand(),
	Run: func(_ *cobra.Command, _ []string) {
		analytics.Track("Dashboard", "cmd", map[string]string{})

		if cliConfig.IsEmpty() {
			cliLogger.Error("Missing Tracetest endpoint configuration")
			os.Exit(1)
			return
		}

		err := utils.OpenBrowser(cliConfig.URL())
		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to open the dashboard url: %s", cliConfig.URL()), zap.Error(err))
			os.Exit(1)
			return
		}

		fmt.Println(fmt.Sprintf("Opening \"%s\" in the default browser", cliConfig.URL()))
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
