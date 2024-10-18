package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/ui"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	GroupID: cmdGroupMisc.ID,
	Use:     "dashboard",
	Short:   "Opens the Tracetest Dashboard URL",
	Long:    "Opens the Tracetest Dashboard URL",
	PreRun:  setupCommand(),
	Run: WithResultHandler(func(_ context.Context, _ *cobra.Command, _ []string) (string, error) {
		if cliConfig.IsEmpty() {
			return "", fmt.Errorf("missing Tracetest endpoint configuration")
		}

		ui := ui.DefaultUI
		err := ui.OpenBrowser(cliConfig.UI())
		if err != nil {
			return "", fmt.Errorf("failed to open the dashboard url: %s", cliConfig.URL())
		}

		return fmt.Sprintf("Opening \"%s\" in the default browser", cliConfig.URL()), nil
	}),
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
