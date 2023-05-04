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
		url := fmt.Sprintf("%s://%s", cliConfig.Scheme, cliConfig.Endpoint)

		err := utils.OpenBrowser(url)
		if err != nil {
			cliLogger.Error(fmt.Sprintf("failed to open the dashboard url: %s", url), zap.Error(err))
			os.Exit(1)
			return
		}

		fmt.Println(fmt.Sprintf("Opening \"%s\" in the default browser", url))
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
